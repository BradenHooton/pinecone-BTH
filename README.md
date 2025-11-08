# Pinecone Recipe Management System

A single-household recipe management and meal planning application designed to centralize recipes, streamline meal planning, and automatically generate organized grocery lists.

## 🌲 Project Overview

Pinecone eliminates decision fatigue and scattered recipe sources by providing:
- **Recipe Management**: Store, search, and organize all household recipes
- **Meal Planning**: Visual calendar for planning meals with future date support
- **Grocery Lists**: Auto-generated lists organized by store department
- **Ingredient Recommendations**: Suggest recipes based on available ingredients
- **Cookbooks**: Create themed collections of recipes

## 🛠 Tech Stack

### Backend
- **Language**: Go 1.21+
- **HTTP Router**: chi with middleware (rate limiting, CORS, auth)
- **Database**: PostgreSQL 16
- **Query Builder**: sqlc (type-safe SQL)
- **Migrations**: goose
- **Authentication**: JWT (HTTP-only cookies) + bcrypt

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **State Management**: Zustand
- **Data Fetching**: TanStack Query
- **Routing**: TanStack Router
- **Forms**: React Hook Form + Zod validation

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Web Server**: Caddy (production)
- **CI/CD**: GitHub Actions

## 🚀 Quick Start

### Prerequisites

| Tool | Version | Installation |
|------|---------|--------------|
| Go | 1.21+ | https://go.dev/doc/install |
| Node.js | 20+ | https://nodejs.org/ |
| Docker | 24+ | https://docs.docker.com/get-docker/ |
| Make | Any | Pre-installed on macOS/Linux |

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/BradenHooton/pinecone-BTH.git
   cd pinecone-BTH
   ```

2. **Install dependencies**
   ```bash
   make install
   ```

3. **Set up environment variables**
   ```bash
   # Backend
   cd backend
   cp .env.example .env.dev
   # Edit .env.dev with your settings (DATABASE_URL, JWT_SECRET, etc.)

   # Frontend
   cd ../frontend
   cp .env.example .env.dev
   # Edit .env.dev with API URL
   ```

4. **Start development environment**
   ```bash
   # From project root
   make dev
   ```

   This will start:
   - PostgreSQL on `localhost:5432`
   - Backend API on `localhost:8080`
   - Frontend on `localhost:5173`

5. **Access the application**
   - Frontend: http://localhost:5173
   - API Health Check: http://localhost:8080/health

## 📚 Development Workflow

### Common Commands

```bash
# Start full stack
make dev

# Run backend only
make backend-run

# Run frontend only
make frontend-dev

# Run tests
make backend-test
make frontend-test

# Database migrations
make migrate-up
make migrate-status
make migrate-create NAME=add_new_table

# Generate sqlc code
make sqlc-generate

# Linting
make backend-lint
make frontend-lint

# View logs
make dev-logs
make dev-logs-api
make dev-logs-web
```

### Project Structure

```
pinecone-BTH/
├── backend/
│   ├── cmd/server/         # Application entry point
│   ├── internal/
│   │   ├── auth/           # Authentication domain
│   │   ├── recipe/         # Recipe domain
│   │   ├── mealplan/       # Meal planning domain
│   │   ├── grocery/        # Grocery list domain
│   │   ├── cookbook/       # Cookbook domain
│   │   ├── nutrition/      # USDA API integration
│   │   ├── models/         # Shared domain models
│   │   ├── middleware/     # HTTP middleware
│   │   ├── config/         # Configuration management
│   │   └── store/
│   │       ├── queries/    # sqlc SQL queries
│   │       └── migrations/ # Database migrations
│   ├── pkg/                # Public libraries (JWT, validator)
│   └── api/                # OpenAPI specification
├── frontend/
│   └── src/
│       ├── components/     # React components
│       ├── hooks/          # Custom React hooks
│       ├── lib/            # API client, utilities
│       ├── routes/         # TanStack Router definitions
│       ├── store/          # Zustand state stores
│       ├── styles/         # Global CSS
│       └── tokens/         # Design tokens
├── docker-compose.yml      # Development environment
├── Makefile                # Common tasks
└── docs/                   # Documentation
```

## 🧪 Testing

This project follows **Test-Driven Development (TDD)** principles.

### Backend Testing

```bash
# Run all tests
make backend-test

