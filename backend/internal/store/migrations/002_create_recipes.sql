-- +goose Up
-- Create grocery_department ENUM
CREATE TYPE grocery_department AS ENUM (
    'produce',
    'meat',
    'seafood',
    'dairy',
    'bakery',
    'frozen',
    'pantry',
    'spices',
    'beverages',
    'other'
);

-- Create recipes table
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    image_url TEXT,
    servings INT NOT NULL CHECK (servings > 0),
    serving_size VARCHAR(100) NOT NULL,
    prep_time_minutes INT CHECK (prep_time_minutes >= 0),
    cook_time_minutes INT CHECK (cook_time_minutes >= 0),
    total_time_minutes INT GENERATED ALWAYS AS (
        COALESCE(prep_time_minutes, 0) + COALESCE(cook_time_minutes, 0)
    ) STORED,
    storage_notes TEXT,
    source VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create recipe_ingredients table
CREATE TABLE recipe_ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    nutrition_id UUID,
    ingredient_name VARCHAR(200) NOT NULL,
    quantity DECIMAL(10, 3) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    department grocery_department NOT NULL DEFAULT 'other',
    order_index INT NOT NULL
);

-- Create recipe_instructions table
CREATE TABLE recipe_instructions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number INT NOT NULL CHECK (step_number > 0),
    instruction TEXT NOT NULL
);

-- Create recipe_tags table
CREATE TABLE recipe_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    tag_name VARCHAR(50) NOT NULL
);

-- Create indexes for recipes
CREATE INDEX idx_recipes_title ON recipes(title) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipes_created_by ON recipes(created_by_user_id);
CREATE INDEX idx_recipes_deleted_at ON recipes(deleted_at);

-- Create indexes for recipe_ingredients
CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_name ON recipe_ingredients(ingredient_name);

-- Create indexes for recipe_instructions
CREATE INDEX idx_recipe_instructions_recipe_id ON recipe_instructions(recipe_id);

-- Create indexes for recipe_tags
CREATE INDEX idx_recipe_tags_recipe_id ON recipe_tags(recipe_id);
CREATE INDEX idx_recipe_tags_tag_name ON recipe_tags(tag_name);

-- +goose Down
DROP INDEX IF EXISTS idx_recipe_tags_tag_name;
DROP INDEX IF EXISTS idx_recipe_tags_recipe_id;
DROP INDEX IF EXISTS idx_recipe_instructions_recipe_id;
DROP INDEX IF EXISTS idx_recipe_ingredients_name;
DROP INDEX IF EXISTS idx_recipe_ingredients_recipe_id;
DROP INDEX IF EXISTS idx_recipes_deleted_at;
DROP INDEX IF EXISTS idx_recipes_created_by;
DROP INDEX IF EXISTS idx_recipes_title;
DROP TABLE IF EXISTS recipe_tags;
DROP TABLE IF EXISTS recipe_instructions;
DROP TABLE IF EXISTS recipe_ingredients;
DROP TABLE IF EXISTS recipes;
DROP TYPE IF EXISTS grocery_department;
