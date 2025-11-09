package models

import (
	"time"

	"github.com/google/uuid"
)

// Recipe represents a recipe in the system
type Recipe struct {
	ID                uuid.UUID  `json:"id"`
	CreatedByUserID   uuid.UUID  `json:"created_by_user_id"`
	Title             string     `json:"title"`
	ImageURL          *string    `json:"image_url,omitempty"`
	Servings          int        `json:"servings"`
	ServingSize       string     `json:"serving_size"`
	PrepTimeMinutes   *int       `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes   *int       `json:"cook_time_minutes,omitempty"`
	TotalTimeMinutes  int        `json:"total_time_minutes"`
	StorageNotes      *string    `json:"storage_notes,omitempty"`
	Source            *string    `json:"source,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Ingredients       []RecipeIngredient   `json:"ingredients,omitempty"`
	Instructions      []RecipeInstruction  `json:"instructions,omitempty"`
	Tags              []RecipeTag          `json:"tags,omitempty"`
}

// RecipeIngredient represents an ingredient in a recipe
type RecipeIngredient struct {
	ID             uuid.UUID  `json:"id"`
	RecipeID       uuid.UUID  `json:"recipe_id"`
	NutritionID    *uuid.UUID `json:"nutrition_id,omitempty"`
	IngredientName string     `json:"ingredient_name"`
	Quantity       float64    `json:"quantity"`
	Unit           string     `json:"unit"`
	Department     string     `json:"department"`
	OrderIndex     int        `json:"order_index"`
}

// RecipeInstruction represents a cooking instruction step
type RecipeInstruction struct {
	ID          uuid.UUID `json:"id"`
	RecipeID    uuid.UUID `json:"recipe_id"`
	StepNumber  int       `json:"step_number"`
	Instruction string    `json:"instruction"`
}

// RecipeTag represents a tag for categorizing recipes
type RecipeTag struct {
	ID       uuid.UUID `json:"id"`
	RecipeID uuid.UUID `json:"recipe_id"`
	TagName  string    `json:"tag_name"`
}

// CreateRecipeRequest represents the request to create a new recipe
type CreateRecipeRequest struct {
	Title           string                     `json:"title"`
	ImageURL        *string                    `json:"image_url,omitempty"`
	Servings        int                        `json:"servings"`
	ServingSize     string                     `json:"serving_size"`
	PrepTimeMinutes *int                       `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes *int                       `json:"cook_time_minutes,omitempty"`
	StorageNotes    *string                    `json:"storage_notes,omitempty"`
	Source          *string                    `json:"source,omitempty"`
	Notes           *string                    `json:"notes,omitempty"`
	Ingredients     []CreateIngredientRequest  `json:"ingredients"`
	Instructions    []CreateInstructionRequest `json:"instructions"`
	Tags            []string                   `json:"tags,omitempty"`
}

// CreateIngredientRequest represents an ingredient in the create request
type CreateIngredientRequest struct {
	IngredientName string     `json:"ingredient_name"`
	Quantity       float64    `json:"quantity"`
	Unit           string     `json:"unit"`
	Department     string     `json:"department"`
	NutritionID    *uuid.UUID `json:"nutrition_id,omitempty"`
}

// CreateInstructionRequest represents an instruction in the create request
type CreateInstructionRequest struct {
	StepNumber  int    `json:"step_number"`
	Instruction string `json:"instruction"`
}

// UpdateRecipeRequest represents the request to update a recipe
type UpdateRecipeRequest struct {
	Title           string                     `json:"title"`
	ImageURL        *string                    `json:"image_url,omitempty"`
	Servings        int                        `json:"servings"`
	ServingSize     string                     `json:"serving_size"`
	PrepTimeMinutes *int                       `json:"prep_time_minutes,omitempty"`
	CookTimeMinutes *int                       `json:"cook_time_minutes,omitempty"`
	StorageNotes    *string                    `json:"storage_notes,omitempty"`
	Source          *string                    `json:"source,omitempty"`
	Notes           *string                    `json:"notes,omitempty"`
	Ingredients     []CreateIngredientRequest  `json:"ingredients"`
	Instructions    []CreateInstructionRequest `json:"instructions"`
	Tags            []string                   `json:"tags,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// RecipeListResponse represents paginated recipe list response
type RecipeListResponse struct {
	Data []Recipe       `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// RecipeSearchParams represents search parameters for recipes
type RecipeSearchParams struct {
	Query      string
	Tags       []string
	Sort       string // "title_asc", "title_desc", "date_asc", "date_desc"
	Limit      int
	Offset     int
}
