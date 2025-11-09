package nutrition

import (
	"context"
	"fmt"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for nutrition data access
type Repository interface {
	CreateNutritionCache(ctx context.Context, fdcID, foodName string, calories, proteinG, carbsG, fiberG, fatG *float64) (*models.NutritionCache, error)
	GetNutritionByFdcID(ctx context.Context, fdcID string) (*models.NutritionCache, error)
	GetNutritionByID(ctx context.Context, id uuid.UUID) (*models.NutritionCache, error)
	SearchNutritionByName(ctx context.Context, query string, limit int) ([]models.NutritionCache, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateNutritionCache creates a new nutrition cache entry
func (r *PostgresRepository) CreateNutritionCache(ctx context.Context, fdcID, foodName string, calories, proteinG, carbsG, fiberG, fatG *float64) (*models.NutritionCache, error) {
	var cache models.NutritionCache
	err := r.db.QueryRow(ctx, `
		INSERT INTO nutrition_cache (
			usda_fdc_id, food_name, calories, protein_g, carbs_g, fiber_g, fat_g
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, usda_fdc_id, food_name, calories, protein_g, carbs_g, fiber_g, fat_g, cached_at
	`, fdcID, foodName, calories, proteinG, carbsG, fiberG, fatG).
		Scan(&cache.ID, &cache.USDAFdcID, &cache.FoodName, &cache.Calories,
			&cache.ProteinG, &cache.CarbsG, &cache.FiberG, &cache.FatG, &cache.CachedAt)
	if err != nil {
		return nil, fmt.Errorf("create nutrition cache: %w", err)
	}
	return &cache, nil
}

// GetNutritionByFdcID retrieves nutrition data by USDA FDC ID
func (r *PostgresRepository) GetNutritionByFdcID(ctx context.Context, fdcID string) (*models.NutritionCache, error) {
	var cache models.NutritionCache
	err := r.db.QueryRow(ctx, `
		SELECT id, usda_fdc_id, food_name, calories, protein_g, carbs_g, fiber_g, fat_g, cached_at
		FROM nutrition_cache
		WHERE usda_fdc_id = $1
	`, fdcID).Scan(&cache.ID, &cache.USDAFdcID, &cache.FoodName, &cache.Calories,
		&cache.ProteinG, &cache.CarbsG, &cache.FiberG, &cache.FatG, &cache.CachedAt)
	if err != nil {
		return nil, fmt.Errorf("get nutrition by fdc id: %w", err)
	}
	return &cache, nil
}

// GetNutritionByID retrieves nutrition data by ID
func (r *PostgresRepository) GetNutritionByID(ctx context.Context, id uuid.UUID) (*models.NutritionCache, error) {
	var cache models.NutritionCache
	err := r.db.QueryRow(ctx, `
		SELECT id, usda_fdc_id, food_name, calories, protein_g, carbs_g, fiber_g, fat_g, cached_at
		FROM nutrition_cache
		WHERE id = $1
	`, id).Scan(&cache.ID, &cache.USDAFdcID, &cache.FoodName, &cache.Calories,
		&cache.ProteinG, &cache.CarbsG, &cache.FiberG, &cache.FatG, &cache.CachedAt)
	if err != nil {
		return nil, fmt.Errorf("get nutrition by id: %w", err)
	}
	return &cache, nil
}

// SearchNutritionByName searches cached nutrition data by food name
func (r *PostgresRepository) SearchNutritionByName(ctx context.Context, query string, limit int) ([]models.NutritionCache, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, usda_fdc_id, food_name, calories, protein_g, carbs_g, fiber_g, fat_g, cached_at
		FROM nutrition_cache
		WHERE food_name ILIKE '%' || $1 || '%'
		ORDER BY food_name ASC
		LIMIT $2
	`, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search nutrition by name: %w", err)
	}
	defer rows.Close()

	var results []models.NutritionCache
	for rows.Next() {
		var cache models.NutritionCache
		if err := rows.Scan(&cache.ID, &cache.USDAFdcID, &cache.FoodName, &cache.Calories,
			&cache.ProteinG, &cache.CarbsG, &cache.FiberG, &cache.FatG, &cache.CachedAt); err != nil {
			return nil, fmt.Errorf("scan nutrition: %w", err)
		}
		results = append(results, cache)
	}

	return results, nil
}
