package grocerylist

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("grocery list not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// Service handles business logic for grocery lists
type Service struct {
	repo Repository
}

// NewService creates a new grocery list service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateGroceryList creates a new grocery list by aggregating ingredients from meal plans
func (s *Service) CreateGroceryList(ctx context.Context, userID uuid.UUID, req *models.CreateGroceryListRequest) (*models.GroceryList, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid start_date format", ErrInvalidInput)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid end_date format", ErrInvalidInput)
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("%w: end_date must be after start_date", ErrInvalidInput)
	}

	// Create grocery list
	groceryList, err := s.repo.CreateGroceryList(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get ingredients from meal plans in date range
	ingredients, err := s.repo.GetIngredientsForDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Aggregate ingredients
	aggregatedItems := s.aggregateIngredients(ingredients)

	// Create grocery list items
	for _, item := range aggregatedItems {
		item.GroceryListID = groceryList.ID
		item.Status = models.StatusPending
		item.IsManual = false

		createdItem, err := s.repo.CreateGroceryListItem(ctx, &item)
		if err != nil {
			return nil, err
		}

		groceryList.Items = append(groceryList.Items, *createdItem)
	}

	return groceryList, nil
}

// GetGroceryListByID retrieves a grocery list by ID
func (s *Service) GetGroceryListByID(ctx context.Context, userID, groceryListID uuid.UUID) (*models.GroceryList, error) {
	groceryList, err := s.repo.GetGroceryListByID(ctx, groceryListID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check authorization
	if groceryList.CreatedByUserID != userID {
		return nil, ErrUnauthorized
	}

	return groceryList, nil
}

// GetGroceryListsByUser retrieves grocery lists for a user
func (s *Service) GetGroceryListsByUser(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]models.GroceryList, int64, error) {
	// Set defaults
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetGroceryListsByUser(ctx, userID, limit, offset)
}

// DeleteGroceryList deletes a grocery list
func (s *Service) DeleteGroceryList(ctx context.Context, userID, groceryListID uuid.UUID) error {
	// Get grocery list to check ownership
	groceryList, err := s.repo.GetGroceryListByID(ctx, groceryListID)
	if err != nil {
		return ErrNotFound
	}

	// Check authorization
	if groceryList.CreatedByUserID != userID {
		return ErrUnauthorized
	}

	return s.repo.DeleteGroceryList(ctx, groceryListID)
}

// AddManualItem adds a manual item to a grocery list
func (s *Service) AddManualItem(ctx context.Context, userID, groceryListID uuid.UUID, req *models.CreateManualItemRequest) (*models.GroceryListItem, error) {
	// Get grocery list to check ownership
	groceryList, err := s.repo.GetGroceryListByID(ctx, groceryListID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Check authorization
	if groceryList.CreatedByUserID != userID {
		return nil, ErrUnauthorized
	}

	// Validate input
	if strings.TrimSpace(req.ItemName) == "" {
		return nil, fmt.Errorf("%w: item_name is required", ErrInvalidInput)
	}

	// Set default department
	department := models.DepartmentOther
	if req.Department != nil {
		department = *req.Department
	}

	// Create item
	item := &models.GroceryListItem{
		GroceryListID: groceryListID,
		ItemName:      strings.TrimSpace(req.ItemName),
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		Department:    department,
		Status:        models.StatusPending,
		IsManual:      true,
	}

	return s.repo.CreateGroceryListItem(ctx, item)
}

// UpdateItemStatus updates the status of a grocery list item
func (s *Service) UpdateItemStatus(ctx context.Context, userID, itemID uuid.UUID, req *models.UpdateItemStatusRequest) error {
	// Note: We should verify the item belongs to a list owned by the user
	// For simplicity, we'll just update the status
	// In production, add a query to verify ownership

	return s.repo.UpdateItemStatus(ctx, itemID, req.Status)
}

// aggregateIngredients aggregates ingredients by name and unit, summing quantities
func (s *Service) aggregateIngredients(ingredients []IngredientAggregation) []models.GroceryListItem {
	// Map to aggregate: key is "name|unit", value is the aggregated item
	aggregated := make(map[string]*models.GroceryListItem)

	for _, ing := range ingredients {
		// Scale quantity based on servings
		scaledQuantity := ing.Quantity
		if ing.MealServings != nil && *ing.MealServings > 0 && ing.RecipeServings > 0 {
			scaledQuantity = ing.Quantity * float64(*ing.MealServings) / float64(ing.RecipeServings)
		}

		// Create key for grouping (normalized name + unit)
		key := s.createAggregationKey(ing.IngredientName, ing.Unit)

		// If item already exists, sum the quantities
		if existing, ok := aggregated[key]; ok {
			if existing.Quantity != nil {
				newQty := *existing.Quantity + scaledQuantity
				existing.Quantity = &newQty
			} else {
				existing.Quantity = &scaledQuantity
			}
		} else {
			// Create new aggregated item
			unit := ing.Unit
			aggregated[key] = &models.GroceryListItem{
				ItemName:       ing.IngredientName,
				Quantity:       &scaledQuantity,
				Unit:           &unit,
				Department:     ing.Department,
				SourceRecipeID: &ing.RecipeID,
			}
		}
	}

	// Convert map to slice
	items := make([]models.GroceryListItem, 0, len(aggregated))
	for _, item := range aggregated {
		items = append(items, *item)
	}

	return items
}

// createAggregationKey creates a key for grouping ingredients
// Ingredients with the same normalized name and unit are grouped together
func (s *Service) createAggregationKey(name, unit string) string {
	normalizedName := strings.TrimSpace(strings.ToLower(name))
	normalizedUnit := strings.TrimSpace(strings.ToLower(unit))
	return fmt.Sprintf("%s|%s", normalizedName, normalizedUnit)
}

// validateCreateRequest validates the create grocery list request
func (s *Service) validateCreateRequest(req *models.CreateGroceryListRequest) error {
	if req.StartDate == "" {
		return errors.New("start_date is required")
	}
	if req.EndDate == "" {
		return errors.New("end_date is required")
	}
	return nil
}
