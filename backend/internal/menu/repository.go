package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for menu recommendation data access
type Repository interface {
	FindRecipesByIngredients(ctx context.Context, ingredients []string) ([]models.Recipe, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// FindRecipesByIngredients finds all recipes that contain any of the provided ingredients
func (r *PostgresRepository) FindRecipesByIngredients(ctx context.Context, ingredients []string) ([]models.Recipe, error) {
	if len(ingredients) == 0 {
		return []models.Recipe{}, nil
	}

	// Normalize ingredients for case-insensitive matching
	normalizedIngredients := make([]string, len(ingredients))
	for i, ing := range ingredients {
		normalizedIngredients[i] = strings.ToLower(strings.TrimSpace(ing))
	}

	// Find all recipes that have at least one matching ingredient
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT r.id, r.created_by_user_id, r.title, r.image_url,
		       r.servings, r.serving_size, r.prep_time_minutes,
		       r.cook_time_minutes, r.total_time_minutes,
		       r.storage_notes, r.source, r.notes,
		       r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		INNER JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		WHERE r.deleted_at IS NULL
		  AND LOWER(ri.ingredient_name) = ANY($1)
		ORDER BY r.created_at DESC
	`, normalizedIngredients)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	recipeMap := make(map[uuid.UUID]*models.Recipe)

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

		recipeMap[recipe.ID] = &recipe
	}

	// Load ingredients for each recipe
	for _, recipe := range recipeMap {
		if err := r.loadIngredients(ctx, recipe); err != nil {
			return nil, err
		}
		if err := r.loadInstructions(ctx, recipe); err != nil {
			return nil, err
		}
		if err := r.loadTags(ctx, recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, *recipe)
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
