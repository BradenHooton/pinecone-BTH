-- +goose Up
-- Create cookbooks table
CREATE TABLE cookbooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create cookbook_recipes junction table
CREATE TABLE cookbook_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cookbook_id UUID NOT NULL REFERENCES cookbooks(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cookbook_id, recipe_id)
);

-- Create indexes
CREATE INDEX idx_cookbooks_user_id ON cookbooks(created_by_user_id);
CREATE INDEX idx_cookbooks_deleted_at ON cookbooks(deleted_at);
CREATE INDEX idx_cookbook_recipes_cookbook_id ON cookbook_recipes(cookbook_id);
CREATE INDEX idx_cookbook_recipes_recipe_id ON cookbook_recipes(recipe_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cookbook_recipes_recipe_id;
DROP INDEX IF EXISTS idx_cookbook_recipes_cookbook_id;
DROP INDEX IF EXISTS idx_cookbooks_deleted_at;
DROP INDEX IF EXISTS idx_cookbooks_user_id;
DROP TABLE IF EXISTS cookbook_recipes;
DROP TABLE IF EXISTS cookbooks;
