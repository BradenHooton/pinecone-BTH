-- name: CreateRecipe :one
INSERT INTO recipes (
    created_by_user_id,
    title,
    image_url,
    servings,
    serving_size,
    prep_time_minutes,
    cook_time_minutes,
    storage_notes,
    source,
    notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetRecipeByID :one
SELECT * FROM recipes
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListRecipes :many
SELECT * FROM recipes
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateRecipe :one
UPDATE recipes
SET title = $2,
    image_url = $3,
    servings = $4,
    serving_size = $5,
    prep_time_minutes = $6,
    cook_time_minutes = $7,
    storage_notes = $8,
    source = $9,
    notes = $10,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteRecipe :exec
UPDATE recipes
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: SearchRecipesByTitle :many
SELECT * FROM recipes
WHERE deleted_at IS NULL
  AND title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Recipe Ingredients

-- name: CreateRecipeIngredient :one
INSERT INTO recipe_ingredients (
    recipe_id,
    nutrition_id,
    ingredient_name,
    quantity,
    unit,
    department,
    order_index
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetRecipeIngredients :many
SELECT * FROM recipe_ingredients
WHERE recipe_id = $1
ORDER BY order_index ASC;

-- name: DeleteRecipeIngredients :exec
DELETE FROM recipe_ingredients
WHERE recipe_id = $1;

-- Recipe Instructions

-- name: CreateRecipeInstruction :one
INSERT INTO recipe_instructions (
    recipe_id,
    step_number,
    instruction
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRecipeInstructions :many
SELECT * FROM recipe_instructions
WHERE recipe_id = $1
ORDER BY step_number ASC;

-- name: DeleteRecipeInstructions :exec
DELETE FROM recipe_instructions
WHERE recipe_id = $1;

-- Recipe Tags

-- name: CreateRecipeTag :one
INSERT INTO recipe_tags (recipe_id, tag_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetRecipeTags :many
SELECT * FROM recipe_tags
WHERE recipe_id = $1;

-- name: DeleteRecipeTags :exec
DELETE FROM recipe_tags
WHERE recipe_id = $1;

-- name: SearchRecipesByTag :many
SELECT DISTINCT r.* FROM recipes r
JOIN recipe_tags rt ON r.id = rt.recipe_id
WHERE r.deleted_at IS NULL
  AND rt.tag_name = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchRecipesByIngredient :many
SELECT DISTINCT r.* FROM recipes r
JOIN recipe_ingredients ri ON r.id = ri.recipe_id
WHERE r.deleted_at IS NULL
  AND ri.ingredient_name ILIKE '%' || $1 || '%'
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountRecipes :one
SELECT COUNT(*) FROM recipes
WHERE deleted_at IS NULL;

-- name: GetRecipesByUserID :many
SELECT * FROM recipes
WHERE created_by_user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
