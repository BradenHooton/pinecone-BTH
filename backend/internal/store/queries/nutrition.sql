-- name: CreateNutritionCache :one
INSERT INTO nutrition_cache (
    usda_fdc_id,
    food_name,
    calories,
    protein_g,
    carbs_g,
    fiber_g,
    fat_g
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetNutritionByFdcID :one
SELECT * FROM nutrition_cache
WHERE usda_fdc_id = $1
LIMIT 1;

-- name: GetNutritionByID :one
SELECT * FROM nutrition_cache
WHERE id = $1
LIMIT 1;

-- name: SearchNutritionByName :many
SELECT * FROM nutrition_cache
WHERE food_name ILIKE '%' || $1 || '%'
ORDER BY food_name ASC
LIMIT $2;

-- name: DeleteOldNutritionCache :exec
DELETE FROM nutrition_cache
WHERE cached_at < NOW() - INTERVAL '90 days';
