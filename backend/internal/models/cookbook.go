package models

import (
	"time"

	"github.com/google/uuid"
)

// Cookbook represents a collection of recipes
type Cookbook struct {
	ID              uuid.UUID `json:"id"`
	CreatedByUserID uuid.UUID `json:"created_by_user_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	RecipeCount     int       `json:"recipe_count"`
	Recipes         []Recipe  `json:"recipes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// CookbookRecipe represents the junction table entry
type CookbookRecipe struct {
	ID         uuid.UUID `json:"id"`
	CookbookID uuid.UUID `json:"cookbook_id"`
	RecipeID   uuid.UUID `json:"recipe_id"`
	Recipe     *Recipe   `json:"recipe,omitempty"`
	AddedAt    time.Time `json:"added_at"`
}

// CreateCookbookRequest represents the request to create a cookbook
type CreateCookbookRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateCookbookRequest represents the request to update a cookbook
type UpdateCookbookRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// CookbookResponse represents the API response for a single cookbook
type CookbookResponse struct {
	Data Cookbook `json:"data"`
}

// CookbookListResponse represents the API response for a list of cookbooks
type CookbookListResponse struct {
	Data []Cookbook `json:"data"`
	Meta struct {
		Total int64 `json:"total"`
	} `json:"meta"`
}
