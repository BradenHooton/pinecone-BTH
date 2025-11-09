package cookbook

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("cookbook not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// Service handles business logic for cookbooks
type Service struct {
	repo Repository
}

// NewService creates a new cookbook service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateCookbook creates a new cookbook
func (s *Service) CreateCookbook(ctx context.Context, userID uuid.UUID, req *models.CreateCookbookRequest) (*models.Cookbook, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	return s.repo.CreateCookbook(ctx, userID, req)
}

// GetCookbookByID retrieves a cookbook by ID
func (s *Service) GetCookbookByID(ctx context.Context, userID, cookbookID uuid.UUID) (*models.Cookbook, error) {
	cookbook, err := s.repo.GetCookbookByID(ctx, cookbookID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check authorization
	if cookbook.CreatedByUserID != userID {
		return nil, ErrUnauthorized
	}

	return cookbook, nil
}

// GetCookbooksByUser retrieves cookbooks for a user
func (s *Service) GetCookbooksByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.Cookbook, int64, error) {
	// Set defaults
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetCookbooksByUser(ctx, userID, limit, offset)
}

// UpdateCookbook updates a cookbook
func (s *Service) UpdateCookbook(ctx context.Context, userID, cookbookID uuid.UUID, req *models.UpdateCookbookRequest) (*models.Cookbook, error) {
	// Validate input
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Get cookbook to check ownership
	cookbook, err := s.repo.GetCookbookByID(ctx, cookbookID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check authorization
	if cookbook.CreatedByUserID != userID {
		return nil, ErrUnauthorized
	}

	return s.repo.UpdateCookbook(ctx, cookbookID, req)
}

// DeleteCookbook deletes a cookbook
func (s *Service) DeleteCookbook(ctx context.Context, userID, cookbookID uuid.UUID) error {
	// Get cookbook to check ownership
	cookbook, err := s.repo.GetCookbookByID(ctx, cookbookID)
	if err != nil {
		return ErrNotFound
	}

	// Check authorization
	if cookbook.CreatedByUserID != userID {
		return ErrUnauthorized
	}

	return s.repo.DeleteCookbook(ctx, cookbookID)
}

// AddRecipeToCookbook adds a recipe to a cookbook
func (s *Service) AddRecipeToCookbook(ctx context.Context, userID, cookbookID, recipeID uuid.UUID) error {
	// Get cookbook to check ownership
	cookbook, err := s.repo.GetCookbookByID(ctx, cookbookID)
	if err != nil {
		return ErrNotFound
	}

	// Check authorization
	if cookbook.CreatedByUserID != userID {
		return ErrUnauthorized
	}

	return s.repo.AddRecipeToCookbook(ctx, cookbookID, recipeID)
}

// RemoveRecipeFromCookbook removes a recipe from a cookbook
func (s *Service) RemoveRecipeFromCookbook(ctx context.Context, userID, cookbookID, recipeID uuid.UUID) error {
	// Get cookbook to check ownership
	cookbook, err := s.repo.GetCookbookByID(ctx, cookbookID)
	if err != nil {
		return ErrNotFound
	}

	// Check authorization
	if cookbook.CreatedByUserID != userID {
		return ErrUnauthorized
	}

	return s.repo.RemoveRecipeFromCookbook(ctx, cookbookID, recipeID)
}

// validateCreateRequest validates the create cookbook request
func (s *Service) validateCreateRequest(req *models.CreateCookbookRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if len(req.Name) > 200 {
		return errors.New("name must be 200 characters or less")
	}

	return nil
}

// validateUpdateRequest validates the update cookbook request
func (s *Service) validateUpdateRequest(req *models.UpdateCookbookRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if len(req.Name) > 200 {
		return errors.New("name must be 200 characters or less")
	}

	return nil
}
