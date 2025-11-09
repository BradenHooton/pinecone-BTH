# Database Schema Documentation
## Pinecone Recipe Management System

**Version:** 1.0  
**Date:** 2025-11-09  
**Database:** PostgreSQL 16

---

## Table of Contents
1. [Overview](#overview)
2. [Entity Relationship Diagram](#entity-relationship-diagram)
3. [Table Definitions](#table-definitions)
4. [Indexes](#indexes)
5. [Enums](#enums)
6. [Constraints](#constraints)
7. [Migration Files](#migration-files)

---

## Overview

### Design Principles

1. **UUIDs for Primary Keys:** Prevents enumeration attacks, enables distributed ID generation
2. **Soft Deletes:** `deleted_at` timestamp allows data recovery and audit trails
3. **Computed Columns:** `total_time_minutes` auto-calculated (prep + cook)
4. **ENUMs for Fixed Lists:** `meal_type`, `grocery_department`, `grocery_item_status`
5. **Normalization:** Ingredients, instructions, tags in separate tables (1:N)
6. **Many-to-Many:** `cookbook_recipes` junction table for flexible organization

### Database Statistics (Expected)

| Metric | Estimate |
|--------|----------|
| **Users** | 2-6 |
| **Recipes** | 100-1,000 |
| **Meal Plans** | ~365 per year |
| **Grocery Lists** | ~52 per year |
| **Cookbooks** | 10-20 |
| **Total Storage** | < 1GB |

---

## Entity Relationship Diagram

```
users (1) ──────────┬─────────── (N) recipes
                    │
                    ├─────────── (N) cookbooks
                    │
                    └─────────── (N) grocery_lists

recipes (1) ────────┬─────────── (N) recipe_ingredients
                    │
                    ├─────────── (N) recipe_instructions
                    │
                    ├─────────── (N) recipe_tags
                    │
                    ├─────────── (N) meal_plan_recipes
                    │
                    └─────────── (N) cookbook_recipes

nutrition_cache (1) ─────────── (N) recipe_ingredients

meal_plans (1) ───────────────── (N) meal_plan_recipes

cookbooks (1) ────────────────── (N) cookbook_recipes

grocery_lists (1) ────────────── (N) grocery_list_items
```

---

## Table Definitions

### users

**Purpose:** Stores household user accounts

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique user identifier |
| email | VARCHAR(255) | NOT NULL, UNIQUE | User email address |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| name | VARCHAR(100) | NOT NULL | User display name |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Account creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMPTZ | NULL | Soft delete timestamp |

**DDL:**
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

---

### recipes

**Purpose:** Stores recipe master data

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique recipe identifier |
| created_by_user_id | UUID | FK → users.id, NOT NULL | Recipe creator |
| title | VARCHAR(200) | NOT NULL | Recipe title |
| image_url | TEXT | NULL | URL or path to recipe image |
| servings | INT | NOT NULL, CHECK > 0 | Number of servings |
| serving_size | VARCHAR(100) | NOT NULL | Serving size description (e.g., "1 cup") |
| prep_time_minutes | INT | NULL, CHECK ≥ 0 | Preparation time in minutes |
| cook_time_minutes | INT | NULL, CHECK ≥ 0 | Cooking time in minutes |
| total_time_minutes | INT | GENERATED | Auto-calculated: prep + cook |
| storage_notes | TEXT | NULL | Storage and freezing instructions |
| source | VARCHAR(500) | NULL | Recipe source (e.g., URL, book) |
| notes | TEXT | NULL | Additional notes |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Recipe creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMPTZ | NULL | Soft delete timestamp |

**DDL:**
```sql
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

CREATE INDEX idx_recipes_title ON recipes(title) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipes_created_by ON recipes(created_by_user_id);
CREATE INDEX idx_recipes_deleted_at ON recipes(deleted_at);
```

---

### recipe_ingredients

**Purpose:** Stores ingredients for each recipe

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique ingredient identifier |
| recipe_id | UUID | FK → recipes.id, NOT NULL | Parent recipe |
| nutrition_id | UUID | FK → nutrition_cache.id, NULL | Linked nutrition data |
| ingredient_name | VARCHAR(200) | NOT NULL | Ingredient name (e.g., "chicken breast") |
| quantity | DECIMAL(10,3) | NOT NULL, CHECK > 0 | Quantity amount |
| unit | VARCHAR(50) | NOT NULL | Measurement unit (e.g., "cups", "lbs") |
| department | grocery_department | NOT NULL, DEFAULT 'other' | Grocery store department |
| order_index | INT | NOT NULL | Display order in recipe |

**DDL:**
```sql
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

CREATE TABLE recipe_ingredients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    nutrition_id UUID REFERENCES nutrition_cache(id),
    ingredient_name VARCHAR(200) NOT NULL,
    quantity DECIMAL(10, 3) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    department grocery_department NOT NULL DEFAULT 'other',
    order_index INT NOT NULL
);

CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_name ON recipe_ingredients(ingredient_name);
```

---

### recipe_instructions

**Purpose:** Stores step-by-step cooking instructions

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique instruction identifier |
| recipe_id | UUID | FK → recipes.id, NOT NULL | Parent recipe |
| step_number | INT | NOT NULL, CHECK > 0 | Step order number |
| instruction | TEXT | NOT NULL | Instruction text |

**DDL:**
```sql
CREATE TABLE recipe_instructions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number INT NOT NULL CHECK (step_number > 0),
    instruction TEXT NOT NULL
);

CREATE INDEX idx_recipe_instructions_recipe_id ON recipe_instructions(recipe_id);
```

---

### recipe_tags

**Purpose:** Stores tags for categorizing recipes

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique tag identifier |
| recipe_id | UUID | FK → recipes.id, NOT NULL | Parent recipe |
| tag_name | VARCHAR(50) | NOT NULL | Tag name (e.g., "vegetarian", "quick") |

**DDL:**
```sql
CREATE TABLE recipe_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    tag_name VARCHAR(50) NOT NULL
);

CREATE INDEX idx_recipe_tags_recipe_id ON recipe_tags(recipe_id);
CREATE INDEX idx_recipe_tags_tag_name ON recipe_tags(tag_name);
```

---

### nutrition_cache

**Purpose:** Caches USDA FoodData Central API nutrition data

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique cache entry identifier |
| usda_fdc_id | VARCHAR(50) | NOT NULL, UNIQUE | USDA FDC ID |
| food_name | VARCHAR(200) | NOT NULL | Food name |
| calories | DECIMAL(10,2) | NULL | Calories per 100g |
| protein_g | DECIMAL(10,2) | NULL | Protein grams per 100g |
| carbs_g | DECIMAL(10,2) | NULL | Carbohydrate grams per 100g |
| fiber_g | DECIMAL(10,2) | NULL | Fiber grams per 100g |
| fat_g | DECIMAL(10,2) | NULL | Fat grams per 100g |
| cached_at | TIMESTAMPTZ | DEFAULT NOW() | Cache timestamp |

**DDL:**
```sql
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

CREATE INDEX idx_nutrition_cache_usda_id ON nutrition_cache(usda_fdc_id);
CREATE INDEX idx_nutrition_cache_food_name ON nutrition_cache(food_name);
```

---

### meal_plans

**Purpose:** Stores meal plan dates (one row per day)

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique meal plan identifier |
| plan_date | DATE | NOT NULL, UNIQUE | Date of meal plan |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |

**DDL:**
```sql
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'snack', 'dinner', 'dessert');

CREATE TABLE meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_date DATE NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_meal_plans_date ON meal_plans(plan_date);
```

---

### meal_plan_recipes

**Purpose:** Links recipes to meal plan slots

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique entry identifier |
| meal_plan_id | UUID | FK → meal_plans.id, NOT NULL | Parent meal plan |
| meal_type | meal_type | NOT NULL | Type of meal |
| recipe_id | UUID | FK → recipes.id, NULL | Linked recipe |
| servings | INT | CHECK > 0 OR NULL | Number of servings |
| out_of_kitchen | BOOLEAN | DEFAULT FALSE | If TRUE, household is eating out |
| order_index | INT | NOT NULL, DEFAULT 0 | Display order (for multiple recipes) |

**Constraint:** Either `out_of_kitchen=TRUE` with `recipe_id=NULL`, or `out_of_kitchen=FALSE` with `recipe_id` and `servings` set.

**DDL:**
```sql
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

CREATE INDEX idx_meal_plan_recipes_meal_plan_id ON meal_plan_recipes(meal_plan_id);
CREATE INDEX idx_meal_plan_recipes_recipe_id ON meal_plan_recipes(recipe_id);
```

---

### grocery_lists

**Purpose:** Stores generated grocery lists

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique grocery list identifier |
| created_by_user_id | UUID | FK → users.id, NOT NULL | User who generated list |
| start_date | DATE | NOT NULL | Start of meal plan range |
| end_date | DATE | NOT NULL | End of meal plan range |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Generation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |

**Constraint:** `end_date >= start_date`

**DDL:**
```sql
CREATE TYPE grocery_item_status AS ENUM ('pending', 'bought', 'have_on_hand');

CREATE TABLE grocery_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_date_range CHECK (end_date >= start_date)
);

CREATE INDEX idx_grocery_lists_user_id ON grocery_lists(created_by_user_id);
CREATE INDEX idx_grocery_lists_dates ON grocery_lists(start_date, end_date);
```

---

### grocery_list_items

**Purpose:** Stores individual items in grocery lists

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique item identifier |
| grocery_list_id | UUID | FK → grocery_lists.id, NOT NULL | Parent grocery list |
| item_name | VARCHAR(200) | NOT NULL | Item name |
| quantity | DECIMAL(10,3) | NULL | Quantity (NULL for manual items) |
| unit | VARCHAR(50) | NULL | Unit of measurement |
| department | grocery_department | NOT NULL, DEFAULT 'other' | Grocery store department |
| status | grocery_item_status | NOT NULL, DEFAULT 'pending' | Purchase status |
| is_manual | BOOLEAN | DEFAULT FALSE | TRUE if user-added item |
| source_recipe_id | UUID | FK → recipes.id, NULL | Origin recipe (for traceability) |

**DDL:**
```sql
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

CREATE INDEX idx_grocery_list_items_list_id ON grocery_list_items(grocery_list_id);
CREATE INDEX idx_grocery_list_items_department ON grocery_list_items(department);
CREATE INDEX idx_grocery_list_items_status ON grocery_list_items(status);
```

---

### cookbooks

**Purpose:** Stores recipe collections (e.g., "Holiday Recipes")

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique cookbook identifier |
| created_by_user_id | UUID | FK → users.id, NOT NULL | Cookbook creator |
| name | VARCHAR(200) | NOT NULL | Cookbook name |
| description | TEXT | NULL | Optional description |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMPTZ | NULL | Soft delete timestamp |

**DDL:**
```sql
CREATE TABLE cookbooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_cookbooks_user_id ON cookbooks(created_by_user_id);
CREATE INDEX idx_cookbooks_deleted_at ON cookbooks(deleted_at);
```

---

### cookbook_recipes

**Purpose:** Many-to-many junction table for cookbooks and recipes

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique entry identifier |
| cookbook_id | UUID | FK → cookbooks.id, NOT NULL | Parent cookbook |
| recipe_id | UUID | FK → recipes.id, NOT NULL | Linked recipe |
| added_at | TIMESTAMPTZ | DEFAULT NOW() | When recipe was added |

**Constraint:** UNIQUE(cookbook_id, recipe_id)

**DDL:**
```sql
CREATE TABLE cookbook_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cookbook_id UUID NOT NULL REFERENCES cookbooks(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cookbook_id, recipe_id)
);

CREATE INDEX idx_cookbook_recipes_cookbook_id ON cookbook_recipes(cookbook_id);
CREATE INDEX idx_cookbook_recipes_recipe_id ON cookbook_recipes(recipe_id);
```

---

## Indexes

### Purpose

Indexes optimize query performance for frequently accessed data patterns.

### Index List

| Table | Index Name | Columns | Purpose |
|-------|-----------|---------|---------|
| users | idx_users_email | email | Login queries |
| users | idx_users_deleted_at | deleted_at | Soft delete filtering |
| recipes | idx_recipes_title | title | Search by title |
| recipes | idx_recipes_created_by | created_by_user_id | User's recipes |
| recipes | idx_recipes_deleted_at | deleted_at | Soft delete filtering |
| recipe_ingredients | idx_recipe_ingredients_recipe_id | recipe_id | Fetch ingredients for recipe |
| recipe_ingredients | idx_recipe_ingredients_name | ingredient_name | Search by ingredient |
| recipe_instructions | idx_recipe_instructions_recipe_id | recipe_id | Fetch instructions for recipe |
| recipe_tags | idx_recipe_tags_recipe_id | recipe_id | Fetch tags for recipe |
| recipe_tags | idx_recipe_tags_tag_name | tag_name | Filter recipes by tag |
| nutrition_cache | idx_nutrition_cache_usda_id | usda_fdc_id | USDA lookup |
| nutrition_cache | idx_nutrition_cache_food_name | food_name | Search nutrition |
| meal_plans | idx_meal_plans_date | plan_date | Date range queries |
| meal_plan_recipes | idx_meal_plan_recipes_meal_plan_id | meal_plan_id | Fetch meals for plan |
| meal_plan_recipes | idx_meal_plan_recipes_recipe_id | recipe_id | Find meal plans using recipe |
| grocery_lists | idx_grocery_lists_user_id | created_by_user_id | User's grocery lists |
| grocery_lists | idx_grocery_lists_dates | start_date, end_date | Date range queries |
| grocery_list_items | idx_grocery_list_items_list_id | grocery_list_id | Fetch items for list |
| grocery_list_items | idx_grocery_list_items_department | department | Group by department |
| grocery_list_items | idx_grocery_list_items_status | status | Filter by status |
| cookbooks | idx_cookbooks_user_id | created_by_user_id | User's cookbooks |
| cookbooks | idx_cookbooks_deleted_at | deleted_at | Soft delete filtering |
| cookbook_recipes | idx_cookbook_recipes_cookbook_id | cookbook_id | Fetch recipes in cookbook |
| cookbook_recipes | idx_cookbook_recipes_recipe_id | recipe_id | Find cookbooks with recipe |

---

## Enums

### grocery_department

**Values:**
- `produce` - Fresh fruits and vegetables
- `meat` - Fresh and packaged meats
- `seafood` - Fresh and frozen seafood
- `dairy` - Milk, cheese, yogurt, eggs
- `bakery` - Bread, buns, pastries
- `frozen` - Frozen meals, vegetables, desserts
- `pantry` - Pasta, rice, canned goods
- `spices` - Herbs, spices, condiments
- `beverages` - Drinks, juices, coffee, tea
- `other` - Miscellaneous items

### meal_type

**Values:**
- `breakfast` - Morning meal
- `lunch` - Midday meal
- `snack` - Between-meal snack
- `dinner` - Evening meal
- `dessert` - After-dinner dessert

### grocery_item_status

**Values:**
- `pending` - Not yet purchased
- `bought` - Purchased at store
- `have_on_hand` - Already have at home

---

## Constraints

### Check Constraints

| Table | Constraint | Logic |
|-------|-----------|-------|
| recipes | servings > 0 | Must have at least 1 serving |
| recipes | prep_time_minutes >= 0 | Cannot be negative |
| recipes | cook_time_minutes >= 0 | Cannot be negative |
| recipe_ingredients | quantity > 0 | Cannot have zero quantity |
| recipe_instructions | step_number > 0 | Steps start at 1 |
| meal_plan_recipes | valid_meal_entry | Either "Out of Kitchen" OR recipe + servings |
| meal_plan_recipes | servings > 0 OR NULL | If recipe present, servings must be positive |
| grocery_lists | end_date >= start_date | Valid date range |

### Foreign Key Constraints

All foreign keys use:
- **ON DELETE CASCADE:** Child records deleted when parent is deleted (recipes, ingredients, etc.)
- **ON DELETE SET NULL:** Orphaned records allowed (meal_plan_recipes.recipe_id, grocery_list_items.source_recipe_id)

### Unique Constraints

| Table | Columns | Purpose |
|-------|---------|---------|
| users | email | One account per email |
| meal_plans | plan_date | One plan per day |
| cookbook_recipes | cookbook_id, recipe_id | No duplicate recipe assignments |
| nutrition_cache | usda_fdc_id | One cache entry per USDA food |

---

## Migration Files

### Migration Order

1. `001_create_users.sql` - User accounts
2. `002_create_recipes.sql` - Recipes, ingredients, instructions, tags
3. `003_create_meal_plans.sql` - Meal plans and meal-recipe links
4. `004_create_grocery_lists.sql` - Grocery lists and items
5. `005_create_cookbooks.sql` - Cookbooks and recipe assignments
6. `006_create_nutrition_cache.sql` - USDA nutrition cache

### Running Migrations

```bash
# Apply all migrations
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# Check status
goose -dir internal/store/migrations postgres "$DATABASE_URL" status

# Rollback last migration
goose -dir internal/store/migrations postgres "$DATABASE_URL" down
```

---

## Backup & Recovery

### Daily Backups

```bash
# Backup script (run via cron)
#!/bin/bash
BACKUP_DIR="/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
pg_dump "$DATABASE_URL" > "$BACKUP_DIR/pinecone_$TIMESTAMP.sql"

# Keep last 7 days
find $BACKUP_DIR -name "pinecone_*.sql" -mtime +7 -delete
```

### Restore from Backup

```bash
# Restore database
psql "$DATABASE_URL" < /backups/pinecone_20250108_120000.sql
```

---

**Document Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-09 | GhostDev | Initial schema documentation |
