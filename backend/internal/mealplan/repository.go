package mealplan

import (
	"context"
	"fmt"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for meal plan data access
type Repository interface {
	GetOrCreateMealPlan(ctx context.Context, date time.Time) (*models.MealPlan, error)
	GetMealPlanByDate(ctx context.Context, date time.Time) (*models.MealPlan, error)
	GetMealPlansByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.MealPlan, error)
	UpdateMealPlan(ctx context.Context, date time.Time, meals []models.CreateMealPlanRecipeRequest) (*models.MealPlan, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetOrCreateMealPlan gets an existing meal plan or creates a new one for the given date
func (r *PostgresRepository) GetOrCreateMealPlan(ctx context.Context, date time.Time) (*models.MealPlan, error) {
	// Normalize date to midnight
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	var mealPlan models.MealPlan
	err := r.db.QueryRow(ctx, `
		INSERT INTO meal_plans (plan_date)
		VALUES ($1)
		ON CONFLICT (plan_date) DO UPDATE SET updated_at = NOW()
		RETURNING id, plan_date, created_at, updated_at
	`, date).Scan(&mealPlan.ID, &mealPlan.PlanDate, &mealPlan.CreatedAt, &mealPlan.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get or create meal plan: %w", err)
	}

	// Load meals
	if err := r.loadMeals(ctx, &mealPlan); err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

// GetMealPlanByDate retrieves a meal plan for a specific date
func (r *PostgresRepository) GetMealPlanByDate(ctx context.Context, date time.Time) (*models.MealPlan, error) {
	// Normalize date to midnight
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	var mealPlan models.MealPlan
	err := r.db.QueryRow(ctx, `
		SELECT id, plan_date, created_at, updated_at
		FROM meal_plans
		WHERE plan_date = $1
	`, date).Scan(&mealPlan.ID, &mealPlan.PlanDate, &mealPlan.CreatedAt, &mealPlan.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get meal plan by date: %w", err)
	}

	// Load meals
	if err := r.loadMeals(ctx, &mealPlan); err != nil {
		return nil, err
	}

	return &mealPlan, nil
}

// GetMealPlansByDateRange retrieves meal plans for a date range
func (r *PostgresRepository) GetMealPlansByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.MealPlan, error) {
	// Normalize dates to midnight
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	rows, err := r.db.Query(ctx, `
		SELECT id, plan_date, created_at, updated_at
		FROM meal_plans
		WHERE plan_date >= $1 AND plan_date <= $2
		ORDER BY plan_date ASC
	`, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get meal plans by date range: %w", err)
	}
	defer rows.Close()

	var mealPlans []models.MealPlan
	for rows.Next() {
		var mp models.MealPlan
		if err := rows.Scan(&mp.ID, &mp.PlanDate, &mp.CreatedAt, &mp.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan meal plan: %w", err)
		}

		// Load meals for this plan
		if err := r.loadMeals(ctx, &mp); err != nil {
			return nil, err
		}

		mealPlans = append(mealPlans, mp)
	}

	return mealPlans, nil
}

// UpdateMealPlan updates a meal plan for a specific date
func (r *PostgresRepository) UpdateMealPlan(ctx context.Context, date time.Time, meals []models.CreateMealPlanRecipeRequest) (*models.MealPlan, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Normalize date to midnight
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	// Get or create meal plan
	var mealPlan models.MealPlan
	err = tx.QueryRow(ctx, `
		INSERT INTO meal_plans (plan_date)
		VALUES ($1)
		ON CONFLICT (plan_date) DO UPDATE SET updated_at = NOW()
		RETURNING id, plan_date, created_at, updated_at
	`, date).Scan(&mealPlan.ID, &mealPlan.PlanDate, &mealPlan.CreatedAt, &mealPlan.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get or create meal plan: %w", err)
	}

	// Delete existing meal plan recipes
	_, err = tx.Exec(ctx, `DELETE FROM meal_plan_recipes WHERE meal_plan_id = $1`, mealPlan.ID)
	if err != nil {
		return nil, fmt.Errorf("delete existing meals: %w", err)
	}

	// Insert new meal plan recipes
	for i, meal := range meals {
		var mpr models.MealPlanRecipe
		err = tx.QueryRow(ctx, `
			INSERT INTO meal_plan_recipes (
				meal_plan_id, meal_type, recipe_id, servings, out_of_kitchen, order_index
			) VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, meal_plan_id, meal_type, recipe_id, servings, out_of_kitchen, order_index
		`, mealPlan.ID, meal.MealType, meal.RecipeID, meal.Servings, meal.OutOfKitchen, i).
			Scan(&mpr.ID, &mpr.MealPlanID, &mpr.MealType, &mpr.RecipeID, &mpr.Servings, &mpr.OutOfKitchen, &mpr.OrderIndex)
		if err != nil {
			return nil, fmt.Errorf("insert meal: %w", err)
		}

		mealPlan.Meals = append(mealPlan.Meals, mpr)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Load full recipe data for the meals
	for i := range mealPlan.Meals {
		if mealPlan.Meals[i].RecipeID != nil {
			if err := r.loadRecipe(ctx, &mealPlan.Meals[i]); err != nil {
				// Log error but don't fail - recipe might have been deleted
				fmt.Printf("failed to load recipe: %v\n", err)
			}
		}
	}

	return &mealPlan, nil
}

// loadMeals loads all meals for a meal plan
func (r *PostgresRepository) loadMeals(ctx context.Context, mealPlan *models.MealPlan) error {
	rows, err := r.db.Query(ctx, `
		SELECT id, meal_plan_id, meal_type, recipe_id, servings, out_of_kitchen, order_index
		FROM meal_plan_recipes
		WHERE meal_plan_id = $1
		ORDER BY meal_type ASC, order_index ASC
	`, mealPlan.ID)
	if err != nil {
		return fmt.Errorf("query meals: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mpr models.MealPlanRecipe
		if err := rows.Scan(&mpr.ID, &mpr.MealPlanID, &mpr.MealType, &mpr.RecipeID, &mpr.Servings, &mpr.OutOfKitchen, &mpr.OrderIndex); err != nil {
			return fmt.Errorf("scan meal: %w", err)
		}

		// Load recipe data if recipe_id is present
		if mpr.RecipeID != nil {
			if err := r.loadRecipe(ctx, &mpr); err != nil {
				// Log error but don't fail - recipe might have been deleted
				fmt.Printf("failed to load recipe: %v\n", err)
			}
		}

		mealPlan.Meals = append(mealPlan.Meals, mpr)
	}

	return nil
}

// loadRecipe loads recipe data for a meal plan recipe
func (r *PostgresRepository) loadRecipe(ctx context.Context, mpr *models.MealPlanRecipe) error {
	if mpr.RecipeID == nil {
		return nil
	}

	var recipe models.Recipe
	err := r.db.QueryRow(ctx, `
		SELECT id, created_by_user_id, title, image_url, servings, serving_size,
		       prep_time_minutes, cook_time_minutes, total_time_minutes,
		       storage_notes, source, notes, created_at, updated_at, deleted_at
		FROM recipes
		WHERE id = $1 AND deleted_at IS NULL
	`, *mpr.RecipeID).Scan(&recipe.ID, &recipe.CreatedByUserID, &recipe.Title, &recipe.ImageURL,
		&recipe.Servings, &recipe.ServingSize, &recipe.PrepTimeMinutes,
		&recipe.CookTimeMinutes, &recipe.TotalTimeMinutes, &recipe.StorageNotes,
		&recipe.Source, &recipe.Notes, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.DeletedAt)
	if err != nil {
		return fmt.Errorf("query recipe: %w", err)
	}

	mpr.Recipe = &recipe
	return nil
}
