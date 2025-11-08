-- +goose Up
-- Create grocery_item_status ENUM
CREATE TYPE grocery_item_status AS ENUM ('pending', 'bought', 'have_on_hand');

-- Create grocery_lists table
CREATE TABLE grocery_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_date_range CHECK (end_date >= start_date)
);

-- Create grocery_list_items table
CREATE TABLE grocery_list_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grocery_list_id UUID NOT NULL REFERENCES grocery_lists(id) ON DELETE CASCADE,
    item_name VARCHAR(200) NOT NULL,
    quantity DECIMAL(10, 3),
    unit VARCHAR(50),
    department grocery_department NOT NULL DEFAULT 'other',
    status grocery_item_status NOT NULL DEFAULT 'pending',
    is_manual BOOLEAN DEFAULT FALSE,
    source_recipe_id UUID REFERENCES recipes(id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_grocery_lists_user_id ON grocery_lists(created_by_user_id);
CREATE INDEX idx_grocery_lists_dates ON grocery_lists(start_date, end_date);
CREATE INDEX idx_grocery_list_items_list_id ON grocery_list_items(grocery_list_id);
CREATE INDEX idx_grocery_list_items_department ON grocery_list_items(department);
CREATE INDEX idx_grocery_list_items_status ON grocery_list_items(status);

-- +goose Down
DROP INDEX IF EXISTS idx_grocery_list_items_status;
DROP INDEX IF EXISTS idx_grocery_list_items_department;
DROP INDEX IF EXISTS idx_grocery_list_items_list_id;
DROP INDEX IF EXISTS idx_grocery_lists_dates;
DROP INDEX IF EXISTS idx_grocery_lists_user_id;
DROP TABLE IF EXISTS grocery_list_items;
DROP TABLE IF EXISTS grocery_lists;
DROP TYPE IF EXISTS grocery_item_status;
