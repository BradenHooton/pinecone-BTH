# 🌲 Pinecone Recipe Management System

A comprehensive recipe management and meal planning application designed to help households organize recipes, plan meals, and generate grocery lists automatically.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Project Status: Planning](https://img.shields.io/badge/status-planning-blue.svg)]()
[![Documentation](https://img.shields.io/badge/docs-complete-brightgreen.svg)]()

---

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Documentation](#documentation)
- [Project Structure](#project-structure)
- [Development](#development)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 Overview

Pinecone solves the everyday challenge of meal planning and grocery shopping for households. It provides:

- **Centralized Recipe Database**: Store all your recipes in one searchable location
- **Visual Meal Planning**: Drag-and-drop calendar interface for weekly meal planning
- **Smart Grocery Lists**: Auto-generated shopping lists organized by store department
- **Ingredient Recommendations**: Suggest recipes based on what you have on hand
- **Recipe Collections**: Organize recipes into themed cookbooks

### Problem Statement

Households struggle with:
- Recipes scattered across bookmarks, notebooks, and websites
- Inefficient grocery shopping (multiple trips, forgotten ingredients)
- Decision fatigue around "What's for dinner?"
- Manual meal planning and list creation

### Solution

Pinecone provides a single-household application that consolidates recipes, streamlines meal planning, and automates grocery list generation.

---

## ✨ Features

### MVP (Phase 1)

- ✅ **User Authentication**: Email/password with JWT-based sessions
- ✅ **Recipe Management**: Full CRUD with rich metadata (nutrition, prep time, tags)
- ✅ **Image Support**: Upload recipe images (≤5MB) or link URLs
- ✅ **Advanced Search**: Filter by tags, prep time, and nutrition values
- ✅ **Nutrition Integration**: USDA FoodData Central API with caching
- ✅ **Meal Planning Calendar**: 5 meal slots per day (Breakfast, Lunch, Snack, Dinner, Dessert)
- ✅ **Future Date Support**: Plan meals weeks or months in advance
- ✅ **"Out of Kitchen" Option**: Mark meals when eating out
- ✅ **Grocery List Generation**: Aggregate ingredients by date range
- ✅ **Department Organization**: Group items by grocery store section
- ✅ **Item Status Tracking**: Mark items as "Bought" or "Have on Hand"
- ✅ **Ingredient-Based Recommendations**: Find recipes with partial ingredient matches
- ✅ **Cookbooks**: Create themed recipe collections

### Future Enhancements (Post-MVP)

- Multi-household support
- Public recipe sharing
- Native mobile apps (iOS/Android)
- Grocery delivery integration
- Automatic recipe import from URLs
- Calorie/macro tracking per user

---

## 🛠 Technology Stack

### Backend

- **Language**: Go 1.21+
- **Web Framework**: chi router with middleware
- **Database**: PostgreSQL 16 (pgx driver)
- **Query Builder**: sqlc (generates type-safe Go from SQL)
- **Migrations**: goose (versioned database migrations)
- **Authentication**: bcrypt + JWT (golang-jwt)
- **Configuration**: viper + godotenv
- **Testing**: testify, go-testcontainers
- **Logging**: log/slog (structured JSON logging)

### Frontend

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **State Management**: Zustand
- **Data Fetching**: TanStack Query
- **Routing**: TanStack Router
- **Forms**: React Hook Form + Zod validation
- **Styling**: CSS with design tokens (forest green, warm off-white)
- **Testing**: Vitest, React Testing Library
- **E2E**: Playwright

### Infrastructure

- **Database**: PostgreSQL 16
- **Containerization**: Docker + Docker Compose
- **Web Server**: Caddy 2 (reverse proxy, automatic HTTPS)
- **CI/CD**: GitHub Actions
- **Error Monitoring**: Sentry
- **Hosting**: VPS (Ubuntu 24.04 LTS)

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 16 (via Docker)
- Git

### Quick Start (Development)

1. **Clone the documentation repository**:
   ```bash
   git clone https://github.com/bradenhooton/pinecone-BTH.git
   cd pinecone-BTH
   ```

2. **Review the documentation**:
   - Start with this README
   - Read `BRD.md` for business requirements
   - Read `TDD.md` for technical architecture
   - Follow `DEVELOPER_ONBOARDING.md` for setup

3. **Set up the database**:
   ```bash
   # Start PostgreSQL
   docker-compose -f docker-compose.dev.yml up -d

   # Verify database is running
   docker ps
   ```

4. **Create implementation repositories**:
   ```bash
   # Backend repository (to be created)
   git clone https://github.com/bradenhooton/pinecone-api.git
   cd pinecone-api
   # Follow setup instructions in DEVELOPER_ONBOARDING.md

   # Frontend repository (to be created)
   git clone https://github.com/bradenhooton/pinecone-web.git
   cd pinecone-web
   # Follow setup instructions in DEVELOPER_ONBOARDING.md
   ```

### Environment Configuration

Copy `.env.example` to `.env.dev` and configure:

```bash
cp .env.example .env.dev
# Edit .env.dev with your settings
```

Key variables:
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: 256-bit secret for JWT tokens
- `USDA_API_KEY`: Get from https://fdc.nal.usda.gov/api-key-signup.html

---

## 📚 Documentation

Comprehensive documentation is available in this repository:

| Document | Description | Audience |
|----------|-------------|----------|
| [BRD.md](BRD.md) | Business Requirements Document | Stakeholders, Product Owners |
| [TDD.md](TDD.md) | Technical Design Document | Developers, Architects |
| [EPIC_BREAKDOWN.md](EPIC_BREAKDOWN.md) | Detailed task breakdown (422 hours) | Developers, Project Managers |
| [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) | 12-week timeline with milestones | All Team Members |
| [DEVELOPER_ONBOARDING.md](DEVELOPER_ONBOARDING.md) | Setup and workflow guide | New Developers |
| [DATABASE_SCHEMA.md](DATABASE_SCHEMA.md) | Complete database reference | Developers, DBAs |
| [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) | Production deployment instructions | DevOps Engineers |
| [DOCUMENTATION_INDEX.md](DOCUMENTATION_INDEX.md) | Documentation overview | All Team Members |

---

## 📁 Project Structure

This is a **documentation-first** repository. Implementation happens in separate repositories:

```
pinecone-project/
├── pinecone-BTH/              # This repository (documentation hub)
│   ├── BRD.md
│   ├── TDD.md
│   ├── EPIC_BREAKDOWN.md
│   ├── PROJECT_ROADMAP.md
│   ├── DATABASE_SCHEMA.md
│   ├── DEVELOPER_ONBOARDING.md
│   ├── DEPLOYMENT_GUIDE.md
│   ├── DOCUMENTATION_INDEX.md
│   ├── .env.example
│   ├── .gitignore
│   ├── config/
│   │   └── grocery_departments.yaml
│   └── docker-compose.dev.yml
│
├── pinecone-api/              # Backend (Go) - To be created
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── api/
│   └── ...
│
└── pinecone-web/              # Frontend (React) - To be created
    ├── src/
    ├── public/
    └── ...
```

---

## 👨‍💻 Development

### Test-Driven Development (TDD)

This project follows strict TDD practices:

1. **RED**: Write a failing test
2. **GREEN**: Write minimum code to pass
3. **REFACTOR**: Improve code quality

**Coverage Targets**:
- Backend: ≥80% (services, repositories)
- Frontend: ≥70% (components, hooks)

### Git Workflow

**Branching Strategy**: GitHub Flow

```bash
# Feature branch
git checkout -b feat/recipe-123/add-search

# Commit frequently with conventional commits
git commit -m "test: add search recipe by title test"
git commit -m "feat: implement search recipe by title"
git commit -m "refactor: extract search logic to helper"

# Push and create PR
git push origin feat/recipe-123/add-search
```

### Running Tests

**Backend**:
```bash
go test ./...                    # All tests
go test -cover ./...             # With coverage
go test -race ./...              # Race detection
```

**Frontend**:
```bash
npm test                         # All tests
npm run test:coverage            # With coverage
```

---

## 🗓 Roadmap

**Total Duration**: 12 weeks (Nov 9, 2025 - Jan 31, 2026)

| Phase | Timeline | Focus |
|-------|----------|-------|
| **Week 1-2** | Nov 9-22 | Foundation & Infrastructure |
| **Week 3** | Nov 23-29 | User Authentication |
| **Week 4-5** | Nov 30-Dec 13 | Recipe Management |
| **Week 6** | Dec 14-20 | Nutrition Integration |
| **Week 7-8** | Dec 21-Jan 3 | Meal Planning |
| **Week 9** | Jan 4-10 | Grocery Lists |
| **Week 10** | Jan 11-17 | Menu Recommendation |
| **Week 11** | Jan 18-24 | Cookbooks |
| **Week 12** | Jan 25-31 | Polish & Launch 🚀 |

See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for detailed milestones and deliverables.

---

## 🤝 Contributing

This is currently a private household project. Contribution guidelines:

1. Follow the TDD workflow
2. Write tests before implementation
3. Follow code style conventions (see DEVELOPER_ONBOARDING.md)
4. Ensure all tests pass before creating PR
5. Update documentation as needed

### Pull Request Process

1. Create feature branch from `main`
2. Write tests (RED)
3. Implement feature (GREEN)
4. Refactor (REFACTOR)
5. Ensure CI passes
6. Request review
7. Squash and merge

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Contact

**Project Owner**: BradenHooton
**GitHub**: https://github.com/bradenhooton

---

## 🙏 Acknowledgments

- **USDA FoodData Central**: Nutrition data API
- **Anthropic Claude**: Documentation generation assistance
- **Open Source Community**: For the excellent tools and libraries

---

## 📊 Project Status

- **Planning**: ✅ Complete
- **Documentation**: ✅ Complete
- **Backend Development**: ⏳ Not Started
- **Frontend Development**: ⏳ Not Started
- **Deployment**: ⏳ Not Started

---

**Ready to build something amazing! 🌲**

For detailed setup instructions, start with [DEVELOPER_ONBOARDING.md](DEVELOPER_ONBOARDING.md).
