package cookbook

import (
	"context"
	"fmt"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for cookbook data access
type Repository interface {
	CreateCookbook(ctx context.Context, userID uuid.UUID, req *models.CreateCookbookRequest) (*models.Cookbook, error)
	GetCookbookByID(ctx context.Context, id uuid.UUID) (*models.Cookbook, error)
	GetCookbooksByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Cookbook, int64, error)
	UpdateCookbook(ctx context.Context, id uuid.UUID, req *models.UpdateCookbookRequest) (*models.Cookbook, error)
	DeleteCookbook(ctx context.Context, id uuid.UUID) error

	AddRecipeToCookbook(ctx context.Context, cookbookID, recipeID uuid.UUID) error
	RemoveRecipeFromCookbook(ctx context.Context, cookbookID, recipeID uuid.UUID) error
	GetCookbookRecipes(ctx context.Context, cookbookID uuid.UUID) ([]models.Recipe, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateCookbook creates a new cookbook
func (r *PostgresRepository) CreateCookbook(ctx context.Context, userID uuid.UUID, req *models.CreateCookbookRequest) (*models.Cookbook, error) {
	var cookbook models.Cookbook
	err := r.db.QueryRow(ctx, `
		INSERT INTO cookbooks (created_by_user_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_by_user_id, name, description, created_at, updated_at, deleted_at
	`, userID, req.Name, req.Description).Scan(
		&cookbook.ID,
		&cookbook.CreatedByUserID,
		&cookbook.Name,
		&cookbook.Description,
		&cookbook.CreatedAt,
		&cookbook.UpdatedAt,
		&cookbook.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create cookbook: %w", err)
	}

	cookbook.RecipeCount = 0
	cookbook.Recipes = []models.Recipe{}

	return &cookbook, nil
}

// GetCookbookByID retrieves a cookbook by ID with all its recipes
func (r *PostgresRepository) GetCookbookByID(ctx context.Context, id uuid.UUID) (*models.Cookbook, error) {
	var cookbook models.Cookbook
	err := r.db.QueryRow(ctx, `
		SELECT id, created_by_user_id, name, description, created_at, updated_at, deleted_at
		FROM cookbooks
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&cookbook.ID,
		&cookbook.CreatedByUserID,
		&cookbook.Name,
		&cookbook.Description,
		&cookbook.CreatedAt,
		&cookbook.UpdatedAt,
		&cookbook.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get cookbook: %w", err)
	}

	// Get recipe count
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM cookbook_recipes
		WHERE cookbook_id = $1
	`, id).Scan(&cookbook.RecipeCount)
	if err != nil {
		return nil, fmt.Errorf("get recipe count: %w", err)
	}

	// Load recipes
	recipes, err := r.GetCookbookRecipes(ctx, id)
	if err != nil {
		return nil, err
	}
	cookbook.Recipes = recipes

	return &cookbook, nil
}

// GetCookbooksByUser retrieves cookbooks for a user with pagination
func (r *PostgresRepository) GetCookbooksByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Cookbook, int64, error) {
	// Get total count
	var total int64
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM cookbooks
		WHERE created_by_user_id = $1 AND deleted_at IS NULL
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count cookbooks: %w", err)
	}

	// Get cookbooks
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.created_by_user_id, c.name, c.description, c.created_at, c.updated_at, c.deleted_at,
		       COALESCE(COUNT(cr.id), 0) as recipe_count
		FROM cookbooks c
		LEFT JOIN cookbook_recipes cr ON c.id = cr.cookbook_id
		WHERE c.created_by_user_id = $1 AND c.deleted_at IS NULL
		GROUP BY c.id, c.created_by_user_id, c.name, c.description, c.created_at, c.updated_at, c.deleted_at
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query cookbooks: %w", err)
	}
	defer rows.Close()

	var cookbooks []models.Cookbook
	for rows.Next() {
		var cb models.Cookbook
		if err := rows.Scan(
			&cb.ID,
			&cb.CreatedByUserID,
			&cb.Name,
			&cb.Description,
			&cb.CreatedAt,
			&cb.UpdatedAt,
			&cb.DeletedAt,
			&cb.RecipeCount,
		); err != nil {
			return nil, 0, fmt.Errorf("scan cookbook: %w", err)
		}
		cookbooks = append(cookbooks, cb)
	}

	return cookbooks, total, nil
}

// UpdateCookbook updates a cookbook
func (r *PostgresRepository) UpdateCookbook(ctx context.Context, id uuid.UUID, req *models.UpdateCookbookRequest) (*models.Cookbook, error) {
	var cookbook models.Cookbook
	err := r.db.QueryRow(ctx, `
		UPDATE cookbooks
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, created_by_user_id, name, description, created_at, updated_at, deleted_at
	`, id, req.Name, req.Description).Scan(
		&cookbook.ID,
		&cookbook.CreatedByUserID,
		&cookbook.Name,
		&cookbook.Description,
		&cookbook.CreatedAt,
		&cookbook.UpdatedAt,
		&cookbook.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update cookbook: %w", err)
	}

	// Get recipe count
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM cookbook_recipes
		WHERE cookbook_id = $1
	`, id).Scan(&cookbook.RecipeCount)
	if err != nil {
		return nil, fmt.Errorf("get recipe count: %w", err)
	}

	return &cookbook, nil
}

// DeleteCookbook soft-deletes a cookbook
func (r *PostgresRepository) DeleteCookbook(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `
		UPDATE cookbooks
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("delete cookbook: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cookbook not found or already deleted")
	}

	return nil
}

