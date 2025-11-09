-- Meal Plans

-- name: CreateMealPlan :one
INSERT INTO meal_plans (plan_date)
VALUES ($1)
RETURNING *;

-- name: GetMealPlanByDate :one
SELECT * FROM meal_plans
WHERE plan_date = $1
LIMIT 1;

-- name: GetMealPlanByID :one
SELECT * FROM meal_plans
WHERE id = $1
LIMIT 1;

-- name: GetMealPlansByDateRange :many
SELECT * FROM meal_plans
WHERE plan_date >= $1 AND plan_date <= $2
ORDER BY plan_date ASC;

-- name: GetOrCreateMealPlan :one
INSERT INTO meal_plans (plan_date)
VALUES ($1)
ON CONFLICT (plan_date) DO UPDATE SET updated_at = NOW()
RETURNING *;

-- Meal Plan Recipes

-- name: CreateMealPlanRecipe :one
INSERT INTO meal_plan_recipes (
    meal_plan_id,
    meal_type,
    recipe_id,
    servings,
    out_of_kitchen,
    order_index
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMealPlanRecipes :many
SELECT * FROM meal_plan_recipes
WHERE meal_plan_id = $1
ORDER BY meal_type ASC, order_index ASC;

-- name: GetMealPlanRecipesByDateRange :many
SELECT mpr.* FROM meal_plan_recipes mpr
JOIN meal_plans mp ON mpr.meal_plan_id = mp.id
WHERE mp.plan_date >= $1 AND mp.plan_date <= $2
ORDER BY mp.plan_date ASC, mpr.meal_type ASC, mpr.order_index ASC;

-- name: DeleteMealPlanRecipes :exec
DELETE FROM meal_plan_recipes
WHERE meal_plan_id = $1;

-- name: DeleteMealPlanRecipesByMealType :exec
DELETE FROM meal_plan_recipes
WHERE meal_plan_id = $1 AND meal_type = $2;
