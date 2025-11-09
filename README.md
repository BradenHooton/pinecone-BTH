# Pinecone Recipe Management System

A comprehensive household recipe management system built with Go, PostgreSQL, React, and TypeScript.

## Features

### ✅ Completed Features

- **User Authentication**
  - Registration with email and password
  - Secure login with JWT tokens (HTTP-only cookies)
  - Bcrypt password hashing (cost 12)
  - 24-hour session expiry

- **Recipe Management**
  - Create, read, update, and delete recipes
  - Rich recipe data: ingredients, instructions, tags
  - Image upload support (5MB max, jpg/png/webp)
  - Search and filter recipes
  - Pagination support
  - Serving size and timing information

- **Nutrition Data Integration**
  - USDA FoodData Central API integration (stub implementation)
  - Nutrition search with 300ms debounce
  - 90-day caching of nutrition data
  - Nutrition facts display for ingredients

- **Meal Planning**
  - Weekly meal plan calendar (7 days x 5 meal types)
  - Meal types: breakfast, lunch, snack, dinner, dessert
  - Multiple recipes per meal slot
  - "Out of Kitchen" option for eating out
  - Date range queries
  - Automatic date normalization to midnight UTC

- **Grocery List Generation**
  - Automatic generation from meal plans
  - Intelligent ingredient aggregation
  - Quantity scaling based on servings
  - 10 grocery departments (produce, meat, seafood, dairy, bakery, frozen, pantry, spices, beverages, other)
  - Item status tracking (pending, bought, have_on_hand)
  - Manual item addition
  - Interactive checkbox UI

- **Ingredient-Based Menu Recommendation**
  - Find recipes based on available ingredients
  - Match scoring: (matched / total) * 100
  - Shows matched and missing ingredients
  - Ingredient normalization (handles plurals)
  - French menu-style presentation
  - Match quality indicators (Excellent/Good/Fair/Partial)

- **Cookbooks**
  - Create and organize recipe collections
  - Add/remove recipes from cookbooks
  - Soft delete support
  - Recipe count tracking
  - Prevent duplicate recipes in same cookbook

## Tech Stack

### Backend
- **Language**: Go 1.21
- **Database**: PostgreSQL 16
- **Router**: Chi v5
- **Database Driver**: pgx/v5
- **Migrations**: Goose
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Password Security**: Bcrypt (cost 12)
- **Logging**: Standard library slog (JSON format)
- **Architecture**: Hexagonal (Handler → Service → Repository)

### Frontend
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite
- **Routing**: TanStack Router
- **Data Fetching**: TanStack Query (React Query)
- **State Management**: Zustand (with persistence)
- **Form Handling**: React Hook Form
- **Validation**: Zod
- **Styling**: Inline styles with design tokens

### Design System
- **Colors**: Forest green (#3A7D44), warm off-white
- **Typography**: Playfair Display (headings), Inter (body)
- **Spacing**: 4px base unit scale

## Security Features

- **Authentication**: JWT tokens in HTTP-only cookies
- **Authorization**: User-scoped data access
- **Password Security**: Bcrypt hashing with cost 12
- **SQL Injection Prevention**: Parameterized queries throughout
- **XSS Prevention**: React auto-escaping
- **CSRF Protection**: SameSite cookies
- **Input Validation**: Server-side validation on all endpoints
- **Soft Deletes**: Data preservation
- **Rate Limiting**: 100 requests per minute per IP
- **CORS**: Configurable allowed origins

## Performance

- **Backend Compilation**: ✓ (no errors)
- **Frontend TypeScript**: ✓ (no type errors)
- **Frontend Bundle**: 398.99 kB (115.34 kB gzipped)
- **Database Indexes**: 24 indexes for optimal query performance
- **Caching**: Nutrition data cached for 90 days
- **Pagination**: All list endpoints support pagination

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 16
- npm or yarn

### Backend Setup
```bash
cd backend
go mod download
goose -dir internal/store/migrations postgres $DATABASE_URL up
go run cmd/server/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## Testing

### Backend
```bash
cd backend
go test ./...
go vet ./...
```

### Frontend
```bash
cd frontend
npm run type-check
npm run build
```

## License

This project is for personal/household use.
