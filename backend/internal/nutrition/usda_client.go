package nutrition

import (
	"context"
	"fmt"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

// USDAClient defines the interface for USDA FoodData Central API
type USDAClient interface {
	Search(ctx context.Context, query string) ([]models.NutritionSearchResult, error)
}

// StubUSDAClient is a stub implementation that returns mock data
// This will be replaced with real API calls when API key is available
type StubUSDAClient struct{}

// NewStubUSDAClient creates a new stub USDA client
func NewStubUSDAClient() *StubUSDAClient {
	return &StubUSDAClient{}
}

// Search returns mock nutrition data based on the query
func (c *StubUSDAClient) Search(ctx context.Context, query string) ([]models.NutritionSearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	// Mock database of common foods
	mockFoods := []models.NutritionSearchResult{
		{
			FdcID:       "123456",
			Description: "Chicken breast, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(120.0),
			ProteinG:    ptr(22.5),
			CarbsG:      ptr(0.0),
			FiberG:      ptr(0.0),
			FatG:        ptr(2.6),
		},
		{
			FdcID:       "123457",
			Description: "Rice, white, long-grain, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(365.0),
			ProteinG:    ptr(7.1),
			CarbsG:      ptr(80.0),
			FiberG:      ptr(1.3),
			FatG:        ptr(0.7),
		},
		{
			FdcID:       "123458",
			Description: "Broccoli, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(34.0),
			ProteinG:    ptr(2.8),
			CarbsG:      ptr(7.0),
			FiberG:      ptr(2.6),
			FatG:        ptr(0.4),
		},
		{
			FdcID:       "123459",
			Description: "Olive oil",
			DataType:    "SR Legacy",
			Calories:    ptr(884.0),
			ProteinG:    ptr(0.0),
			CarbsG:      ptr(0.0),
			FiberG:      ptr(0.0),
			FatG:        ptr(100.0),
		},
		{
			FdcID:       "123460",
			Description: "Milk, whole, 3.25% milkfat",
			DataType:    "SR Legacy",
			Calories:    ptr(61.0),
			ProteinG:    ptr(3.2),
			CarbsG:      ptr(4.8),
			FiberG:      ptr(0.0),
			FatG:        ptr(3.3),
		},
		{
			FdcID:       "123461",
			Description: "Egg, whole, raw, fresh",
			DataType:    "SR Legacy",
			Calories:    ptr(143.0),
			ProteinG:    ptr(12.6),
			CarbsG:      ptr(0.7),
			FiberG:      ptr(0.0),
			FatG:        ptr(9.5),
		},
		{
			FdcID:       "123462",
			Description: "Tomato, red, ripe, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(18.0),
			ProteinG:    ptr(0.9),
			CarbsG:      ptr(3.9),
			FiberG:      ptr(1.2),
			FatG:        ptr(0.2),
		},
		{
			FdcID:       "123463",
			Description: "Onion, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(40.0),
			ProteinG:    ptr(1.1),
			CarbsG:      ptr(9.3),
			FiberG:      ptr(1.7),
			FatG:        ptr(0.1),
		},
		{
			FdcID:       "123464",
			Description: "Garlic, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(149.0),
			ProteinG:    ptr(6.4),
			CarbsG:      ptr(33.1),
			FiberG:      ptr(2.1),
			FatG:        ptr(0.5),
		},
		{
			FdcID:       "123465",
			Description: "Pasta, dry, enriched",
			DataType:    "SR Legacy",
			Calories:    ptr(371.0),
			ProteinG:    ptr(13.0),
			CarbsG:      ptr(74.7),
			FiberG:      ptr(3.2),
			FatG:        ptr(1.5),
		},
		{
			FdcID:       "123466",
			Description: "Ground beef, 80% lean meat / 20% fat, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(254.0),
			ProteinG:    ptr(17.2),
			CarbsG:      ptr(0.0),
			FiberG:      ptr(0.0),
			FatG:        ptr(20.0),
		},
		{
			FdcID:       "123467",
			Description: "Salmon, Atlantic, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(142.0),
			ProteinG:    ptr(19.8),
			CarbsG:      ptr(0.0),
			FiberG:      ptr(0.0),
			FatG:        ptr(6.3),
		},
		{
			FdcID:       "123468",
			Description: "Potato, flesh and skin, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(77.0),
			ProteinG:    ptr(2.0),
			CarbsG:      ptr(17.5),
			FiberG:      ptr(2.1),
			FatG:        ptr(0.1),
		},
		{
			FdcID:       "123469",
			Description: "Carrot, raw",
			DataType:    "SR Legacy",
			Calories:    ptr(41.0),
			ProteinG:    ptr(0.9),
			CarbsG:      ptr(9.6),
			FiberG:      ptr(2.8),
			FatG:        ptr(0.2),
		},
		{
			FdcID:       "123470",
			Description: "Cheese, cheddar",
			DataType:    "SR Legacy",
			Calories:    ptr(403.0),
			ProteinG:    ptr(22.9),
			CarbsG:      ptr(3.1),
			FiberG:      ptr(0.0),
			FatG:        ptr(33.3),
		},
	}

	// Filter results based on query (case-insensitive)
	lowerQuery := strings.ToLower(query)
	var results []models.NutritionSearchResult

	for _, food := range mockFoods {
		if strings.Contains(strings.ToLower(food.Description), lowerQuery) {
			results = append(results, food)
		}
	}

	// Limit to 10 results
	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}

// Helper function to create pointer to float64
func ptr(f float64) *float64 {
	return &f
}
