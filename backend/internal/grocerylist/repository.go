package grocerylist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for grocery list data access
type Repository interface {
	CreateGroceryList(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*models.GroceryList, error)
	GetGroceryListByID(ctx context.Context, id uuid.UUID) (*models.GroceryList, error)
	GetGroceryListsByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.GroceryList, int64, error)
	DeleteGroceryList(ctx context.Context, id uuid.UUID) error

	CreateGroceryListItem(ctx context.Context, item *models.GroceryListItem) (*models.GroceryListItem, error)
	GetGroceryListItems(ctx context.Context, groceryListID uuid.UUID) ([]models.GroceryListItem, error)
	UpdateItemStatus(ctx context.Context, itemID uuid.UUID, status models.GroceryItemStatus) error
	DeleteGroceryListItem(ctx context.Context, itemID uuid.UUID) error

	GetIngredientsForDateRange(ctx context.Context, startDate, endDate time.Time) ([]IngredientAggregation, error)
}

// IngredientAggregation represents an ingredient from a recipe in a meal plan
type IngredientAggregation struct {
	IngredientName  string
	Quantity        float64
	Unit            string
	Department      models.GroceryDepartment
	RecipeID        uuid.UUID
	RecipeTitle     string
	MealServings    *int
	RecipeServings  int
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateGroceryList creates a new grocery list
func (r *PostgresRepository) CreateGroceryList(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*models.GroceryList, error) {
	// Normalize dates to midnight
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	var groceryList models.GroceryList
	err := r.db.QueryRow(ctx, `
		INSERT INTO grocery_lists (created_by_user_id, start_date, end_date)
		VALUES ($1, $2, $3)
		RETURNING id, created_by_user_id, start_date, end_date, created_at, updated_at
	`, userID, startDate, endDate).Scan(
		&groceryList.ID,
		&groceryList.CreatedByUserID,
		&groceryList.StartDate,
		&groceryList.EndDate,
		&groceryList.CreatedAt,
		&groceryList.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create grocery list: %w", err)
	}

	return &groceryList, nil
}

// GetGroceryListByID retrieves a grocery list by ID
func (r *PostgresRepository) GetGroceryListByID(ctx context.Context, id uuid.UUID) (*models.GroceryList, error) {
	var groceryList models.GroceryList
	err := r.db.QueryRow(ctx, `
		SELECT id, created_by_user_id, start_date, end_date, created_at, updated_at
		FROM grocery_lists
		WHERE id = $1
	`, id).Scan(
		&groceryList.ID,
		&groceryList.CreatedByUserID,
		&groceryList.StartDate,
		&groceryList.EndDate,
		&groceryList.CreatedAt,
		&groceryList.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get grocery list: %w", err)
	}

	// Load items
	items, err := r.GetGroceryListItems(ctx, id)
	if err != nil {
		return nil, err
	}
	groceryList.Items = items

	return &groceryList, nil
}

// GetGroceryListsByUser retrieves grocery lists for a user with pagination
func (r *PostgresRepository) GetGroceryListsByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.GroceryList, int64, error) {
	// Get total count
	var total int64
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM grocery_lists
		WHERE created_by_user_id = $1
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count grocery lists: %w", err)
	}

	// Get grocery lists
	rows, err := r.db.Query(ctx, `
		SELECT id, created_by_user_id, start_date, end_date, created_at, updated_at
		FROM grocery_lists
		WHERE created_by_user_id = $1
		ORDER BY start_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query grocery lists: %w", err)
	}
	defer rows.Close()

	var groceryLists []models.GroceryList
	for rows.Next() {
		var gl models.GroceryList
		if err := rows.Scan(
			&gl.ID,
			&gl.CreatedByUserID,
			&gl.StartDate,
			&gl.EndDate,
			&gl.CreatedAt,
			&gl.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan grocery list: %w", err)
		}

		// Load items for each list
		items, err := r.GetGroceryListItems(ctx, gl.ID)
		if err != nil {
			return nil, 0, err
		}
		gl.Items = items

		groceryLists = append(groceryLists, gl)
	}

	return groceryLists, total, nil
}

// DeleteGroceryList deletes a grocery list
func (r *PostgresRepository) DeleteGroceryList(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM grocery_lists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete grocery list: %w", err)
	}
	return nil
}

// CreateGroceryListItem creates a new grocery list item
func (r *PostgresRepository) CreateGroceryListItem(ctx context.Context, item *models.GroceryListItem) (*models.GroceryListItem, error) {
	var newItem models.GroceryListItem
	err := r.db.QueryRow(ctx, `
		INSERT INTO grocery_list_items (
			grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id
	`, item.GroceryListID, item.ItemName, item.Quantity, item.Unit, item.Department, item.Status, item.IsManual, item.SourceRecipeID).Scan(
		&newItem.ID,
		&newItem.GroceryListID,
		&newItem.ItemName,
		&newItem.Quantity,
		&newItem.Unit,
		&newItem.Department,
		&newItem.Status,
		&newItem.IsManual,
		&newItem.SourceRecipeID,
	)
	if err != nil {
		return nil, fmt.Errorf("create grocery list item: %w", err)
	}

	return &newItem, nil
}

// GetGroceryListItems retrieves all items for a grocery list
func (r *PostgresRepository) GetGroceryListItems(ctx context.Context, groceryListID uuid.UUID) ([]models.GroceryListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id
		FROM grocery_list_items
		WHERE grocery_list_id = $1
		ORDER BY department, item_name
	`, groceryListID)
	if err != nil {
		return nil, fmt.Errorf("query grocery list items: %w", err)
	}
	defer rows.Close()

	var items []models.GroceryListItem
	for rows.Next() {
		var item models.GroceryListItem
		if err := rows.Scan(
			&item.ID,
			&item.GroceryListID,
			&item.ItemName,
			&item.Quantity,
			&item.Unit,
			&item.Department,
			&item.Status,
			&item.IsManual,
			&item.SourceRecipeID,
		); err != nil {
			return nil, fmt.Errorf("scan grocery list item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// UpdateItemStatus updates the status of a grocery list item
func (r *PostgresRepository) UpdateItemStatus(ctx context.Context, itemID uuid.UUID, status models.GroceryItemStatus) error {
	_, err := r.db.Exec(ctx, `
		UPDATE grocery_list_items
		SET status = $2
		WHERE id = $1
	`, itemID, status)
	if err != nil {
		return fmt.Errorf("update item status: %w", err)
	}
	return nil
}

// DeleteGroceryListItem deletes a grocery list item
func (r *PostgresRepository) DeleteGroceryListItem(ctx context.Context, itemID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM grocery_list_items WHERE id = $1`, itemID)
	if err != nil {
		return fmt.Errorf("delete grocery list item: %w", err)
	}
	return nil
}

// GetIngredientsForDateRange retrieves all ingredients from recipes in meal plans within a date range
func (r *PostgresRepository) GetIngredientsForDateRange(ctx context.Context, startDate, endDate time.Time) ([]IngredientAggregation, error) {
	// Normalize dates to midnight
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	rows, err := r.db.Query(ctx, `
		SELECT
			ri.ingredient_name,
			ri.quantity,
			ri.unit,
			ri.department,
			r.id as recipe_id,
			r.title as recipe_title,
			mpr.servings as meal_servings,
			r.servings as recipe_servings
		FROM meal_plans mp
		JOIN meal_plan_recipes mpr ON mp.id = mpr.meal_plan_id
		JOIN recipes r ON mpr.recipe_id = r.id
		JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		WHERE mp.plan_date >= $1
		  AND mp.plan_date <= $2
		  AND mpr.out_of_kitchen = false
		  AND r.deleted_at IS NULL
		ORDER BY ri.department, ri.ingredient_name
	`, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []IngredientAggregation
	for rows.Next() {
		var ing IngredientAggregation
		if err := rows.Scan(
			&ing.IngredientName,
			&ing.Quantity,
			&ing.Unit,
			&ing.Department,
			&ing.RecipeID,
			&ing.RecipeTitle,
			&ing.MealServings,
			&ing.RecipeServings,
		); err != nil {
			return nil, fmt.Errorf("scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ing)
	}

	return ingredients, nil
}

// normalizeUnit normalizes units for comparison (lowercase, trimmed)
func normalizeUnit(unit string) string {
	return strings.TrimSpace(strings.ToLower(unit))
}

// normalizeIngredientName normalizes ingredient names for comparison
func normalizeIngredientName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}