// AddRecipeToCookbook adds a recipe to a cookbook
func (r *PostgresRepository) AddRecipeToCookbook(ctx context.Context, cookbookID, recipeID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO cookbook_recipes (cookbook_id, recipe_id)
		VALUES ($1, $2)
		ON CONFLICT (cookbook_id, recipe_id) DO NOTHING
	`, cookbookID, recipeID)
	if err != nil {
		return fmt.Errorf("add recipe to cookbook: %w", err)
	}

	return nil
}

// RemoveRecipeFromCookbook removes a recipe from a cookbook
func (r *PostgresRepository) RemoveRecipeFromCookbook(ctx context.Context, cookbookID, recipeID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM cookbook_recipes
		WHERE cookbook_id = $1 AND recipe_id = $2
	`, cookbookID, recipeID)
	if err != nil {
		return fmt.Errorf("remove recipe from cookbook: %w", err)
	}

	return nil
}

// GetCookbookRecipes retrieves all recipes in a cookbook
func (r *PostgresRepository) GetCookbookRecipes(ctx context.Context, cookbookID uuid.UUID) ([]models.Recipe, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.created_by_user_id, r.title, r.image_url,
		       r.servings, r.serving_size, r.prep_time_minutes,
		       r.cook_time_minutes, r.total_time_minutes,
		       r.storage_notes, r.source, r.notes,
		       r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		INNER JOIN cookbook_recipes cr ON r.id = cr.recipe_id
		WHERE cr.cookbook_id = $1 AND r.deleted_at IS NULL
		ORDER BY cr.added_at DESC
	`, cookbookID)
	if err != nil {
		return nil, fmt.Errorf("query cookbook recipes: %w", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(
			&recipe.ID,
			&recipe.CreatedByUserID,
			&recipe.Title,
			&recipe.ImageURL,
			&recipe.Servings,
			&recipe.ServingSize,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.TotalTimeMinutes,
			&recipe.StorageNotes,
			&recipe.Source,
			&recipe.Notes,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
			&recipe.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}

		// Load ingredients
		if err := r.loadIngredients(ctx, &recipe); err != nil {
			return nil, err
		}

		// Load instructions
		if err := r.loadInstructions(ctx, &recipe); err != nil {
			return nil, err
		}

		// Load tags
		if err := r.loadTags(ctx, &recipe); err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// loadIngredients loads all ingredients for a recipe
func (r *PostgresRepository) loadIngredients(ctx context.Context, recipe *models.Recipe) error {
	rows, err := r.db.Query(ctx, `
		SELECT id, recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
		FROM recipe_ingredients
		WHERE recipe_id = $1
		ORDER BY order_index ASC
	`, recipe.ID)
	if err != nil {
		return fmt.Errorf("query ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient models.RecipeIngredient
		if err := rows.Scan(
			&ingredient.ID,
			&ingredient.RecipeID,
			&ingredient.NutritionID,
			&ingredient.IngredientName,
			&ingredient.Quantity,
			&ingredient.Unit,
			&ingredient.Department,
			&ingredient.OrderIndex,
		); err != nil {
			return fmt.Errorf("scan ingredient: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	return nil
}

// loadInstructions loads all instructions for a recipe
func (r *PostgresRepository) loadInstructions(ctx context.Context, recipe *models.Recipe) error {
	rows, err := r.db.Query(ctx, `
		SELECT id, recipe_id, step_number, instruction
		FROM recipe_instructions
		WHERE recipe_id = $1
		ORDER BY step_number ASC
	`, recipe.ID)
	if err != nil {
		return fmt.Errorf("query instructions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var instruction models.RecipeInstruction
		if err := rows.Scan(
			&instruction.ID,
			&instruction.RecipeID,
			&instruction.StepNumber,
			&instruction.Instruction,
		); err != nil {
			return fmt.Errorf("scan instruction: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, instruction)
	}

	return nil
}

// loadTags loads all tags for a recipe
func (r *PostgresRepository) loadTags(ctx context.Context, recipe *models.Recipe) error {
	rows, err := r.db.Query(ctx, `
		SELECT id, recipe_id, tag_name
		FROM recipe_tags
		WHERE recipe_id = $1
		ORDER BY tag_name ASC
	`, recipe.ID)
	if err != nil {
		return fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.RecipeTag
		if err := rows.Scan(
			&tag.ID,
			&tag.RecipeID,
			&tag.TagName,
		); err != nil {
			return fmt.Errorf("scan tag: %w", err)
		}
		recipe.Tags = append(recipe.Tags, tag)
	}

	return nil
}
