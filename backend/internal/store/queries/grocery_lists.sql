-- Grocery List Queries

-- Create a new grocery list
-- name: CreateGroceryList :one
INSERT INTO grocery_lists (
    created_by_user_id, start_date, end_date
) VALUES ($1, $2, $3)
RETURNING id, created_by_user_id, start_date, end_date, created_at, updated_at;

-- Get grocery list by ID
-- name: GetGroceryListByID :one
SELECT id, created_by_user_id, start_date, end_date, created_at, updated_at
FROM grocery_lists
WHERE id = $1;

-- Get grocery lists for a user
-- name: GetGroceryListsByUser :many
SELECT id, created_by_user_id, start_date, end_date, created_at, updated_at
FROM grocery_lists
WHERE created_by_user_id = $1
ORDER BY start_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- Count grocery lists for a user
-- name: CountGroceryListsByUser :one
SELECT COUNT(*)
FROM grocery_lists
WHERE created_by_user_id = $1;

-- Delete grocery list
-- name: DeleteGroceryList :exec
DELETE FROM grocery_lists
WHERE id = $1;

-- Create grocery list item
-- name: CreateGroceryListItem :one
INSERT INTO grocery_list_items (
    grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id;

-- Get items for a grocery list
-- name: GetGroceryListItems :many
SELECT id, grocery_list_id, item_name, quantity, unit, department, status, is_manual, source_recipe_id
FROM grocery_list_items
WHERE grocery_list_id = $1
ORDER BY department, item_name;

-- Update item status
-- name: UpdateItemStatus :exec
UPDATE grocery_list_items
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- Delete grocery list item
-- name: DeleteGroceryListItem :exec
DELETE FROM grocery_list_items
WHERE id = $1;

-- Get recipe ingredients for aggregation
-- This query gets all ingredients from recipes in meal plans within a date range
-- Used by the service to aggregate quantities
-- name: GetIngredientsForDateRange :many
SELECT
    ri.ingredient_name,
    ri.quantity,
    ri.unit,
    ri.department,
    r.id as recipe_id,
    r.title as recipe_title,
    mpr.servings as meal_servings,
    r.servings as recipe_servings
FROM meal_plans mp
JOIN meal_plan_recipes mpr ON mp.id = mpr.meal_plan_id
JOIN recipes r ON mpr.recipe_id = r.id
JOIN recipe_ingredients ri ON r.id = ri.recipe_id
WHERE mp.plan_date >= $1
  AND mp.plan_date <= $2
  AND mpr.out_of_kitchen = false
  AND r.deleted_at IS NULL
ORDER BY ri.department, ri.ingredient_name;
