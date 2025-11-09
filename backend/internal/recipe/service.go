package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
)

var (
	//ErrRecipeNotFound is returned when a recipe is not found
	ErrRecipeNotFound = errors.New("recipe not found")
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
	// ErrUnauthorized is returned when a user tries to modify a recipe they don't own
	ErrUnauthorized = errors.New("unauthorized")
)

// Service defines the recipe business logic
type Service struct {
	repo Repository
}

// NewService creates a new recipe service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateRecipe creates a new recipe
func (s *Service) CreateRecipe(ctx context.Context, userID uuid.UUID, req *models.CreateRecipeRequest) (*models.Recipe, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	return s.repo.CreateRecipe(ctx, userID, req)
}

// GetRecipeByID retrieves a recipe by ID
func (s *Service) GetRecipeByID(ctx context.Context, id uuid.UUID) (*models.Recipe, error) {
	recipe, err := s.repo.GetRecipeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecipeNotFound, err)
	}
	return recipe, nil
}

// ListRecipes retrieves recipes with search/filter/sort
func (s *Service) ListRecipes(ctx context.Context, params *models.RecipeSearchParams) (*models.RecipeListResponse, error) {
	// Set default pagination
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	recipes, total, err := s.repo.ListRecipes(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.RecipeListResponse{
		Data: recipes,
		Meta: models.PaginationMeta{
			Total:  total,
			Limit:  int64(params.Limit),
			Offset: int64(params.Offset),
		},
	}, nil
}

// UpdateRecipe updates a recipe
func (s *Service) UpdateRecipe(ctx context.Context, userID, id uuid.UUID, req *models.UpdateRecipeRequest) (*models.Recipe, error) {
	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Check if recipe exists and user owns it
	existing, err := s.repo.GetRecipeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRecipeNotFound, err)
	}

	if existing.CreatedByUserID != userID {
		return nil, ErrUnauthorized
	}

	return s.repo.UpdateRecipe(ctx, id, req)
}

// DeleteRecipe deletes a recipe
func (s *Service) DeleteRecipe(ctx context.Context, userID, id uuid.UUID) error {
	// Check if recipe exists and user owns it
	existing, err := s.repo.GetRecipeByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRecipeNotFound, err)
	}

	if existing.CreatedByUserID != userID {
		return ErrUnauthorized
	}

	return s.repo.DeleteRecipe(ctx, id)
}

// GetRecipesByUserID retrieves recipes created by a user
func (s *Service) GetRecipesByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Recipe, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetRecipesByUserID(ctx, userID, limit, offset)
}

// validateCreateRequest validates the create recipe request
func (s *Service) validateCreateRequest(req *models.CreateRecipeRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if req.Servings <= 0 {
		return errors.New("servings must be greater than 0")
	}
	if req.ServingSize == "" {
		return errors.New("serving size is required")
	}
	if req.PrepTimeMinutes != nil && *req.PrepTimeMinutes < 0 {
		return errors.New("prep time cannot be negative")
	}
	if req.CookTimeMinutes != nil && *req.CookTimeMinutes < 0 {
		return errors.New("cook time cannot be negative")
	}
	if len(req.Ingredients) == 0 {
		return errors.New("at least one ingredient is required")
	}
	if len(req.Instructions) == 0 {
		return errors.New("at least one instruction is required")
	}

	// Validate ingredients
	for i, ing := range req.Ingredients {
		if ing.IngredientName == "" {
			return fmt.Errorf("ingredient %d: name is required", i)
		}
		if ing.Quantity <= 0 {
			return fmt.Errorf("ingredient %d: quantity must be greater than 0", i)
		}
		if ing.Unit == "" {
			return fmt.Errorf("ingredient %d: unit is required", i)
		}
	}

	// Validate instructions
	for i, inst := range req.Instructions {
		if inst.StepNumber <= 0 {
			return fmt.Errorf("instruction %d: step number must be greater than 0", i)
		}
		if inst.Instruction == "" {
			return fmt.Errorf("instruction %d: text is required", i)
		}
	}

	return nil
}

// validateUpdateRequest validates the update recipe request
func (s *Service) validateUpdateRequest(req *models.UpdateRecipeRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if req.Servings <= 0 {
		return errors.New("servings must be greater than 0")
	}
	if req.ServingSize == "" {
		return errors.New("serving size is required")
	}
	if req.PrepTimeMinutes != nil && *req.PrepTimeMinutes < 0 {
		return errors.New("prep time cannot be negative")
	}
	if req.CookTimeMinutes != nil && *req.CookTimeMinutes < 0 {
		return errors.New("cook time cannot be negative")
	}
	if len(req.Ingredients) == 0 {
		return errors.New("at least one ingredient is required")
	}
	if len(req.Instructions) == 0 {
		return errors.New("at least one instruction is required")
	}

	// Validate ingredients
	for i, ing := range req.Ingredients {
		if ing.IngredientName == "" {
			return fmt.Errorf("ingredient %d: name is required", i)
		}
		if ing.Quantity <= 0 {
			return fmt.Errorf("ingredient %d: quantity must be greater than 0", i)
		}
		if ing.Unit == "" {
			return fmt.Errorf("ingredient %d: unit is required", i)
		}
	}

	// Validate instructions
	for i, inst := range req.Instructions {
		if inst.StepNumber <= 0 {
			return fmt.Errorf("instruction %d: step number must be greater than 0", i)
		}
		if inst.Instruction == "" {
			return fmt.Errorf("instruction %d: text is required", i)
		}
	}

	return nil
}
