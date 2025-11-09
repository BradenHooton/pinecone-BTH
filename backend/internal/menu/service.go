package menu

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

// Service handles business logic for menu recommendations
type Service struct {
	repo Repository
}

// NewService creates a new menu recommendation service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// RecommendRecipes recommends recipes based on available ingredients
func (s *Service) RecommendRecipes(ctx context.Context, req *models.RecommendRecipesRequest) ([]models.RecipeRecommendation, error) {
	// Validate input
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Normalize provided ingredients
	providedIngredients := s.normalizeIngredients(req.Ingredients)

	// Find recipes that contain any of the provided ingredients
	recipes, err := s.repo.FindRecipesByIngredients(ctx, req.Ingredients)
	if err != nil {
		return nil, err
	}

	// Score and rank recipes
	recommendations := make([]models.RecipeRecommendation, 0, len(recipes))

	for _, recipe := range recipes {
		recommendation := s.scoreRecipe(recipe, providedIngredients)
		recommendations = append(recommendations, recommendation)
	}

	// Sort by match score (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].MatchScore > recommendations[j].MatchScore
	})

	return recommendations, nil
}

// scoreRecipe calculates the match score for a recipe
func (s *Service) scoreRecipe(recipe models.Recipe, providedIngredients map[string]bool) models.RecipeRecommendation {
	if len(recipe.Ingredients) == 0 {
		return models.RecipeRecommendation{
			Recipe:             recipe,
			MatchScore:         0,
			MatchedIngredients: []string{},
			MissingIngredients: []string{},
		}
	}

	var matchedIngredients []string
	var missingIngredients []string

	// Check each recipe ingredient against provided ingredients
	for _, ingredient := range recipe.Ingredients {
		normalizedName := s.normalizeIngredientName(ingredient.IngredientName)
		if providedIngredients[normalizedName] {
			matchedIngredients = append(matchedIngredients, ingredient.IngredientName)
		} else {
			missingIngredients = append(missingIngredients, ingredient.IngredientName)
		}
	}

	// Calculate match score: (matched / total) * 100
	matchScore := (float64(len(matchedIngredients)) / float64(len(recipe.Ingredients))) * 100

	return models.RecipeRecommendation{
		Recipe:             recipe,
		MatchScore:         matchScore,
		MatchedIngredients: matchedIngredients,
		MissingIngredients: missingIngredients,
	}
}

// normalizeIngredients creates a map of normalized ingredient names for quick lookup
func (s *Service) normalizeIngredients(ingredients []string) map[string]bool {
	normalized := make(map[string]bool)
	for _, ing := range ingredients {
		normalizedName := s.normalizeIngredientName(ing)
		normalized[normalizedName] = true
	}
	return normalized
}

// normalizeIngredientName normalizes an ingredient name for comparison
func (s *Service) normalizeIngredientName(name string) string {
	// Convert to lowercase and trim whitespace
	normalized := strings.ToLower(strings.TrimSpace(name))

	// Remove common plurals (basic approach)
	// In a production system, you might use a more sophisticated stemming algorithm
	if strings.HasSuffix(normalized, "ies") {
		normalized = strings.TrimSuffix(normalized, "ies") + "y"
	} else if strings.HasSuffix(normalized, "es") {
		normalized = strings.TrimSuffix(normalized, "es")
	} else if strings.HasSuffix(normalized, "s") && !strings.HasSuffix(normalized, "ss") {
		normalized = strings.TrimSuffix(normalized, "s")
	}

	return normalized
}

// validateRequest validates the recommendation request
func (s *Service) validateRequest(req *models.RecommendRecipesRequest) error {
	if len(req.Ingredients) == 0 {
		return errors.New("at least one ingredient is required")
	}

	// Check for empty ingredient names
	for _, ing := range req.Ingredients {
		if strings.TrimSpace(ing) == "" {
			return errors.New("ingredient names cannot be empty")
		}
	}

	return nil
}
