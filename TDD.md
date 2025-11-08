# Technical Design Document (TDD)
## Pinecone Recipe Management & Meal Planning System

**Version:** 1.0  
**Date:** 2025-11-08  
**Status:** Approved  
**Author:** GhostDev

---

## Table of Contents
1. [System Architecture](#1-system-architecture)
2. [Database Design](#2-database-design)
3. [API Specification](#3-api-specification)
4. [Authentication & Security](#4-authentication--security)
5. [External Integrations](#5-external-integrations)
6. [Configuration Management](#6-configuration-management)
7. [Testing Strategy](#7-testing-strategy)
8. [Deployment Architecture](#8-deployment-architecture)

---

## 1. System Architecture

### 1.1 Technology Stack

#### Backend
- **Language:** Go 1.21+
- **HTTP Router:** chi, chi middleware, go-chi/httprate
- **Database Driver:** pgx
- **Query Builder:** sqlc
- **Migrations:** goose
- **Authentication:** bcrypt, golang-jwt
- **Configuration:** viper, godotenv
- **Logging:** log/slog
- **Error Handling:** cockroachdb/errors
- **Testing:** testify, go-testcontainers
- **Validation:** go-playground/validator

#### Frontend
- **Framework:** React 18+ with TypeScript
- **Build Tool:** Vite
- **State Management:** Zustand
- **Data Fetching:** TanStack Query
- **Routing:** TanStack Router
- **Forms:** React Hook Form + Zod
- **Styling:** CSS with Design Tokens
- **Testing:** Vitest, React Testing Library
- **E2E Testing:** Playwright

#### Infrastructure
- **Database:** PostgreSQL 16
- **Containerization:** Docker, Docker Compose
- **Web Server:** Caddy 2
- **CI/CD:** GitHub Actions
- **Error Reporting:** Sentry

### 1.2 Architecture Principles

1. **Modular Monolith:** Single Go application with clean package boundaries
2. **Hexagonal Architecture:** Services use interfaces for repositories
3. **Test-Driven Development:** Tests written before implementation
4. **Security First:** JWT authentication, bcrypt passwords, parameterized queries
5. **API Contract First:** OpenAPI specification drives implementation

### 1.3 Package Structure

#### Backend (`pinecone-api`)
```
pinecone-api/
├── api/                    # OpenAPI specification
│   └── openapi.yaml
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── config/
│   └── grocery_departments.yaml
├── docs/
│   ├── BRD.md
│   ├── TDD.md
│   └── diagrams/
├── internal/
│   ├── auth/               # Authentication domain
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── recipe/             # Recipe domain
│   ├── mealplan/           # Meal planning domain
│   ├── grocery/            # Grocery list domain
│   ├── cookbook/           # Cookbook domain
│   ├── nutrition/          # Nutrition integration
│   ├── models/             # Shared domain models
│   ├── middleware/         # HTTP middleware
│   ├── config/             # Configuration loading
│   └── store/
│       ├── queries/        # sqlc SQL files
│       └── migrations/     # Goose migrations
├── pkg/                    # Public libraries
│   ├── jwt/
│   └── validator/
├── uploads/                # Recipe images (gitignored)
├── .env.dev                # Local environment (gitignored)
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── sqlc.yaml
```

#### Frontend (`pinecone-web`)
```
pinecone-web/
├── public/
├── src/
│   ├── components/
│   │   ├── common/         # Reusable UI components
│   │   ├── recipe/
│   │   ├── mealplan/
│   │   ├── grocery/
│   │   └── menu/
│   ├── hooks/              # Custom React hooks
│   ├── lib/                # API client, utilities
│   ├── routes/             # TanStack Router definitions
│   ├── store/              # Zustand state stores
│   ├── styles/             # Global CSS
│   ├── tokens/             # Design tokens
│   ├── App.tsx
│   └── main.tsx
├── .env.dev
├── docker-compose.yml
├── Dockerfile
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── vitest.config.ts
```

---

## 2. Database Design

### 2.1 Key Design Decisions

1. **UUIDs for Primary Keys:** Prevents enumeration attacks, enables distributed ID generation
2. **Soft Deletes:** `deleted_at` timestamp allows data recovery and audit trails
3. **Computed Columns:** `total_time_minutes` auto-calculated (prep + cook)
4. **ENUMs for Fixed Lists:** `meal_type`, `grocery_department`, `grocery_item_status`
5. **Normalization:** Recipe ingredients, instructions, tags in separate tables (1:N)
6. **Many-to-Many:** `cookbook_recipes` junction table for flexible organization

### 2.2 Core Tables

#### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

#### recipes
```sql
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id),
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
```

#### recipe_ingredients
```sql
CREATE TYPE grocery_department AS ENUM (
    'produce', 'meat', 'seafood', 'dairy', 'bakery',
    'frozen', 'pantry', 'spices', 'beverages', 'other'
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
```

#### meal_plans
```sql
CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'snack', 'dinner', 'dessert');

CREATE TABLE meal_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_date DATE NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

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
```

#### grocery_lists
```sql
CREATE TYPE grocery_item_status AS ENUM ('pending', 'bought', 'have_on_hand');

CREATE TABLE grocery_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_date_range CHECK (end_date >= start_date)
);

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
```

#### cookbooks
```sql
CREATE TABLE cookbooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by_user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE cookbook_recipes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cookbook_id UUID NOT NULL REFERENCES cookbooks(id) ON DELETE CASCADE,
    recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cookbook_id, recipe_id)
);
```

#### nutrition_cache
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
```

### 2.3 Indexes

```sql
-- Performance optimization indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipes_title ON recipes(title) WHERE deleted_at IS NULL;
CREATE INDEX idx_recipes_created_by ON recipes(created_by_user_id);
CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id);
CREATE INDEX idx_recipe_ingredients_name ON recipe_ingredients(ingredient_name);
CREATE INDEX idx_recipe_tags_recipe_id ON recipe_tags(recipe_id);
CREATE INDEX idx_recipe_tags_tag_name ON recipe_tags(tag_name);
CREATE INDEX idx_meal_plans_date ON meal_plans(plan_date);
CREATE INDEX idx_grocery_lists_dates ON grocery_lists(start_date, end_date);
CREATE INDEX idx_nutrition_cache_food_name ON nutrition_cache(food_name);
```

---

## 3. API Specification

### 3.1 Base URL
- **Development:** `http://localhost:8080/api/v1`
- **Production:** `https://pinecone.example.com/api/v1`

### 3.2 Authentication
- **Type:** Cookie-based JWT
- **Cookie Name:** `jwt_token`
- **Cookie Properties:** HttpOnly, Secure, SameSite=Strict
- **Token Expiry:** 24 hours

### 3.3 Core Endpoints

#### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login with credentials
- `POST /auth/logout` - Logout and invalidate token

#### Recipes
- `GET /recipes` - List recipes (with pagination, search, filter, sort)
- `POST /recipes` - Create recipe
- `GET /recipes/{id}` - Get recipe by ID
- `PUT /recipes/{id}` - Update recipe
- `DELETE /recipes/{id}` - Soft-delete recipe
- `POST /recipes/upload-image` - Upload recipe image

#### Nutrition
- `GET /nutrition/search?query={query}` - Search USDA nutrition data

#### Meal Plans
- `GET /meal-plans?start_date={date}&end_date={date}` - Get meal plans for range
- `GET /meal-plans/{date}` - Get meal plan for specific date
- `PUT /meal-plans/{date}` - Create/update meal plan

#### Grocery Lists
- `POST /grocery-lists` - Generate grocery list from meal plan
- `PATCH /grocery-lists/{list_id}/items/{item_id}` - Update item status
- `POST /grocery-lists/{list_id}/items` - Add manual item

#### Menu Recommendation
- `POST /menu/recommend` - Get recipe recommendations based on ingredients

#### Cookbooks
- `GET /cookbooks` - List all cookbooks
- `POST /cookbooks` - Create cookbook
- `GET /cookbooks/{id}` - Get cookbook with recipes
- `DELETE /cookbooks/{id}` - Soft-delete cookbook
- `POST /cookbooks/{cookbook_id}/recipes/{recipe_id}` - Add recipe to cookbook
- `DELETE /cookbooks/{cookbook_id}/recipes/{recipe_id}` - Remove recipe from cookbook

### 3.4 Standard Response Format

**Success Response:**
```json
{
  "data": { ... },
  "meta": {
    "timestamp": "2025-11-08T12:00:00Z",
    "request_id": "uuid"
  }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": [...]
  },
  "meta": {
    "timestamp": "2025-11-08T12:00:00Z",
    "request_id": "uuid"
  }
}
```

---

## 4. Authentication & Security

### 4.1 Password Security
- **Hashing Algorithm:** bcrypt
- **Cost Factor:** 12
- **Salt:** Automatically generated per password

### 4.2 JWT Configuration
```go
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Token expiry: 24 hours
expirationTime := time.Now().Add(24 * time.Hour)
```

### 4.3 Security Headers
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: no-referrer-when-downgrade
Strict-Transport-Security: max-age=31536000
```

### 4.4 Rate Limiting
- **Rate:** 100 requests per minute per IP
- **Implementation:** go-chi/httprate

### 4.5 CORS Configuration
```go
cors.Options{
    AllowedOrigins:   []string{"http://localhost:5173", "https://pinecone.example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
    AllowedHeaders:   []string{"Content-Type"},
    AllowCredentials: true,
}
```

---

## 5. External Integrations

### 5.1 USDA FoodData Central API

**Purpose:** Fetch nutrition data for recipe ingredients

**Endpoint:** `https://api.nal.usda.gov/fdc/v1/foods/search`

**Authentication:** API Key (free tier: 1,000 requests/hour)

**Example Request:**
```bash
GET https://api.nal.usda.gov/fdc/v1/foods/search?query=chicken+breast&api_key=YOUR_KEY
```

**Caching Strategy:**
1. Check local `nutrition_cache` table first
2. If cache miss, fetch from USDA API
3. Store result in cache with 90-day expiry
4. Return cached or fresh data

**Error Handling:**
- **Rate Limit Exceeded:** Return cached results only + warning
- **API Down:** Allow manual nutrition entry
- **No Match Found:** Proceed without nutrition data

---

## 6. Configuration Management

### 6.1 Environment Variables

**File:** `.env.dev` (local), GitHub Secrets (production)

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pinecone_dev?sslmode=disable

# JWT
JWT_SECRET=your-256-bit-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

# USDA API
USDA_API_KEY=your-usda-api-key
USDA_API_BASE_URL=https://api.nal.usda.gov/fdc/v1

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# File Uploads
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE_MB=5

# CORS
ALLOWED_ORIGINS=http://localhost:5173,https://pinecone.example.com

# Logging
LOG_LEVEL=info  # debug, info, warn, error
```

### 6.2 YAML Configuration

**File:** `config/grocery_departments.yaml`

```yaml
departments:
  - id: produce
    name: Produce
    description: Fresh fruits and vegetables
    order: 1
  - id: meat
    name: Meat & Poultry
    order: 2
  - id: dairy
    name: Dairy & Eggs
    order: 4
  # ... additional departments
```

---

## 7. Testing Strategy

### 7.1 Test-Driven Development Workflow

```
1. RED: Write failing test
2. GREEN: Write minimum code to pass
3. REFACTOR: Improve code quality
4. COMMIT: Save working state
```

### 7.2 Test Layers

#### Unit Tests
- **Purpose:** Test individual functions in isolation
- **Backend:** `testify` (assertions, mocking)
- **Frontend:** `vitest`, `@testing-library/react`
- **Coverage Target:** ≥80% (backend services/repos), ≥70% (frontend)

#### Integration Tests
- **Purpose:** Test component interactions
- **Backend:** `go-testcontainers` (real PostgreSQL in Docker)
- **Scope:** Repository layer queries against real DB

#### E2E Tests
- **Purpose:** Test complete user flows
- **Tool:** Playwright
- **Critical Flows:**
  - User registration and login
  - Recipe creation and editing
  - Meal plan creation and grocery list generation
  - Ingredient-based menu recommendation

### 7.3 Test Examples

**Backend Unit Test:**
```go
func TestRecipeService_Create_CalculatesTotalTime(t *testing.T) {
    // ARRANGE
    mockRepo := new(MockRecipeRepository)
    service := recipe.NewService(mockRepo)
    
    recipeInput := &models.Recipe{
        Title:           "Grilled Chicken",
        PrepTimeMinutes: 15,
        CookTimeMinutes: 30,
    }
    
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // ACT
    err := service.CreateRecipe(context.Background(), recipeInput)
    
    // ASSERT
    assert.NoError(t, err)
    assert.Equal(t, 45, recipeInput.TotalTimeMinutes)
    mockRepo.AssertExpectations(t)
}
```

**Frontend Component Test:**
```typescript
describe('RecipeCard', () => {
  it('should display recipe title and prep time', () => {
    const recipe = {
      id: '123',
      title: 'Grilled Chicken',
      prep_time_minutes: 15,
      cook_time_minutes: 30,
    };

    render(<RecipeCard recipe={recipe} />);

    expect(screen.getByText('Grilled Chicken')).toBeInTheDocument();
    expect(screen.getByText('15 min prep')).toBeInTheDocument();
  });
});
```

---

## 8. Deployment Architecture

### 8.1 Production Stack

```
[Internet] → [Caddy (Port 80/443)]
              ├── /api/* → [Go API Container (Port 8080)]
              └── /* → [React Static Files]
                         
[Go API] → [PostgreSQL Container (Port 5432)]
[Go API] → [File Storage: /var/www/pinecone/uploads]
```

### 8.2 Docker Compose Configuration

**File:** `docker-compose.prod.yml`

```yaml
version: '3.8'

services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: pinecone
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    image: ghcr.io/bhooton/pinecone-api:latest
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@db:5432/pinecone
      JWT_SECRET: ${JWT_SECRET}
      USDA_API_KEY: ${USDA_API_KEY}
    volumes:
      - ./uploads:/app/uploads
    depends_on:
      db:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]

  caddy:
    image: caddy:2-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - ./web/dist:/var/www/html
      - caddy_data:/data

volumes:
  postgres_data:
  caddy_data:
```

### 8.3 CI/CD Pipeline

**Triggers:**
- Push to `main` → Deploy to Staging
- Git tag `v*` → Deploy to Production

**Steps:**
1. Checkout code
2. Run tests (backend + frontend)
3. Build Docker images
4. Push to GitHub Container Registry
5. SSH to server and pull images
6. Restart services via `docker-compose up -d`
7. Health check
8. Notify on success/failure

---

## Appendix A: Design Tokens

### Colors
```typescript
export const colors = {
  primary: '#3A7D44',        // Forest green
  background: '#FAF9F6',     // Warm off-white
  text: '#121212',           // Black
  textSecondary: '#4A4A4A',
  border: '#D1D1D1',
  error: '#D32F2F',
  success: '#388E3C',
};
```

### Typography
```typescript
export const typography = {
  fontFamily: {
    heading: '"Playfair Display", serif',
    body: '"Inter", sans-serif',
  },
  fontSize: {
    xs: '0.75rem',
    sm: '0.875rem',
    base: '1rem',
    lg: '1.125rem',
    xl: '1.25rem',
    '2xl': '1.5rem',
  },
};
```

### Spacing
```typescript
export const spacing = {
  xs: '4px',
  sm: '8px',
  md: '16px',
  lg: '24px',
  xl: '32px',
  '2xl': '48px',
};
```

---

**Document Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-08 | GhostDev | Initial TDD creation and approval |
