-- +goose Up
-- Create nutrition_cache table
CREATE TABLE nutrition_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usda_fdc_id VARCHAR(50) NOT NULL UNIQUE,
    food_name VARCHAR(200) NOT NULL,
    calories DECIMAL(10, 2),
    protein_g DECIMAL(10, 2),
    carbs_g DECIMAL(10, 2),
    fiber_g DECIMAL(10, 2),
    fat_g DECIMAL(10, 2),
    cached_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_nutrition_cache_usda_id ON nutrition_cache(usda_fdc_id);
CREATE INDEX idx_nutrition_cache_food_name ON nutrition_cache(food_name);

-- Add foreign key to recipe_ingredients (referencing nutrition_cache)
ALTER TABLE recipe_ingredients
    ADD CONSTRAINT fk_recipe_ingredients_nutrition
    FOREIGN KEY (nutrition_id) REFERENCES nutrition_cache(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE recipe_ingredients DROP CONSTRAINT IF EXISTS fk_recipe_ingredients_nutrition;
DROP INDEX IF EXISTS idx_nutrition_cache_food_name;
DROP INDEX IF EXISTS idx_nutrition_cache_usda_id;
DROP TABLE IF EXISTS nutrition_cache;
