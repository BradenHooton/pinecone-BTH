package models

import (
	"time"

	"github.com/google/uuid"
)

// NutritionCache represents cached nutrition data from USDA FoodData Central
type NutritionCache struct {
	ID         uuid.UUID  `json:"id"`
	USDAFdcID  string     `json:"usda_fdc_id"`
	FoodName   string     `json:"food_name"`
	Calories   *float64   `json:"calories,omitempty"`
	ProteinG   *float64   `json:"protein_g,omitempty"`
	CarbsG     *float64   `json:"carbs_g,omitempty"`
	FiberG     *float64   `json:"fiber_g,omitempty"`
	FatG       *float64   `json:"fat_g,omitempty"`
	CachedAt   time.Time  `json:"cached_at"`
}

// NutritionSearchResult represents a search result from USDA API
type NutritionSearchResult struct {
	FdcID       string   `json:"fdc_id"`
	Description string   `json:"description"`
	DataType    string   `json:"data_type"`
	Calories    *float64 `json:"calories,omitempty"`
	ProteinG    *float64 `json:"protein_g,omitempty"`
	CarbsG      *float64 `json:"carbs_g,omitempty"`
	FiberG      *float64 `json:"fiber_g,omitempty"`
	FatG        *float64 `json:"fat_g,omitempty"`
}

// NutritionSearchResponse represents the response for nutrition search
type NutritionSearchResponse struct {
	Data []NutritionSearchResult `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

// RecipeNutrition represents calculated nutrition for a recipe
type RecipeNutrition struct {
	TotalCalories  float64 `json:"total_calories"`
	TotalProteinG  float64 `json:"total_protein_g"`
	TotalCarbsG    float64 `json:"total_carbs_g"`
	TotalFiberG    float64 `json:"total_fiber_g"`
	TotalFatG      float64 `json:"total_fat_g"`
	PerServingCalories float64 `json:"per_serving_calories"`
	PerServingProteinG float64 `json:"per_serving_protein_g"`
	PerServingCarbsG   float64 `json:"per_serving_carbs_g"`
	PerServingFiberG   float64 `json:"per_serving_fiber_g"`
	PerServingFatG     float64 `json:"per_serving_fat_g"`
}
