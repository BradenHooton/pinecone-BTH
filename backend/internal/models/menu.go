package models

import "github.com/google/uuid"

// RecommendRecipesRequest represents the request to get recipe recommendations
type RecommendRecipesRequest struct {
	Ingredients []string `json:"ingredients"`
}

// RecipeRecommendation represents a recommended recipe with match information
type RecipeRecommendation struct {
	Recipe             Recipe   `json:"recipe"`
	MatchScore         float64  `json:"match_score"`          // Percentage of ingredients matched (0-100)
	MatchedIngredients []string `json:"matched_ingredients"`  // Ingredients that match
	MissingIngredients []string `json:"missing_ingredients"`  // Ingredients not in user's list
}

// RecommendRecipesResponse represents the API response for recipe recommendations
type RecommendRecipesResponse struct {
	Data []RecipeRecommendation `json:"data"`
	Meta struct {
		ProvidedIngredients []string `json:"provided_ingredients"`
		TotalRecipesFound   int      `json:"total_recipes_found"`
	} `json:"meta"`
}

// IngredientMatch represents an ingredient match for scoring
type IngredientMatch struct {
	RecipeID           uuid.UUID
	RecipeTitle        string
	IngredientName     string
	TotalIngredients   int
	MatchedIngredients int
}
