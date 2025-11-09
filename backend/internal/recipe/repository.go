package recipe

import (
	"context"
	"fmt"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for recipe data access
type Repository interface {
	CreateRecipe(ctx context.Context, userID uuid.UUID, req *models.CreateRecipeRequest) (*models.Recipe, error)
	GetRecipeByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error)
	ListRecipes(ctx context.Context, params *models.RecipeSearchParams) ([]models.Recipe, int64, error)
	UpdateRecipe(ctx context.Context, id uuid.UUID, req *models.UpdateRecipeRequest) (*models.Recipe, error)
	DeleteRecipe(ctx context.Context, id uuid.UUID) error
	GetRecipesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Recipe, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateRecipe creates a new recipe with all its related data
func (r *PostgresRepository) CreateRecipe(ctx context.Context, userID uuid.UUID, req *models.CreateRecipeRequest) (*models.Recipe, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create the recipe
	var recipe models.Recipe
	err = tx.QueryRow(ctx, `
		INSERT INTO recipes (
			created_by_user_id, title, image_url, servings, serving_size,
			prep_time_minutes, cook_time_minutes, storage_notes, source, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_by_user_id, title, image_url, servings, serving_size,
		          prep_time_minutes, cook_time_minutes, total_time_minutes,
		          storage_notes, source, notes, created_at, updated_at, deleted_at
	`, userID, req.Title, req.ImageURL, req.Servings, req.ServingSize,
		req.PrepTimeMinutes, req.CookTimeMinutes, req.StorageNotes, req.Source, req.Notes).
		Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
			&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
			&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("insert recipe: %w", err)
	}

	// Create ingredients
	for i, ing := range req.Ingredients {
		var ingredient models.RecipeIngredient
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_ingredients (
				recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
		`, recipe.ID, ing.NutritionID, ing.IngredientName, ing.Quantity, ing.Unit, ing.Department, i).
			Scan(&ingredient.ID, &ingredient.RecipeID, &ingredient.NutritionID,
				&ingredient.IngredientName, &ingredient.Quantity, &ingredient.Unit,
				&ingredient.Department, &ingredient.OrderIndex)
		if err != nil {
			return nil, fmt.Errorf("insert ingredient: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Create instructions
	for _, inst := range req.Instructions {
		var instruction models.RecipeInstruction
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_instructions (recipe_id, step_number, instruction)
			VALUES ($1, $2, $3)
			RETURNING id, recipe_id, step_number, instruction
		`, recipe.ID, inst.StepNumber, inst.Instruction).
			Scan(&instruction.ID, &instruction.RecipeID, &instruction.StepNumber, &instruction.Instruction)
		if err != nil {
			return nil, fmt.Errorf("insert instruction: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, instruction)
	}

	// Create tags
	for _, tagName := range req.Tags {
		var tag models.RecipeTag
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_tags (recipe_id, tag_name)
			VALUES ($1, $2)
			RETURNING id, recipe_id, tag_name
		`, recipe.ID, tagName).
			Scan(&tag.ID, &tag.RecipeID, &tag.TagName)
		if err != nil {
			return nil, fmt.Errorf("insert tag: %w", err)
		}
		recipe.Tags = append(recipe.Tags, tag)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &recipe, nil
}

// GetRecipeByID retrieves a recipe with all its related data
func (r *PostgresRepository) GetRecipeByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.QueryRow(ctx, `
		SELECT id, created_by_user_id, title, image_url, servings, serving_size,
		       prep_time_minutes, cook_time_minutes, total_time_minutes,
		       storage_notes, source, notes, created_at, updated_at, deleted_at
		FROM recipes
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
		&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
		&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
		&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("query recipe: %w", err)
	}

	// Get ingredients
	rows, err := r.db.Query(ctx, `
		SELECT id, recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
		FROM recipe_ingredients
		WHERE recipe_id = $1
		ORDER BY order_index ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("query ingredients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ing models.RecipeIngredient
		if err := rows.Scan(&ing.ID, &ing.RecipeID, &ing.NutritionID, &ing.IngredientName,
			&ing.Quantity, &ing.Unit, &ing.Department, &ing.OrderIndex); err != nil {
			return nil, fmt.Errorf("scan ingredient: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, ing)
	}

	// Get instructions
	rows, err = r.db.Query(ctx, `
		SELECT id, recipe_id, step_number, instruction
		FROM recipe_instructions
		WHERE recipe_id = $1
		ORDER BY step_number ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("query instructions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var inst models.RecipeInstruction
		if err := rows.Scan(&inst.ID, &inst.RecipeID, &inst.StepNumber, &inst.Instruction); err != nil {
			return nil, fmt.Errorf("scan instruction: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, inst)
	}

	// Get tags
	rows, err = r.db.Query(ctx, `
		SELECT id, recipe_id, tag_name
		FROM recipe_tags
		WHERE recipe_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.RecipeTag
		if err := rows.Scan(&tag.ID, &tag.RecipeID, &tag.TagName); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		recipe.Tags = append(recipe.Tags, tag)
	}

	return &recipe, nil
}

// ListRecipes retrieves recipes with optional search/filter/sort
func (r *PostgresRepository) ListRecipes(ctx context.Context, params *models.RecipeSearchParams) ([]models.Recipe, int64, error) {
	query := `
		SELECT DISTINCT r.id, r.created_by_user_id, r.title, r.image_url, r.servings, r.serving_size,
		       r.prep_time_minutes, r.cook_time_minutes, r.total_time_minutes,
		       r.storage_notes, r.source, r.notes, r.created_at, r.updated_at, r.deleted_at
		FROM recipes r
		LEFT JOIN recipe_tags rt ON r.id = rt.recipe_id
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		WHERE r.deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 0

	// Add search filter
	if params.Query != "" {
		argCount++
		query += fmt.Sprintf(" AND (r.title ILIKE $%d OR ri.ingredient_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+params.Query+"%")
	}

	// Add tag filter
	if len(params.Tags) > 0 {
		argCount++
		query += fmt.Sprintf(" AND rt.tag_name = ANY($%d)", argCount)
		args = append(args, params.Tags)
	}

	// Add sorting
	switch params.Sort {
	case "title_asc":
		query += " ORDER BY r.title ASC"
	case "title_desc":
		query += " ORDER BY r.title DESC"
	case "date_asc":
		query += " ORDER BY r.created_at ASC"
	case "date_desc", "":
		query += " ORDER BY r.created_at DESC"
	}

	// Add pagination
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, params.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, params.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
			&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
			&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt); err != nil {
			return nil, 0, fmt.Errorf("scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(DISTINCT r.id) FROM recipes r LEFT JOIN recipe_tags rt ON r.id = rt.recipe_id LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id WHERE r.deleted_at IS NULL"
	countArgs := []interface{}{}
	argCount = 0

	if params.Query != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND (r.title ILIKE $%d OR ri.ingredient_name ILIKE $%d)", argCount, argCount)
		countArgs = append(countArgs, "%"+params.Query+"%")
	}

	if len(params.Tags) > 0 {
		argCount++
		countQuery += fmt.Sprintf(" AND rt.tag_name = ANY($%d)", argCount)
		countArgs = append(countArgs, params.Tags)
	}

	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count recipes: %w", err)
	}

	return recipes, total, nil
}

// UpdateRecipe updates a recipe and all its related data
func (r *PostgresRepository) UpdateRecipe(ctx context.Context, id uuid.UUID, req *models.UpdateRecipeRequest) (*models.Recipe, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update the recipe
	var recipe models.Recipe
	err = tx.QueryRow(ctx, `
		UPDATE recipes
		SET title = $2, image_url = $3, servings = $4, serving_size = $5,
		    prep_time_minutes = $6, cook_time_minutes = $7, storage_notes = $8,
		    source = $9, notes = $10, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, created_by_user_id, title, image_url, servings, serving_size,
		          prep_time_minutes, cook_time_minutes, total_time_minutes,
		          storage_notes, source, notes, created_at, updated_at, deleted_at
	`, id, req.Title, req.ImageURL, req.Servings, req.ServingSize,
		req.PrepTimeMinutes, req.CookTimeMinutes, req.StorageNotes, req.Source, req.Notes).
		Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
			&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
			&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("update recipe: %w", err)
	}

	// Delete old ingredients, instructions, and tags
	_, err = tx.Exec(ctx, "DELETE FROM recipe_ingredients WHERE recipe_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("delete ingredients: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM recipe_instructions WHERE recipe_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("delete instructions: %w", err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM recipe_tags WHERE recipe_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("delete tags: %w", err)
	}

	// Create new ingredients
	for i, ing := range req.Ingredients {
		var ingredient models.RecipeIngredient
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_ingredients (
				recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, recipe_id, nutrition_id, ingredient_name, quantity, unit, department, order_index
		`, id, ing.NutritionID, ing.IngredientName, ing.Quantity, ing.Unit, ing.Department, i).
			Scan(&ingredient.ID, &ingredient.RecipeID, &ingredient.NutritionID,
				&ingredient.IngredientName, &ingredient.Quantity, &ingredient.Unit,
				&ingredient.Department, &ingredient.OrderIndex)
		if err != nil {
			return nil, fmt.Errorf("insert ingredient: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Create new instructions
	for _, inst := range req.Instructions {
		var instruction models.RecipeInstruction
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_instructions (recipe_id, step_number, instruction)
			VALUES ($1, $2, $3)
			RETURNING id, recipe_id, step_number, instruction
		`, id, inst.StepNumber, inst.Instruction).
			Scan(&instruction.ID, &instruction.RecipeID, &instruction.StepNumber, &instruction.Instruction)
		if err != nil {
			return nil, fmt.Errorf("insert instruction: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, instruction)
	}

	// Create new tags
	for _, tagName := range req.Tags {
		var tag models.RecipeTag
		err = tx.QueryRow(ctx, `
			INSERT INTO recipe_tags (recipe_id, tag_name)
			VALUES ($1, $2)
			RETURNING id, recipe_id, tag_name
		`, id, tagName).
			Scan(&tag.ID, &tag.RecipeID, &tag.TagName)
		if err != nil {
			return nil, fmt.Errorf("insert tag: %w", err)
		}
		recipe.Tags = append(recipe.Tags, tag)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &recipe, nil
}

// DeleteRecipe soft-deletes a recipe
func (r *PostgresRepository) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE recipes
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete recipe: %w", err)
	}
	return nil
}

// GetRecipesByUserID retrieves recipes created by a specific user
func (r *PostgresRepository) GetRecipesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Recipe, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, created_by_user_id, title, image_url, servings, serving_size,
		       prep_time_minutes, cook_time_minutes, total_time_minutes,
		       storage_notes, source, notes, created_at, updated_at, deleted_at
		FROM recipes
		WHERE created_by_user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
			&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
			&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
