package models

import (
	"time"

	"github.com/google/uuid"
)

// MealType represents the type of meal
type MealType string

const (
	MealTypeBreakfast MealType = "breakfast"
	MealTypeLunch     MealType = "lunch"
	MealTypeSnack     MealType = "snack"
	MealTypeDinner    MealType = "dinner"
	MealTypeDessert   MealType = "dessert"
)

// MealPlan represents a meal plan for a specific date
type MealPlan struct {
	ID        uuid.UUID  `json:"id"`
	PlanDate  time.Time  `json:"plan_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Meals     []MealPlanRecipe `json:"meals,omitempty"`
}

// MealPlanRecipe represents a recipe in a meal slot
type MealPlanRecipe struct {
	ID            uuid.UUID  `json:"id"`
	MealPlanID    uuid.UUID  `json:"meal_plan_id"`
	MealType      MealType   `json:"meal_type"`
	RecipeID      *uuid.UUID `json:"recipe_id,omitempty"`
	Recipe        *Recipe    `json:"recipe,omitempty"`
	Servings      *int       `json:"servings,omitempty"`
	OutOfKitchen  bool       `json:"out_of_kitchen"`
	OrderIndex    int        `json:"order_index"`
}

// CreateMealPlanRecipeRequest represents a request to add a recipe to a meal slot
type CreateMealPlanRecipeRequest struct {
	MealType     MealType   `json:"meal_type"`
	RecipeID     *uuid.UUID `json:"recipe_id,omitempty"`
	Servings     *int       `json:"servings,omitempty"`
	OutOfKitchen bool       `json:"out_of_kitchen"`
}

// UpdateMealPlanRequest represents a request to update a meal plan for a specific date
type UpdateMealPlanRequest struct {
	Meals []CreateMealPlanRecipeRequest `json:"meals"`
}

// MealPlanResponse represents the response for a single meal plan
type MealPlanResponse struct {
	Data MealPlan `json:"data"`
}

// MealPlanListResponse represents the response for multiple meal plans
type MealPlanListResponse struct {
	Data []MealPlan `json:"data"`
	Meta struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"meta"`
}
