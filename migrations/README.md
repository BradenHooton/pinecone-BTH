# Database Migrations

## Overview

This directory serves as a placeholder and reference for database migrations.

The **actual migration files** will be created in the `pinecone-api` repository under `internal/store/migrations/` during **Epic 1: Foundation & Infrastructure Setup** (User Story 1.2: Database Setup & Migrations).

## Migration Files to be Created

The following migration files will be implemented in the API repository:

### 001_create_users.sql
Creates the `users` table with authentication fields.

**Up Migration:**
- Creates `uuid-ossp` extension
- Creates `users` table with columns:
  - id (UUID, primary key)
  - email (VARCHAR, unique)
  - password_hash (VARCHAR)
  - name (VARCHAR)
  - created_at, updated_at, deleted_at (TIMESTAMPTZ)
- Creates indexes on email and deleted_at

### 002_create_recipes.sql
Creates recipe-related tables.

**Up Migration:**
- Creates `grocery_department` ENUM type
- Creates `recipes` table
- Creates `recipe_ingredients` table
- Creates `recipe_instructions` table
- Creates `recipe_tags` table
- Creates all necessary indexes

### 003_create_meal_plans.sql
Creates meal planning tables.

**Up Migration:**
- Creates `meal_type` ENUM type
- Creates `meal_plans` table
- Creates `meal_plan_recipes` table with "Out of Kitchen" support
- Creates indexes

### 004_create_grocery_lists.sql
Creates grocery list tables.

**Up Migration:**
- Creates `grocery_item_status` ENUM type
- Creates `grocery_lists` table
- Creates `grocery_list_items` table
- Creates indexes

### 005_create_cookbooks.sql
Creates cookbook (recipe collection) tables.

**Up Migration:**
- Creates `cookbooks` table
- Creates `cookbook_recipes` junction table (many-to-many)
- Creates indexes

### 006_create_nutrition_cache.sql
Creates nutrition data caching table.

**Up Migration:**
- Creates `nutrition_cache` table for USDA API caching
- Creates indexes on usda_fdc_id and food_name

## Migration Tool

Migrations will be managed using [Goose](https://github.com/pressly/goose) with the following commands:

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run all migrations (up)
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# Check migration status
goose -dir internal/store/migrations postgres "$DATABASE_URL" status

# Rollback last migration
goose -dir internal/store/migrations postgres "$DATABASE_URL" down

# Create new migration
goose -dir internal/store/migrations create migration_name sql
```

## Documentation Reference

For the complete database schema definition, see:
- **DATABASE_SCHEMA.md** - Full schema with DDL statements, constraints, and indexes
- **TDD.md** - Section 2 (Database Design)
- **EPIC_BREAKDOWN.md** - User Story 1.2 (Database Setup & Migrations)

## Implementation Timeline

**Epic**: 1 (Foundation & Infrastructure)
**User Story**: 1.2 (Database Setup & Migrations)
**Estimated Effort**: 12 hours
**Target Completion**: Week 1-2 of development

## Notes

- All migrations use UUID primary keys for security
- Soft deletes implemented via `deleted_at` timestamps
- Foreign keys use `ON DELETE CASCADE` or `ON DELETE SET NULL` appropriately
- Computed columns (e.g., `total_time_minutes`) are database-generated
- All ENUMs are defined before tables that reference them

## Migration Best Practices

1. **Always test migrations** on development database first
2. **Write down migrations** for every up migration
3. **Never modify existing migrations** - create new ones instead
4. **Back up production data** before running migrations
5. **Run migrations sequentially** in numbered order
6. **Document complex migrations** with comments in SQL
7. **Test rollback procedures** before deploying

---

**Status**: Migrations pending - will be created in `pinecone-api` repository during Epic 1