# Run with coverage
make backend-test-coverage

# Target: ≥80% coverage for services and repositories
```

### Frontend Testing

```bash
# Run all tests
make frontend-test

# Run with UI
make frontend-test-ui

# Run with coverage
make frontend-test-coverage

# Target: ≥70% coverage for components and hooks
```

## 🗄 Database

### Migrations

```bash
# Apply all migrations
make migrate-up

# Check migration status
make migrate-status

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=create_new_table
```

### Schema

The database includes 13 tables:
- `users` - User accounts
- `recipes` - Recipe master data
- `recipe_ingredients` - Recipe ingredients
- `recipe_instructions` - Step-by-step instructions
- `recipe_tags` - Recipe categorization
- `meal_plans` - Meal plan dates
- `meal_plan_recipes` - Recipes scheduled in meal plans
- `grocery_lists` - Generated grocery lists
- `grocery_list_items` - Individual grocery items
- `cookbooks` - Recipe collections
- `cookbook_recipes` - Cookbook-recipe associations
- `nutrition_cache` - Cached USDA nutrition data

## 🎨 Design System

### Colors
- **Primary**: Forest Green (#3A7D44)
- **Background**: Warm Off-white (#FAF9F6)
- **Text**: Black (#121212)

### Typography
- **Headings**: Playfair Display (serif)
- **Body**: Inter (sans-serif)

### Tokens
All design tokens are defined in `frontend/src/tokens/`:
- `colors.ts` - Color palette
- `typography.ts` - Font families, sizes, weights
- `spacing.ts` - Spacing scale, breakpoints, shadows

## 📖 Documentation

- [Business Requirements (BRD)](BRD.md)
- [Technical Design (TDD)](TDD.md)
- [Database Schema](DATABASE_SCHEMA.md)
- [Epic Breakdown](EPIC_BREAKDOWN.md)
- [Project Roadmap](PROJECT_ROADMAP.md)
- [Developer Onboarding](DEVELOPER_ONBOARDING.md)
- [Deployment Guide](DEPLOYMENT_GUIDE.md)
- [Documentation Index](DOCUMENTATION_INDEX.md)

## 🚢 Deployment

### Production Build

```bash
# Build Docker images
make docker-build

# Or use Docker Compose
docker-compose -f docker-compose.prod.yml up -d
```

See [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) for detailed production deployment instructions.

## 🤝 Contributing

This project follows strict coding standards:

1. **TDD Workflow**: Write tests before implementation (Red → Green → Refactor)
2. **Code Quality**: All PRs must pass linting and tests
3. **Coverage**: Meet minimum coverage targets (Backend ≥80%, Frontend ≥70%)
4. **Conventional Commits**: Use semantic commit messages

### Pull Request Process

1. Create a feature branch from `develop`
2. Write tests first (TDD)
3. Implement feature
4. Ensure all tests pass
5. Run linting (`make backend-lint frontend-lint`)
6. Submit PR with clear description

## 📋 Roadmap

**Current Status**: Epic 1 Complete - Foundation & Infrastructure ✓

### Upcoming Milestones

- **M2** (Week 3): User Authentication
- **M3** (Week 4-5): Recipe Management
- **M4** (Week 6): Nutrition Data Integration
- **M5** (Week 7-8): Meal Planning
- **M6** (Week 9): Grocery List Generation
- **M7** (Week 10): Ingredient Recommendations
- **M8** (Week 11): Cookbooks
- **M9** (Week 12): Production Launch

See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for the complete timeline.

## 📝 License

Private household project - Not licensed for public use.

## 👨‍💻 Author

**Braden Hooton**
GitHub: [@BradenHooton](https://github.com/BradenHooton)

---

**Built with ❤️ using Test-Driven Development**
