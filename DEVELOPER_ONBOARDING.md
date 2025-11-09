# Developer Onboarding Guide
## Pinecone Recipe Management System

**Version:** 1.0  
**Date:** 2025-11-09

---

## Table of Contents
1. [Welcome](#welcome)
2. [Prerequisites](#prerequisites)
3. [Local Development Setup](#local-development-setup)
4. [Development Workflow](#development-workflow)
5. [Testing Standards](#testing-standards)
6. [Code Style & Conventions](#code-style--conventions)
7. [Pull Request Process](#pull-request-process)
8. [Common Tasks](#common-tasks)
9. [Troubleshooting](#troubleshooting)

---

## Welcome

Welcome to the Pinecone project! This guide will help you set up your local development environment and understand our development practices.

### Project Overview

Pinecone is a single-household recipe management and meal planning application built with:
- **Backend:** Go, chi router, PostgreSQL, sqlc
- **Frontend:** React, TypeScript, Vite, TanStack Query/Router
- **Infrastructure:** Docker, Caddy, GitHub Actions

### Team Philosophy

1. **Test-Driven Development (TDD):** Write tests before implementation
2. **Clean Architecture:** Clear separation of concerns (handler → service → repository)
3. **Security First:** JWT authentication, bcrypt passwords, parameterized queries
4. **Documentation:** Code is read more than written

---

## Prerequisites

### Required Software

| Tool | Version | Installation |
|------|---------|--------------|
| **Go** | 1.21+ | https://go.dev/doc/install |
| **Node.js** | 20+ | https://nodejs.org/ |
| **Docker** | 24+ | https://docs.docker.com/get-docker/ |
| **Docker Compose** | 2.20+ | Included with Docker Desktop |
| **Git** | 2.40+ | https://git-scm.com/downloads |
| **Make** (optional) | Any | Pre-installed on macOS/Linux, Windows: http://gnuwin32.sourceforge.net/packages/make.htm |

### Optional Tools

- **DBeaver** or **pgAdmin** - PostgreSQL GUI client
- **Postman** or **Insomnia** - API testing
- **VS Code** - Recommended IDE with Go and TypeScript extensions

### Verify Installation

```bash
go version       # Should show 1.21 or higher
node --version   # Should show v20 or higher
docker --version # Should show 24 or higher
git --version    # Should show 2.40 or higher
```

---

## Local Development Setup

### Step 1: Clone Repositories

```bash
# Create project directory
mkdir ~/pinecone-project
cd ~/pinecone-project

# Clone repositories
git clone https://github.com/bradenhooton/pinecone-api.git
git clone https://github.com/bradenhooton/pinecone-web.git
```

### Step 2: Backend Setup

```bash
cd pinecone-api

# Install Go dependencies
go mod download

# Install development tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Copy environment template
cp .env.example .env.dev

# Edit .env.dev with your settings
# Required: DATABASE_URL, JWT_SECRET, USDA_API_KEY
```

**`.env.dev` Example:**
```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pinecone_dev?sslmode=disable
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRY_HOURS=24
USDA_API_KEY=your-usda-api-key-from-https://fdc.nal.usda.gov/api-key-signup.html
USDA_API_BASE_URL=https://api.nal.usda.gov/fdc/v1
SERVER_PORT=8080
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE_MB=5
ALLOWED_ORIGINS=http://localhost:5173
LOG_LEVEL=debug
```

### Step 3: Database Setup

```bash
# Start PostgreSQL via Docker Compose
docker-compose up -d

# Wait for database to be ready
sleep 5

# Run migrations
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# Verify migrations
goose -dir internal/store/migrations postgres "$DATABASE_URL" status

# Optional: Seed test data
psql "$DATABASE_URL" < internal/store/migrations/seed_test_data.sql
```

### Step 4: Generate sqlc Code

```bash
# Generate Go code from SQL queries
sqlc generate

# Verify generation
ls internal/store/queries/*.sql.go
```

### Step 5: Run Backend

```bash
# Run the server
go run cmd/server/main.go

# Or use Makefile (if available)
make run

# Server should start on http://localhost:8080
# Test health endpoint: curl http://localhost:8080/health
```

### Step 6: Frontend Setup

```bash
cd ../pinecone-web

# Install dependencies
npm install

# Copy environment template
cp .env.example .env.dev

# Edit .env.dev
# VITE_API_BASE_URL=http://localhost:8080/api/v1

# Generate API types from OpenAPI spec
npm run generate-types

# This runs: npx openapi-typescript ../pinecone-api/api/openapi.yaml -o src/lib/api-types.ts
```

### Step 7: Run Frontend

```bash
# Start Vite dev server
npm run dev

# Frontend should start on http://localhost:5173
# Open browser to http://localhost:5173
```

### Step 8: Verify Everything Works

1. **Backend Health Check:**
   ```bash
   curl http://localhost:8080/health
   # Expected: {"status": "ok"}
   ```

2. **Frontend:** Navigate to `http://localhost:5173`
3. **Database:** Connect with DBeaver to `localhost:5432`, database `pinecone_dev`

---

## Development Workflow

### Git Branching Strategy

We follow **GitHub Flow:**

1. `main` branch is always deployable (production-ready)
2. All work happens on feature branches
3. Feature branches merge via Pull Request

### Branch Naming Convention

```
feat/[story-id]/[short-description]    # New features
fix/[bug-id]/[short-description]       # Bug fixes
chore/[task-description]               # Tooling, refactoring
docs/[documentation-update]            # Documentation only
```

**Examples:**
```
feat/auth-101/add-login-endpoint
fix/bug-202/login-button-crash
chore/update-dependencies
docs/api-documentation-update
```

### Daily Workflow

1. **Pull latest changes:**
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create feature branch:**
   ```bash
   git checkout -b feat/recipe-123/add-search
   ```

3. **Write test (RED):**
   ```bash
   # Create test file: internal/recipe/service_test.go
   # Write failing test
   go test ./internal/recipe
   # Test should fail
   ```

4. **Write implementation (GREEN):**
   ```bash
   # Write code in internal/recipe/service.go
   go test ./internal/recipe
   # Test should pass
   ```

5. **Refactor (REFACTOR):**
   ```bash
   # Improve code quality
   go test ./internal/recipe
   # Test should still pass
   ```

6. **Commit frequently:**
   ```bash
   git add .
   git commit -m "test: add search recipe by title test"
   git commit -m "feat: implement search recipe by title"
   git commit -m "refactor: extract search logic to helper"
   ```

7. **Push and open PR:**
   ```bash
   git push origin feat/recipe-123/add-search
   # Open Pull Request on GitHub
   ```

---

## Testing Standards

### Test-Driven Development (TDD)

**Mandatory:** All code must be test-driven.

**Process:**
1. **RED:** Write a failing test
2. **GREEN:** Write minimum code to pass
3. **REFACTOR:** Improve code quality

### Backend Testing

#### Unit Tests

```go
// File: internal/recipe/service_test.go
package recipe_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRecipeService_Create_ValidRecipe(t *testing.T) {
    // ARRANGE
    service := NewTestService()
    recipe := &models.Recipe{Title: "Test Recipe"}
    
    // ACT
    err := service.Create(context.Background(), recipe)
    
    // ASSERT
    assert.NoError(t, err)
    assert.NotNil(t, recipe.ID)
}
```

**Run tests:**
```bash
go test ./...                    # Run all tests
go test ./internal/recipe        # Run specific package
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage
go test -race ./...              # Race condition detection
```

#### Integration Tests

```go
// File: internal/recipe/repository_integration_test.go
// +build integration

func TestRecipeRepository_Create_Integration(t *testing.T) {
    // ARRANGE: Spin up real Postgres with testcontainers
    container := setupTestContainer(t)
    defer container.Terminate(context.Background())
    
    repo := NewRepository(container.DB)
    recipe := &models.Recipe{Title: "Integration Test"}
    
    // ACT
    err := repo.Create(context.Background(), recipe)
    
    // ASSERT
    assert.NoError(t, err)
    
    // Verify in database
    fetched, _ := repo.GetByID(context.Background(), recipe.ID)
    assert.Equal(t, "Integration Test", fetched.Title)
}
```

**Run integration tests:**
```bash
go test -tags=integration ./...
```

### Frontend Testing

#### Component Tests

```typescript
// File: src/components/recipe/RecipeCard.test.tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { RecipeCard } from './RecipeCard';

describe('RecipeCard', () => {
  it('should render recipe title', () => {
    // ARRANGE
    const recipe = { id: '1', title: 'Test Recipe' };
    
    // ACT
    render(<RecipeCard recipe={recipe} />);
    
    // ASSERT
    expect(screen.getByText('Test Recipe')).toBeInTheDocument();
  });
});
```

**Run tests:**
```bash
npm test                  # Run all tests
npm test -- RecipeCard    # Run specific file
npm run test:coverage     # With coverage
```

### Coverage Targets

| Layer | Target | Command |
|-------|--------|---------|
| Backend Services | ≥80% | `go test -cover ./internal/*/service.go` |
| Backend Repositories | ≥80% | `go test -cover ./internal/*/repository.go` |
| Frontend Components | ≥70% | `npm run test:coverage` |

---

## Code Style & Conventions

### Go Style

**Follow:** [Effective Go](https://go.dev/doc/effective_go)

**Key Rules:**
- Use `gofmt` (automatic with save in VS Code)
- Package names: lowercase, single word (`recipe`, not `recipeService`)
- Interfaces: `-er` suffix (`Creator`, `Fetcher`)
- Error handling: Always check errors, wrap with context
- Comments: Exported functions must have doc comments

**Example:**
```go
// CreateRecipe creates a new recipe in the database.
// It returns an error if the recipe title is empty or if the database operation fails.
func (s *Service) CreateRecipe(ctx context.Context, recipe *models.Recipe) error {
    if recipe.Title == "" {
        return errors.New("recipe title cannot be empty")
    }
    
    if err := s.repo.Create(ctx, recipe); err != nil {
        return errors.Wrap(err, "failed to create recipe")
    }
    
    return nil
}
```

### TypeScript Style

**Follow:** [TypeScript Best Practices](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html)

**Key Rules:**
- Use TypeScript strict mode
- Prefer `const` over `let`
- Use arrow functions for callbacks
- Explicit return types for functions
- No `any` types (use `unknown` if necessary)

**Example:**
```typescript
// Good
export const fetchRecipes = async (): Promise<Recipe[]> => {
  const response = await apiClient.get<Recipe[]>('/recipes');
  return response.data;
};

// Bad
export const fetchRecipes = async () => {
  const response = await apiClient.get('/recipes');
  return response.data;
};
```

### React Conventions

- **Functional components only** (no class components)
- **Hooks:** Prefix with `use` (`useRecipes`, `useAuth`)
- **File naming:** PascalCase for components (`RecipeCard.tsx`)
- **Folder structure:** Group by feature, not by type

**Example:**
```typescript
// src/components/recipe/RecipeCard.tsx
export const RecipeCard: React.FC<RecipeCardProps> = ({ recipe }) => {
  const { colors } = useDesignTokens();
  
  return (
    <div style={{ borderColor: colors.primary }}>
      <h3>{recipe.title}</h3>
      <p>{recipe.prepTime} min prep</p>
    </div>
  );
};
```

### SQL Conventions

- **Keywords:** UPPERCASE (`SELECT`, `FROM`, `WHERE`)
- **Table names:** snake_case, plural (`recipes`, `meal_plans`)
- **Column names:** snake_case (`created_at`, `prep_time_minutes`)
- **Comments:** Describe purpose of complex queries

**Example:**
```sql
-- name: GetRecipeByID :one
-- Fetches a single recipe by its UUID.
-- Returns an error if the recipe is soft-deleted.
SELECT id, title, created_at
FROM recipes
WHERE id = $1 AND deleted_at IS NULL;
```

---

## Pull Request Process

### Before Opening a PR

- [ ] All tests pass locally (`go test ./...` and `npm test`)
- [ ] Linter passes (`golangci-lint run` and `npm run lint`)
- [ ] Code coverage meets targets (≥80% backend, ≥70% frontend)
- [ ] Commit messages follow conventions (see below)
- [ ] Branch is up-to-date with `main`

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `test`: Adding tests
- `refactor`: Code refactoring
- `docs`: Documentation changes
- `chore`: Tooling, dependencies

**Example:**
```
feat(recipe): add search by ingredient

Implement search functionality that allows users to find recipes
by ingredient name. Uses case-insensitive LIKE query with sqlc.

Closes #123
```

### PR Template

When opening a PR, fill out this template:

```markdown
## Description
(What does this PR do? Why is it needed?)

## Related Task
- Fixes: #[Ticket-ID]

## Changes
- Created `/api/v1/recipes` endpoint
- Added `SearchRecipes` service method
- Updated OpenAPI spec

## How to Test
1. Check out this branch
2. Run `go run cmd/server/main.go`
3. Execute: `curl "http://localhost:8080/api/v1/recipes?search=chicken"`
4. **Expected Result:** JSON array of recipes containing "chicken"

## Checklist
- [ ] Tests added and passing
- [ ] Linter passes
- [ ] Documentation updated (if needed)
- [ ] OpenAPI spec updated (if API change)
```

### PR Review Checklist

Reviewers should verify:

- [ ] **Solves the Problem:** Code addresses the linked task
- [ ] **Clean Code:** Readable, maintainable, follows conventions
- [ ] **Tests Added:** Sufficient unit/integration tests
- [ ] **CI Checks Pass:** All automated checks green
- [ ] **Documentation:** New functionality documented
- [ ] **Security:** No obvious vulnerabilities (SQL injection, XSS, etc.)

### Approval Process

1. PR must have at least **1 approval** from team member
2. All CI checks must be **green**
3. All review comments must be **resolved**
4. Branch must be **up-to-date** with `main`

### Merging

- Use **"Squash and merge"** for feature branches
- Use **"Merge commit"** for hotfixes (preserves history)
- Delete branch after merge

---

## Common Tasks

### Add a New API Endpoint

1. **Update OpenAPI spec:**
   ```yaml
   # api/openapi.yaml
   /recipes/{id}:
     get:
       summary: Get recipe by ID
       # ... rest of spec
   ```

2. **Generate types (frontend):**
   ```bash
   cd pinecone-web
   npm run generate-types
   ```

3. **Write SQL query:**
   ```sql
   -- internal/store/queries/recipes.sql
   -- name: GetRecipeByID :one
   SELECT * FROM recipes WHERE id = $1;
   ```

4. **Generate sqlc code:**
   ```bash
   cd pinecone-api
   sqlc generate
   ```

5. **TDD: Write tests → Implement → Refactor**
   - `internal/recipe/service_test.go`
   - `internal/recipe/service.go`
   - `internal/recipe/handler_test.go`
   - `internal/recipe/handler.go`

6. **Register route:**
   ```go
   // cmd/server/main.go
   r.Get("/api/v1/recipes/{id}", recipeHandler.GetByID)
   ```

7. **Test manually:**
   ```bash
   curl http://localhost:8080/api/v1/recipes/some-uuid
   ```

### Run Database Migrations

```bash
# Check migration status
goose -dir internal/store/migrations postgres "$DATABASE_URL" status

# Run all pending migrations
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# Rollback last migration
goose -dir internal/store/migrations postgres "$DATABASE_URL" down

# Create new migration
goose -dir internal/store/migrations create add_new_table sql
```

### Reset Local Database

```bash
# Stop containers
docker-compose down

# Remove database volume
docker volume rm pinecone-api_postgres_data

# Start fresh
docker-compose up -d
goose -dir internal/store/migrations postgres "$DATABASE_URL" up
```

### Update Dependencies

**Backend:**
```bash
go get -u ./...                 # Update all
go get -u github.com/lib/pq     # Update specific
go mod tidy                     # Clean up
```

**Frontend:**
```bash
npm update                      # Update all
npm install react@latest        # Update specific
npm outdated                    # Check versions
```

### Run Linters

**Backend:**
```bash
golangci-lint run
golangci-lint run --fix  # Auto-fix issues
```

**Frontend:**
```bash
npm run lint
npm run lint:fix  # Auto-fix issues
```

### Generate Code Coverage Report

**Backend:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Frontend:**
```bash
npm run test:coverage
# Opens HTML report in browser
```

---

## Troubleshooting

### "Database connection refused"

**Problem:** Backend can't connect to PostgreSQL

**Solution:**
```bash
# Check if container is running
docker ps

# If not, start it
docker-compose up -d

# Check logs
docker-compose logs db

# Verify DATABASE_URL in .env.dev
echo $DATABASE_URL
```

### "Port 8080 already in use"

**Problem:** Another process is using port 8080

**Solution:**
```bash
# Find process using port
lsof -i :8080      # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process
kill -9 <PID>

# Or change SERVER_PORT in .env.dev
```

### "sqlc: command not found"

**Problem:** sqlc not installed

**Solution:**
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Add to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin
```

### "Module not found" errors in Go

**Problem:** Dependencies not downloaded

**Solution:**
```bash
go mod download
go mod tidy
```

### "npm ERR! peer dependency" errors

**Problem:** Dependency version conflicts

**Solution:**
```bash
rm -rf node_modules package-lock.json
npm install --legacy-peer-deps
```

### "CORS error" in browser

**Problem:** Frontend can't access backend API

**Solution:**
1. Check `ALLOWED_ORIGINS` in backend `.env.dev`
2. Ensure it includes `http://localhost:5173`
3. Restart backend server

### Tests failing with "database locked"

**Problem:** Parallel tests accessing same DB

**Solution:**
```bash
# Run tests sequentially
go test -p 1 ./...

# Or use testcontainers (spins up isolated DB per test)
```

---

## Getting Help

### Resources

- **Documentation:** `docs/` folder in repositories
- **API Reference:** `api/openapi.yaml`
- **BRD:** Business requirements and user stories
- **TDD:** Technical architecture and design decisions

### Contact

- **Project Owner:** BHooton
- **Questions:** Open GitHub issue with `question` label

---

## Quick Reference

### Backend Commands
```bash
go run cmd/server/main.go           # Start server
go test ./...                       # Run tests
go test -cover ./...                # Test with coverage
golangci-lint run                   # Lint
sqlc generate                       # Generate DB code
goose -dir internal/store/migrations postgres "$DATABASE_URL" up  # Migrate
```

### Frontend Commands
```bash
npm run dev                         # Start dev server
npm test                            # Run tests
npm run test:coverage               # Test with coverage
npm run lint                        # Lint
npm run build                       # Production build
npm run generate-types              # Generate API types
```

### Docker Commands
```bash
docker-compose up -d                # Start all services
docker-compose down                 # Stop all services
docker-compose logs -f api          # View logs
docker-compose exec db psql -U postgres pinecone_dev  # Access DB
```

---

**Welcome aboard! If you have questions, don't hesitate to ask. Happy coding! 🚀**
