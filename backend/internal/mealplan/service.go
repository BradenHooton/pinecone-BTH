package mealplan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

var (
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
	// ErrNotFound is returned when a meal plan is not found
	ErrNotFound = errors.New("meal plan not found")
)

// Service defines the meal plan business logic
type Service struct {
	repo Repository
}

// NewService creates a new meal plan service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetMealPlanByDate retrieves a meal plan for a specific date
func (s *Service) GetMealPlanByDate(ctx context.Context, date time.Time) (*models.MealPlan, error) {
	// Try to get existing meal plan
	mealPlan, err := s.repo.GetMealPlanByDate(ctx, date)
	if err != nil {
		// If not found, create a new empty meal plan
		return s.repo.GetOrCreateMealPlan(ctx, date)
	}
	return mealPlan, nil
}

// GetMealPlansByDateRange retrieves meal plans for a date range
func (s *Service) GetMealPlansByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.MealPlan, error) {
	if startDate.After(endDate) {
		return nil, fmt.Errorf("%w: start date must be before end date", ErrInvalidInput)
	}

	// Limit range to 90 days
	maxDuration := 90 * 24 * time.Hour
	if endDate.Sub(startDate) > maxDuration {
		return nil, fmt.Errorf("%w: date range cannot exceed 90 days", ErrInvalidInput)
	}

	mealPlans, err := s.repo.GetMealPlansByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Fill in missing dates with empty meal plans
	filledPlans := s.fillMissingDates(mealPlans, startDate, endDate)

	return filledPlans, nil
}

// UpdateMealPlan updates a meal plan for a specific date
func (s *Service) UpdateMealPlan(ctx context.Context, date time.Time, req *models.UpdateMealPlanRequest) (*models.MealPlan, error) {
	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	return s.repo.UpdateMealPlan(ctx, date, req.Meals)
}

// validateUpdateRequest validates the update meal plan request
func (s *Service) validateUpdateRequest(req *models.UpdateMealPlanRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	for i, meal := range req.Meals {
		// Validate meal type
		if !isValidMealType(meal.MealType) {
			return fmt.Errorf("meal %d: invalid meal type '%s'", i, meal.MealType)
		}

		// Validate out_of_kitchen constraint
		if meal.OutOfKitchen {
			if meal.RecipeID != nil {
				return fmt.Errorf("meal %d: cannot have recipe_id when out_of_kitchen is true", i)
			}
			if meal.Servings != nil {
				return fmt.Errorf("meal %d: cannot have servings when out_of_kitchen is true", i)
			}
		} else {
			if meal.RecipeID == nil {
				return fmt.Errorf("meal %d: recipe_id required when out_of_kitchen is false", i)
			}
			if meal.Servings == nil {
				return fmt.Errorf("meal %d: servings required when out_of_kitchen is false", i)
			}
			if *meal.Servings <= 0 {
				return fmt.Errorf("meal %d: servings must be greater than 0", i)
			}
		}
	}

	return nil
}

// isValidMealType checks if a meal type is valid
func isValidMealType(mt models.MealType) bool {
	switch mt {
	case models.MealTypeBreakfast, models.MealTypeLunch, models.MealTypeSnack, models.MealTypeDinner, models.MealTypeDessert:
		return true
	default:
		return false
	}
}

// fillMissingDates fills in missing dates in the meal plan range with empty meal plans
func (s *Service) fillMissingDates(mealPlans []models.MealPlan, startDate, endDate time.Time) []models.MealPlan {
	// Create a map of existing meal plans by date
	plansByDate := make(map[string]models.MealPlan)
	for _, mp := range mealPlans {
		dateStr := mp.PlanDate.Format("2006-01-02")
		plansByDate[dateStr] = mp
	}

	// Fill in all dates in the range
	var filledPlans []models.MealPlan
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateStr := currentDate.Format("2006-01-02")

		if mp, exists := plansByDate[dateStr]; exists {
			filledPlans = append(filledPlans, mp)
		} else {
			// Create empty meal plan for missing date
			emptyPlan := models.MealPlan{
				PlanDate: currentDate,
				Meals:    []models.MealPlanRecipe{},
			}
			filledPlans = append(filledPlans, emptyPlan)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return filledPlans
}
