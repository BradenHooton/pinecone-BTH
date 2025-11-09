package nutrition

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/jackc/pgx/v5"
)

var (
	// ErrNotFound is returned when nutrition data is not found
	ErrNotFound = errors.New("nutrition data not found")
)

// Service defines the nutrition business logic
type Service struct {
	repo       Repository
	usdaClient USDAClient
}

// NewService creates a new nutrition service
func NewService(repo Repository, usdaClient USDAClient) *Service {
	return &Service{
		repo:       repo,
		usdaClient: usdaClient,
	}
}

// Search searches for nutrition data
// First checks cache, then falls back to USDA API
func (s *Service) Search(ctx context.Context, query string) ([]models.NutritionSearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("query cannot be empty")
	}

	// Search cache first
	cachedResults, err := s.repo.SearchNutritionByName(ctx, query, 10)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Log error but continue to API search
		fmt.Printf("cache search error: %v\n", err)
	}

	// Convert cached results to search results
	var results []models.NutritionSearchResult
	for _, cached := range cachedResults {
		results = append(results, models.NutritionSearchResult{
			FdcID:       cached.USDAFdcID,
			Description: cached.FoodName,
			DataType:    "Cached",
			Calories:    cached.Calories,
			ProteinG:    cached.ProteinG,
			CarbsG:      cached.CarbsG,
			FiberG:      cached.FiberG,
			FatG:        cached.FatG,
		})
	}

	// If we have cached results, return them
	if len(results) > 0 {
		return results, nil
	}

	// Otherwise, search USDA API
	apiResults, err := s.usdaClient.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("usda api search: %w", err)
	}

	// Cache the API results
	for _, result := range apiResults {
		_, err := s.repo.CreateNutritionCache(
			ctx,
			result.FdcID,
			result.Description,
			result.Calories,
			result.ProteinG,
			result.CarbsG,
			result.FiberG,
			result.FatG,
		)
		if err != nil {
			// Log error but continue (don't fail the request)
			fmt.Printf("failed to cache nutrition data: %v\n", err)
		}
	}

	return apiResults, nil
}

// GetByFdcID retrieves nutrition data by USDA FDC ID
// First checks cache, then falls back to USDA API
func (s *Service) GetByFdcID(ctx context.Context, fdcID string) (*models.NutritionCache, error) {
	// Check cache first
	cached, err := s.repo.GetNutritionByFdcID(ctx, fdcID)
	if err == nil {
		return cached, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get from cache: %w", err)
	}

	// Not in cache, search API
	results, err := s.usdaClient.Search(ctx, fdcID)
	if err != nil {
		return nil, fmt.Errorf("usda api search: %w", err)
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	// Use the first result
	result := results[0]

	// Cache it
	cached, err = s.repo.CreateNutritionCache(
		ctx,
		result.FdcID,
		result.Description,
		result.Calories,
		result.ProteinG,
		result.CarbsG,
		result.FiberG,
		result.FatG,
	)
	if err != nil {
		return nil, fmt.Errorf("cache nutrition data: %w", err)
	}

	return cached, nil
}

// CalculateRecipeNutrition calculates total and per-serving nutrition for a recipe
func (s *Service) CalculateRecipeNutrition(ctx context.Context, recipe *models.Recipe) (*models.RecipeNutrition, error) {
	if recipe == nil {
		return nil, errors.New("recipe cannot be nil")
	}

	nutrition := &models.RecipeNutrition{}

	// Sum up nutrition from all ingredients that have nutrition data
	for _, ingredient := range recipe.Ingredients {
		if ingredient.NutritionID == nil {
			// Skip ingredients without nutrition data
			continue
		}

		// Get nutrition data
		nutritionData, err := s.repo.GetNutritionByID(ctx, *ingredient.NutritionID)
		if err != nil {
			// Skip if not found, but log the error
			fmt.Printf("nutrition data not found for ingredient %s: %v\n", ingredient.IngredientName, err)
			continue
		}

		// Add to totals (nutrition data is per 100g, so multiply by quantity/100)
		// Note: This is a simplified calculation. In a real app, you'd need unit conversion
		scalingFactor := ingredient.Quantity / 100.0

		if nutritionData.Calories != nil {
			nutrition.TotalCalories += *nutritionData.Calories * scalingFactor
		}
		if nutritionData.ProteinG != nil {
			nutrition.TotalProteinG += *nutritionData.ProteinG * scalingFactor
		}
		if nutritionData.CarbsG != nil {
			nutrition.TotalCarbsG += *nutritionData.CarbsG * scalingFactor
		}
		if nutritionData.FiberG != nil {
			nutrition.TotalFiberG += *nutritionData.FiberG * scalingFactor
		}
		if nutritionData.FatG != nil {
			nutrition.TotalFatG += *nutritionData.FatG * scalingFactor
		}
	}

	// Calculate per-serving values
	if recipe.Servings > 0 {
		servings := float64(recipe.Servings)
		nutrition.PerServingCalories = nutrition.TotalCalories / servings
		nutrition.PerServingProteinG = nutrition.TotalProteinG / servings
		nutrition.PerServingCarbsG = nutrition.TotalCarbsG / servings
		nutrition.PerServingFiberG = nutrition.TotalFiberG / servings
		nutrition.PerServingFatG = nutrition.TotalFatG / servings
	}

	return nutrition, nil
}
