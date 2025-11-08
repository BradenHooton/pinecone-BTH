-- +goose Up
-- Create meal_type ENUM
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'snack', 'dinner', 'dessert');

-- Create meal_plans table
CREATE TABLE meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_date DATE NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create meal_plan_recipes table
CREATE TABLE meal_plan_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meal_plan_id UUID NOT NULL REFERENCES meal_plans(id) ON DELETE CASCADE,
    meal_type meal_type NOT NULL,
    recipe_id UUID REFERENCES recipes(id) ON DELETE SET NULL,
    servings INT CHECK (servings > 0 OR recipe_id IS NULL),
    out_of_kitchen BOOLEAN DEFAULT FALSE,
    order_index INT NOT NULL DEFAULT 0,
    CONSTRAINT valid_meal_entry CHECK (
        (out_of_kitchen = TRUE AND recipe_id IS NULL) OR
        (out_of_kitchen = FALSE AND recipe_id IS NOT NULL AND servings IS NOT NULL)
    )
);

-- Create indexes
CREATE INDEX idx_meal_plans_date ON meal_plans(plan_date);
CREATE INDEX idx_meal_plan_recipes_meal_plan_id ON meal_plan_recipes(meal_plan_id);
CREATE INDEX idx_meal_plan_recipes_recipe_id ON meal_plan_recipes(recipe_id);

-- +goose Down
DROP INDEX IF EXISTS idx_meal_plan_recipes_recipe_id;
DROP INDEX IF EXISTS idx_meal_plan_recipes_meal_plan_id;
DROP INDEX IF EXISTS idx_meal_plans_date;
DROP TABLE IF EXISTS meal_plan_recipes;
DROP TABLE IF EXISTS meal_plans;
DROP TYPE IF EXISTS meal_type;
