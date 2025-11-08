# Epic Breakdown & Task Stories
## Pinecone Recipe Management System

**Version:** 1.0  
**Date:** 2025-11-08  
**Total Estimated Effort:** 422 hours (~10.5 weeks)

---

## Table of Contents
1. [Epic 1: Foundation & Infrastructure](#epic-1-foundation--infrastructure-52-hours)
2. [Epic 2: User Authentication](#epic-2-user-authentication-42-hours)
3. [Epic 3: Recipe Management](#epic-3-recipe-management-76-hours)
4. [Epic 4: Nutrition Data Integration](#epic-4-nutrition-data-integration-32-hours)
5. [Epic 5: Meal Planning](#epic-5-meal-planning-40-hours)
6. [Epic 6: Grocery List Generation](#epic-6-grocery-list-generation-38-hours)
7. [Epic 7: Ingredient-Based Menu Recommendation](#epic-7-ingredient-based-menu-recommendation-24-hours)
8. [Epic 8: Cookbooks](#epic-8-cookbooks-36-hours)
9. [Epic 9: Polish, Deployment & Production](#epic-9-polish-deployment--production-82-hours)

---

## EPIC 1: Foundation & Infrastructure (52 hours)

**Duration:** 2 weeks  
**Priority:** P0 (Blocking)

### User Story 1.1: Repository Setup & Structure (4h)

**As a** developer  
**I want** standardized repository structure  
**So that** the codebase is organized and follows best practices

**Acceptance Criteria:**
- [ ] `pinecone-api` repository created with Standard Go Project Layout
- [ ] `pinecone-web` repository created with feature-based structure
- [ ] All repos have `.gitignore` files
- [ ] All repos have `README.md` with setup instructions

**Tasks:**
1. Create GitHub repositories (3x)
2. Initialize Go module: `go mod init github.com/bhooton/pinecone-api`
3. Initialize Vite project: `npm create vite@latest pinecone-web -- --template react-ts`
4. Create directory structures per TDD
5. Add `.gitignore` files
6. Create initial `README.md` files
7. Push initial commits

---

### User Story 1.2: Database Setup & Migrations (12h)

**As a** developer  
**I want** PostgreSQL with versioned migrations  
**So that** schema changes are tracked

**Acceptance Criteria:**
- [ ] PostgreSQL 16 running via Docker Compose
- [ ] Goose installed and configured
- [ ] All 6 migrations created and applied
- [ ] `sqlc.yaml` configuration created
- [ ] Test data seeding script created

**Tasks:**
1. Write `docker-compose.yml` for local development
2. Install Goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`
3. Create `internal/store/migrations/` directory
4. Write migration `001_create_users.sql`
5. Write migration `002_create_recipes.sql`
6. Write migration `003_create_meal_plans.sql`
7. Write migration `004_create_grocery_lists.sql`
8. Write migration `005_create_cookbooks.sql`
9. Write migration `006_create_nutrition_cache.sql`
10. Run migrations: `goose -dir internal/store/migrations postgres $DATABASE_URL up`
11. Verify schema
12. Create `sqlc.yaml`
13. Write seed data script

---

### User Story 1.3: OpenAPI Specification (8h)

**As a** developer  
**I want** complete OpenAPI 3.0 specification  
**So that** API contracts are documented

**Acceptance Criteria:**
- [ ] `api/openapi.yaml` created
- [ ] All endpoints documented
- [ ] Security scheme defined
- [ ] Spec validates without errors

**Tasks:**
1. Create `api/openapi.yaml`
2. Define `info`, `servers`, `components.securitySchemes`
3. Define all schemas (User, Recipe, MealPlan, etc.)
4. Document all endpoints
5. Validate: `npx @apidevtools/swagger-cli validate api/openapi.yaml`

---

### User Story 1.4: Configuration Management (6h)

**As a** developer  
**I want** centralized configuration  
**So that** environment settings are managed securely

**Acceptance Criteria:**
- [ ] `.env.dev` file created
- [ ] `config/grocery_departments.yaml` created
- [ ] `internal/config/config.go` loads `.env` with Viper
- [ ] Configuration validated on startup

**Tasks:**
1. Create `.env.dev` template
2. Install Viper: `go get github.com/spf13/viper`
3. Install godotenv: `go get github.com/joho/godotenv`
4. Write `internal/config/config.go`
5. **TDD:** Write `internal/config/config_test.go`
6. Create `config/grocery_departments.yaml`
7. Write `internal/config/departments.go`
8. **TDD:** Write `internal/config/departments_test.go`
9. Add config validation in `cmd/server/main.go`

---

### User Story 1.5: JWT Middleware & Auth Utilities (8h)

**As a** developer  
**I want** reusable JWT utilities  
**So that** protected routes are secured

**Acceptance Criteria:**
- [ ] `pkg/jwt/jwt.go` generates and validates JWT tokens
- [ ] JWT tokens in HTTP-only cookies (24h expiry)
- [ ] `internal/middleware/auth.go` validates JWT
- [ ] Unit tests with ≥80% coverage

**Tasks:**
1. **TDD:** Write `pkg/jwt/jwt_test.go` first
   - Test: Generate token for user ID
   - Test: Validate valid token
   - Test: Reject expired token
   - Test: Reject tampered token
2. Install JWT: `go get github.com/golang-jwt/jwt/v5`
3. Write `pkg/jwt/jwt.go` (pass tests)
4. **TDD:** Write `internal/middleware/auth_test.go`
   - Test: Allow valid JWT
   - Test: Reject missing cookie
   - Test: Reject invalid JWT
   - Test: Extract user ID to context
5. Write `internal/middleware/auth.go`
6. Run coverage: `go test -cover ./pkg/jwt ./internal/middleware`

---

### User Story 1.6: Logging & Request ID Middleware (4h)

**As a** developer  
**I want** structured JSON logging with request tracing  
**So that** all API calls are debuggable

**Acceptance Criteria:**
- [ ] `log/slog` configured for JSON output
- [ ] Request ID middleware generates unique ID
- [ ] Request ID in all log entries
- [ ] Logger middleware logs: method, path, status, duration

**Tasks:**
1. Write `internal/middleware/request_id.go`
2. Write `internal/middleware/logger.go`
3. Configure `slog.NewJSONHandler()` in `cmd/server/main.go`
4. Test manually: `curl http://localhost:8080/health`

---

### User Story 1.7: CI Pipeline (6h)

**As a** developer  
**I want** automated testing on every PR  
**So that** code quality is enforced

**Acceptance Criteria:**
- [ ] `.github/workflows/ci.yml` created for backend
- [ ] CI runs on every PR to `main`
- [ ] Steps: Checkout → Setup Go → Lint → Test → Build
- [ ] `.github/workflows/ci-frontend.yml` created

**Tasks:**
1. Create `.github/workflows/ci.yml` (backend)
2. Add golangci-lint step
3. Add test step with coverage
4. Add build step
5. Create `.github/workflows/ci-frontend.yml`
6. Add ESLint step
7. Add Vitest step
8. Add Vite build step
9. Verify CI runs

---

### User Story 1.8: Design Tokens (4h)

**As a** developer  
**I want** centralized design token system  
**So that** branding is consistent

**Acceptance Criteria:**
- [ ] `src/tokens/colors.ts` defines forest green, warm off-white
- [ ] `src/tokens/typography.ts` defines fonts
- [ ] `src/tokens/spacing.ts` defines spacing scale
- [ ] `src/styles/global.css` applies tokens

**Tasks:**
1. Create `src/tokens/colors.ts`
2. Create `src/tokens/typography.ts`
3. Create `src/tokens/spacing.ts`
4. Create `src/styles/global.css` with CSS variables
5. Import tokens in components

---

## EPIC 2: User Authentication (42 hours)

**Duration:** 1 week  
**Priority:** P0 (Blocking)

### User Story 2.1: User Registration Backend (10h)

**As a** new user  
**I want** to register with email and password  
**So that** I can create an account

**Acceptance Criteria:**
- [ ] `POST /api/v1/auth/register` endpoint created
- [ ] Email validation (RFC 5322)
- [ ] Password validation (min 8 chars)
- [ ] Password hashed with bcrypt (cost 12)
- [ ] Duplicate email returns `409 Conflict`
- [ ] Test coverage ≥80%

**Tasks:**
1. **TDD:** Write `internal/auth/service_test.go`
   - Test: Register valid user
   - Test: Reject duplicate email
   - Test: Reject weak password
   - Test: Hash password
2. Write `internal/auth/service.go`
3. **TDD:** Write `internal/auth/repository_test.go`
4. Create sqlc query: `internal/store/queries/users.sql`
5. Run `sqlc generate`
6. Write `internal/auth/repository.go`
7. **TDD:** Write `internal/auth/handler_test.go`
8. Write `internal/auth/handler.go`
9. Register route in `cmd/server/main.go`
10. Integration test with testcontainers

---

### User Story 2.2: User Login Backend (8h)

**As a** registered user  
**I want** to login with email and password  
**So that** I can access protected features

**Acceptance Criteria:**
- [ ] `POST /api/v1/auth/login` endpoint created
- [ ] Password verified with bcrypt
- [ ] Success returns `200 OK` + JWT cookie
- [ ] Invalid credentials return `401 Unauthorized`
- [ ] Test coverage ≥80%

**Tasks:**
1. **TDD:** Write `internal/auth/service_test.go` (login tests)
2. Update `internal/auth/service.go` with `Login()`
3. **TDD:** Write `internal/auth/handler_test.go` (login tests)
4. Update `internal/auth/handler.go` with `Login()`
5. Register route
6. Test manually

---

### User Story 2.3: User Logout Backend (4h)

**As a** logged-in user  
**I want** to logout  
**So that** my session is securely ended

**Acceptance Criteria:**
- [ ] `POST /api/v1/auth/logout` endpoint created
- [ ] Success returns `204 No Content` + clears cookie
- [ ] Cookie cleared with `MaxAge=-1`

**Tasks:**
1. **TDD:** Write `internal/auth/handler_test.go` (logout test)
2. Write `internal/auth/handler.go` `Logout()`
3. Register protected route
4. Test manually

---

### User Story 2.4: Registration Form Frontend (8h)

**As a** new user  
**I want** a registration form  
**So that** I can create an account

**Acceptance Criteria:**
- [ ] `/register` route created
- [ ] Form fields: email, password, name
- [ ] Client-side validation with Zod
- [ ] API call to `POST /api/v1/auth/register`
- [ ] Success redirects to `/login`
- [ ] Component test passes

**Tasks:**
1. Install dependencies
2. Generate API types: `npx openapi-typescript`
3. Create `src/routes/register.tsx`
4. Create Zod schema
5. Build form with React Hook Form
6. Create `src/hooks/useAuth.ts`
7. **TDD:** Write `src/routes/register.test.tsx`
8. Implement component
9. Style with design tokens

---

### User Story 2.5: Login Form Frontend (8h)

**As a** registered user  
**I want** a login form  
**So that** I can authenticate

**Acceptance Criteria:**
- [ ] `/login` route created
- [ ] Form fields: email, password
- [ ] Client-side validation with Zod
- [ ] Success stores auth state + redirects
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/login.tsx`
2. Create Zod schema
3. Build form
4. Update `src/hooks/useAuth.ts`
5. Create `src/store/authStore.ts` (Zustand)
6. **TDD:** Write `src/routes/login.test.tsx`
7. Implement component
8. Style

---

### User Story 2.6: Protected Route Wrapper Frontend (4h)

**As a** developer  
**I want** route guard for authenticated pages  
**So that** unauthenticated users are redirected

**Acceptance Criteria:**
- [ ] `src/components/common/ProtectedRoute.tsx` created
- [ ] Checks auth state from Zustand
- [ ] Redirects to `/login` if not authenticated
- [ ] Test passes

**Tasks:**
1. Create `src/components/common/ProtectedRoute.tsx`
2. Read auth state
3. Conditional render/redirect
4. **TDD:** Write test
5. Wrap routes

---

## EPIC 3: Recipe Management (76 hours)

**Duration:** 2 weeks  
**Priority:** P1

### User Story 3.1: Recipe CRUD Endpoints Backend (20h)

**As a** user  
**I want** to create, read, update, delete recipes  
**So that** I can manage my recipe database

**Acceptance Criteria:**
- [ ] `POST /api/v1/recipes` creates recipe
- [ ] `GET /api/v1/recipes` lists recipes (pagination)
- [ ] `GET /api/v1/recipes/{id}` fetches single recipe
- [ ] `PUT /api/v1/recipes/{id}` updates recipe
- [ ] `DELETE /api/v1/recipes/{id}` soft-deletes
- [ ] All endpoints require authentication
- [ ] Test coverage ≥80%

**Tasks:**
1. Create sqlc queries: `internal/store/queries/recipes.sql`
2. Run `sqlc generate`
3. **TDD:** Write `internal/recipe/service_test.go`
4. Write `internal/recipe/service.go`
5. **TDD:** Write `internal/recipe/repository_test.go`
6. Write `internal/recipe/repository.go`
7. **TDD:** Write `internal/recipe/handler_test.go`
8. Write `internal/recipe/handler.go`
9. Register routes
10. Integration test

---

### User Story 3.2: Recipe Image Upload Backend (6h)

**As a** user  
**I want** to upload recipe images  
**So that** I can visually identify recipes

**Acceptance Criteria:**
- [ ] `POST /api/v1/recipes/upload-image` endpoint
- [ ] Accepts `multipart/form-data`
- [ ] Validates file type (jpg, jpeg, png, webp)
- [ ] Validates size (≤5MB)
- [ ] Saves to `uploads/` with UUID filename
- [ ] Returns `201 Created` + `image_url`

**Tasks:**
1. **TDD:** Write `internal/recipe/handler_test.go` (image tests)
2. Write `internal/recipe/handler.go` `UploadImage()`
3. Generate UUID filename
4. Configure Caddy to serve `/uploads`
5. Test manually

---

### User Story 3.3: Recipe Search, Filter, Sort Backend (10h)

**As a** user  
**I want** to search and filter recipes  
**So that** I can quickly find recipes

**Acceptance Criteria:**
- [ ] `GET /api/v1/recipes?search=chicken` searches title/ingredients
- [ ] `GET /api/v1/recipes?tags=vegetarian` filters by tags
- [ ] `GET /api/v1/recipes?sort=title_asc` sorts results
- [ ] Pagination with `limit` and `offset`
- [ ] Test coverage ≥80%

**Tasks:**
1. Update sqlc query with search logic
2. Run `sqlc generate`
3. **TDD:** Write service tests
4. Update service with `SearchRecipes()`
5. **TDD:** Write handler tests
6. Update handler
7. Test manually

---

### User Story 3.4: Recipe Grid UI Frontend (10h)

**As a** user  
**I want** to view recipes in a grid  
**So that** I can browse my collection

**Acceptance Criteria:**
- [ ] `/recipes` route displays grid
- [ ] Each card shows: image, title, prep/cook time
- [ ] Responsive grid (1-3+ columns)
- [ ] Click card navigates to detail
- [ ] Loading state with skeletons
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/recipes/index.tsx`
2. Create `src/hooks/useRecipes.ts`
3. **TDD:** Write `src/components/recipe/RecipeCard.test.tsx`
4. Create `RecipeCard.tsx`
5. **TDD:** Write `RecipeGrid.test.tsx`
6. Create `RecipeGrid.tsx`
7. Style with CSS Grid

---

### User Story 3.5: Recipe Detail Page Frontend (10h)

**As a** user  
**I want** to view full recipe details  
**So that** I can see ingredients, instructions, nutrition

**Acceptance Criteria:**
- [ ] `/recipes/{id}` route displays details
- [ ] Displays all recipe fields
- [ ] "Edit" button (if user is creator)
- [ ] "Delete" button with confirmation
- [ ] Nutrition table format
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/recipes/$recipeId.tsx`
2. Fetch recipe
3. **TDD:** Write `RecipeDetail.test.tsx`
4. Create `RecipeDetail.tsx`
5. Style elegantly

---

### User Story 3.6: Recipe Form (Create/Edit) Frontend (20h)

**As a** user  
**I want** form to create/edit recipes  
**So that** I can manage recipe data

**Acceptance Criteria:**
- [ ] `/recipes/new` route displays creation form
- [ ] `/recipes/{id}/edit` pre-fills for editing
- [ ] All fields match OpenAPI schema
- [ ] Ingredient builder (add/remove)
- [ ] Instruction builder (add/remove)
- [ ] Tag input (chips)
- [ ] Image upload or URL input
- [ ] Nutrition lookup modal
- [ ] Zod validation
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/recipes/new.tsx`
2. Create `src/routes/recipes/$recipeId/edit.tsx`
3. Create Zod schema
4. **TDD:** Write `RecipeForm.test.tsx`
5. Create `RecipeForm.tsx`
6. Create `IngredientBuilder.tsx`
7. Create `InstructionBuilder.tsx`
8. Create `NutritionSearchModal.tsx`
9. Style

---

## EPIC 4: Nutrition Data Integration (32 hours)

**Duration:** 1 week  
**Priority:** P1

### User Story 4.1: USDA API Client Backend (6h)

**Acceptance Criteria:**
- [ ] `internal/nutrition/usda_client.go` implements search
- [ ] API key from `.env`
- [ ] Parses JSON response
- [ ] Handles errors gracefully
- [ ] Unit and integration tests

**Tasks:**
1. **TDD:** Write `internal/nutrition/usda_client_test.go`
2. Write `internal/nutrition/usda_client.go`
3. Test with real API

---

### User Story 4.2: Nutrition Cache Backend (8h)

**Acceptance Criteria:**
- [ ] Cache checks before API call
- [ ] Cache hit returns immediately
- [ ] Cache miss fetches and stores
- [ ] 90-day expiry
- [ ] Test coverage ≥80%

**Tasks:**
1. **TDD:** Write `internal/nutrition/service_test.go`
2. Write `internal/nutrition/service.go`
3. Create sqlc queries
4. Run `sqlc generate`
5. **TDD:** Write repository tests
6. Write repository

---

### User Story 4.3: Nutrition Search Endpoint Backend (4h)

**Acceptance Criteria:**
- [ ] `GET /api/v1/nutrition/search?query={query}`
- [ ] Returns cached + USDA results
- [ ] Protected route
- [ ] Tests pass

**Tasks:**
1. **TDD:** Write handler tests
2. Write handler
3. Register route
4. Test manually

---

### User Story 4.4: Recipe Nutrition Calculation Backend (6h)

**Acceptance Criteria:**
- [ ] Nutrition auto-calculated from ingredients
- [ ] Calculation: sum ingredients → divide by servings
- [ ] Handles missing data
- [ ] Tests pass

**Tasks:**
1. **TDD:** Write service tests
2. Update service with `CalculateNutrition()`
3. Call in `GetRecipeByID()`
4. Update OpenAPI schema

---

### User Story 4.5: Nutrition Search Modal Frontend (8h)

**Acceptance Criteria:**
- [ ] Modal opens on "Search Nutrition" button
- [ ] Search with debounce (300ms)
- [ ] Displays results
- [ ] User selects result
- [ ] Component test passes

**Tasks:**
1. **TDD:** Write `NutritionSearchModal.test.tsx`
2. Create `NutritionSearchModal.tsx`
3. Create `src/hooks/useNutrition.ts`
4. Add to `IngredientBuilder.tsx`
5. Style modal

---

## EPIC 5: Meal Planning (40 hours)

**Duration:** 1.5 weeks  
**Priority:** P1

### User Story 5.1: Meal Plan CRUD Endpoints Backend (16h)

**Acceptance Criteria:**
- [ ] `GET /api/v1/meal-plans?start_date={date}&end_date={date}`
- [ ] `GET /api/v1/meal-plans/{date}`
- [ ] `PUT /api/v1/meal-plans/{date}`
- [ ] Supports 5 meal types
- [ ] Multiple recipes per meal OR "Out of Kitchen"
- [ ] Test coverage ≥80%

**Tasks:**
1. Create sqlc queries
2. Run `sqlc generate`
3. **TDD:** Write service tests
4. Write service
5. **TDD:** Write repository tests
6. Write repository
7. **TDD:** Write handler tests
8. Write handler
9. Register routes
10. Integration test

---

### User Story 5.2: Meal Plan Calendar UI Frontend (14h)

**Acceptance Criteria:**
- [ ] `/mealplan` displays 7-day calendar
- [ ] 5 meal slots per day
- [ ] Displays recipes or "Out of Kitchen"
- [ ] Navigation: Prev/Next Week, Today
- [ ] Click slot opens modal
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/mealplan/index.tsx`
2. Create `src/hooks/useMealPlan.ts`
3. **TDD:** Write `MealPlanCalendar.test.tsx`
4. Create `MealPlanCalendar.tsx`
5. **TDD:** Write `MealSlot.test.tsx`
6. Create `MealSlot.tsx`
7. Style with CSS Grid

---

### User Story 5.3: Meal Slot Modal Frontend (10h)

**Acceptance Criteria:**
- [ ] Modal displays recipe dropdown
- [ ] User selects recipe + servings
- [ ] "Add Another Recipe" button
- [ ] "Out of Kitchen" checkbox
- [ ] Save updates meal plan
- [ ] Component test passes

**Tasks:**
1. **TDD:** Write `MealSlotModal.test.tsx`
2. Create `MealSlotModal.tsx`
3. Use React Hook Form
4. TanStack Query mutation
5. Style

---

## EPIC 6: Grocery List Generation (38 hours)

**Duration:** 1 week  
**Priority:** P1

### User Story 6.1: Grocery List Generation Backend (16h)

**Acceptance Criteria:**
- [ ] `POST /api/v1/grocery-lists` generates list
- [ ] Aggregates recipes from meal plan
- [ ] Sums quantities
- [ ] Groups by department
- [ ] Excludes "Out of Kitchen"
- [ ] Test coverage ≥80%

**Tasks:**
1. Create sqlc queries
2. Run `sqlc generate`
3. **TDD:** Write service tests (aggregation logic)
4. Write service
5. **TDD:** Write repository tests
6. Write repository
7. **TDD:** Write handler tests
8. Write handler
9. Register route
10. Integration test

---

### User Story 6.2: Item Status Update Backend (4h)

**Acceptance Criteria:**
- [ ] `PATCH /api/v1/grocery-lists/{list_id}/items/{item_id}`
- [ ] Updates status (bought/have_on_hand/pending)
- [ ] Tests pass

**Tasks:**
1. Create sqlc query
2. **TDD:** Write handler tests
3. Write handler
4. Register route

---

### User Story 6.3: Manual Item Addition Backend (4h)

**Acceptance Criteria:**
- [ ] `POST /api/v1/grocery-lists/{list_id}/items`
- [ ] Adds manual item
- [ ] `is_manual = true`
- [ ] Tests pass

**Tasks:**
1. **TDD:** Write handler tests
2. Write handler
3. Register route

---

### User Story 6.4: Grocery List UI Frontend (14h)

**Acceptance Criteria:**
- [ ] `/grocery` route displays list
- [ ] Date range selector
- [ ] Items grouped by department
- [ ] Checkbox updates status
- [ ] "Add Manual Item" button
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/grocery/index.tsx`
2. Create `src/hooks/useGroceryList.ts`
3. **TDD:** Write `GroceryList.test.tsx`
4. Create `GroceryList.tsx`
5. **TDD:** Write `GroceryDepartment.test.tsx`
6. Create `GroceryDepartment.tsx`
7. **TDD:** Write `GroceryItem.test.tsx`
8. Create `GroceryItem.tsx`
9. Create `AddManualItemModal.tsx`
10. Style

---

## EPIC 7: Ingredient-Based Menu Recommendation (24 hours)

**Duration:** 1 week  
**Priority:** P2

### User Story 7.1: Recommendation Endpoint Backend (12h)

**Acceptance Criteria:**
- [ ] `POST /api/v1/menu/recommend` accepts ingredients
- [ ] Returns ranked recipes with match scores
- [ ] Match score = matched/total * 100
- [ ] Includes missing ingredients
- [ ] Partial matches included
- [ ] Test coverage ≥80%

**Tasks:**
1. **TDD:** Write service tests (match scoring)
2. Write service `RecommendRecipes()`
3. **TDD:** Write handler tests
4. Write handler
5. Register route
6. Integration test

---

### User Story 7.2: Recommended Menu UI Frontend (12h)

**Acceptance Criteria:**
- [ ] `/menu` route displays input form
- [ ] User enters ingredients (chips)
- [ ] Results in French menu style
- [ ] Shows: title, match score, missing ingredients
- [ ] Click navigates to recipe
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/menu/index.tsx`
2. Create `src/hooks/useMenu.ts`
3. **TDD:** Write `RecommendedMenu.test.tsx`
4. Create `RecommendedMenu.tsx`
5. **TDD:** Write `MenuCard.test.tsx`
6. Create `MenuCard.tsx`
7. Style with elegant aesthetic

---

## EPIC 8: Cookbooks (36 hours)

**Duration:** 1 week  
**Priority:** P2

### User Story 8.1: Cookbook CRUD Endpoints Backend (14h)

**Acceptance Criteria:**
- [ ] `POST /api/v1/cookbooks` creates cookbook
- [ ] `GET /api/v1/cookbooks` lists all
- [ ] `GET /api/v1/cookbooks/{id}` with recipes
- [ ] `DELETE /api/v1/cookbooks/{id}` soft-deletes
- [ ] `POST /api/v1/cookbooks/{cookbook_id}/recipes/{recipe_id}` adds recipe
- [ ] `DELETE` removes recipe
- [ ] Test coverage ≥80%

**Tasks:**
1. Create sqlc queries
2. Run `sqlc generate`
3. **TDD:** Write service tests
4. Write service
5. **TDD:** Write repository tests
6. Write repository
7. **TDD:** Write handler tests
8. Write handler
9. Register routes
10. Integration test

---

### User Story 8.2: Cookbook List UI Frontend (10h)

**Acceptance Criteria:**
- [ ] `/cookbooks` displays list
- [ ] Cards show: name, description, recipe count
- [ ] Click navigates to detail
- [ ] "New Cookbook" button
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/cookbooks/index.tsx`
2. Create `src/hooks/useCookbooks.ts`
3. **TDD:** Write `CookbookCard.test.tsx`
4. Create `CookbookCard.tsx`
5. **TDD:** Write `CookbookList.test.tsx`
6. Create `CookbookList.tsx`
7. Create `CreateCookbookModal.tsx`
8. Style

---

### User Story 8.3: Cookbook Detail Page Frontend (12h)

**Acceptance Criteria:**
- [ ] `/cookbooks/{id}` displays details
- [ ] Shows: name, description, creator
- [ ] Grid of recipe cards
- [ ] "Add Recipe" button
- [ ] "Remove" button per recipe
- [ ] "Delete Cookbook" button
- [ ] Component test passes

**Tasks:**
1. Create `src/routes/cookbooks/$cookbookId.tsx`
2. **TDD:** Write `CookbookDetail.test.tsx`
3. Create `CookbookDetail.tsx`
4. Create `AddRecipeModal.tsx`
5. Style

---

## EPIC 9: Polish, Deployment & Production (82 hours)

**Duration:** 1.5 weeks  
**Priority:** P0

### User Story 9.1: E2E Test Suite (16h)

**Acceptance Criteria:**
- [ ] E2E tests for: registration, login, logout
- [ ] E2E tests for: recipe CRUD
- [ ] E2E tests for: meal plan, grocery list
- [ ] E2E tests for: recommendation
- [ ] E2E tests for: cookbooks
- [ ] All tests pass in CI

**Tasks:**
1. Create `pinecone-e2e` repository
2. Install Playwright
3. Write `e2e/auth.spec.ts`
4. Write `e2e/recipes.spec.ts`
5. Write `e2e/mealplan.spec.ts`
6. Write `e2e/menu.spec.ts`
7. Write `e2e/cookbooks.spec.ts`
8. Add to CI pipeline
9. Run seed script before tests

---

### User Story 9.2: Production Docker Compose (12h)

**Acceptance Criteria:**
- [ ] `docker-compose.prod.yml` created
- [ ] Services: db, api, caddy
- [ ] Health checks for all
- [ ] Database backups automated
- [ ] Logs persisted

**Tasks:**
1. Create `docker-compose.prod.yml`
2. Create `Caddyfile`
3. Write `Dockerfile` for backend
4. Write `Dockerfile` for frontend
5. Configure health checks
6. Write backup script
7. Document deployment

---

### User Story 9.3: CD Pipeline (10h)

**Acceptance Criteria:**
- [ ] `.github/workflows/cd.yml` created
- [ ] Triggered on push to `main` (staging)
- [ ] Triggered on tag `v*` (production)
- [ ] Builds and pushes Docker images
- [ ] SSHs and restarts services
- [ ] Notifies on success/failure

**Tasks:**
1. Create `.github/workflows/cd.yml`
2. Add Docker build/push
3. Add SSH deployment
4. Configure GitHub Environments
5. Add secrets
6. Test staging deployment
7. Test production deployment

---

### User Story 9.4: Sentry Integration (6h)

**Acceptance Criteria:**
- [ ] Sentry DSN in `.env`
- [ ] Backend initialized in `main.go`
- [ ] Frontend initialized in `main.tsx`
- [ ] Error boundary for frontend
- [ ] All errors include user ID and request ID
- [ ] Alerts configured

**Tasks:**
1. Create Sentry project
2. Install backend SDK
3. Initialize in `main.go`
4. Install frontend SDK
5. Initialize in `main.tsx`
6. Add error boundary
7. Configure alerts

---

### User Story 9.5: Performance Optimization (10h)

**Acceptance Criteria:**
- [ ] Database indexes added
- [ ] API responses < 500ms
- [ ] Recipe search < 300ms
- [ ] Frontend bundle < 500KB
- [ ] Lighthouse score ≥90
- [ ] Images lazy-loaded

**Tasks:**
1. Add DB indexes
2. Run Go benchmarks
3. Optimize slow queries
4. Run Lighthouse audit
5. Code-split routes
6. Optimize bundle
7. Add lazy loading

---

### User Story 9.6: Security Audit (8h)

**Acceptance Criteria:**
- [ ] `go vet` passes
- [ ] `golangci-lint` passes
- [ ] `npm audit` clean
- [ ] OWASP Top 10 checklist complete
- [ ] SQL injection verified (sqlc)
- [ ] XSS verified (React escaping)
- [ ] CSRF verified (SameSite cookies)

**Tasks:**
1. Run `go vet`
2. Run `golangci-lint`
3. Run `npm audit` and fix
4. Complete OWASP checklist
5. Security review auth flow
6. Penetration test
7. Document in `docs/SECURITY.md`

---

### User Story 9.7: User Acceptance Testing (12h)

**Acceptance Criteria:**
- [ ] 2-3 household members test
- [ ] UAT checklist completed
- [ ] Feedback documented
- [ ] Critical bugs fixed
- [ ] Success metrics baseline captured

**Tasks:**
1. Deploy to staging
2. Create UAT checklist
3. Invite testers
4. Collect feedback
5. Triage issues
6. Fix P0 issues
7. Document in `docs/UAT_RESULTS.md`

---

### User Story 9.8: Documentation Finalization (8h)

**Acceptance Criteria:**
- [ ] `README.md` complete
- [ ] `docs/BRD.md` finalized
- [ ] `docs/TDD.md` finalized
- [ ] `docs/DEPLOYMENT.md` created
- [ ] `docs/DEVELOPMENT.md` created
- [ ] `docs/API.md` created
- [ ] `docs/SECURITY.md` created

**Tasks:**
1. Write `README.md`
2. Finalize BRD/TDD
3. Write `docs/DEPLOYMENT.md`
4. Write `docs/DEVELOPMENT.md`
5. Write `docs/API.md`
6. Write `docs/SECURITY.md`
7. Proofread all

---

## Summary

| Epic | Hours | Weeks |
|------|-------|-------|
| Epic 1: Foundation | 52 | 2 |
| Epic 2: Authentication | 42 | 1 |
| Epic 3: Recipe Management | 76 | 2 |
| Epic 4: Nutrition Integration | 32 | 1 |
| Epic 5: Meal Planning | 40 | 1.5 |
| Epic 6: Grocery Lists | 38 | 1 |
| Epic 7: Recommendation | 24 | 1 |
| Epic 8: Cookbooks | 36 | 1 |
| Epic 9: Polish & Deployment | 82 | 1.5 |
| **TOTAL** | **422** | **~10.5** |

**Timeline:** Nov 9, 2025 → Jan 31, 2026 (12 weeks with buffer)

---

**Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-08 | GhostDev | Initial epic breakdown |
