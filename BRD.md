# Business Requirements Document (BRD)
## Pinecone Recipe Management & Meal Planning System

**Version:** 1.0  
**Date:** 2025-11-09  
**Status:** Approved  
**Project Owner:** BHooton

---

## Table of Contents
1. [Project Overview](#1-project-overview)
2. [Stakeholders](#2-stakeholders)
3. [Business Requirements](#3-business-requirements)
4. [Non-Functional Requirements](#4-non-functional-requirements)
5. [Assumptions, Risks, and Constraints](#5-assumptions-risks-and-constraints)
6. [Success Metrics](#6-success-metrics)
7. [Document Sign-off](#7-document-sign-off)

---

## 1. Project Overview

### 1.1 Business Problem
A household lacks a unified system to manage recipes, plan meals for the week, and automatically generate grocery lists. This results in:
- Time wasted searching for recipes across scattered sources (bookmarks, notebooks, websites)
- Inefficient grocery shopping (multiple trips, forgotten ingredients, duplicate purchases)
- Difficulty answering "What's for dinner?" without manual planning effort

**Solution:** Pinecone provides a single-household recipe database, visual meal planning calendar, and auto-generated grocery lists organized by store department.

### 1.2 Project Objectives
1. **Centralize Recipe Storage:** Consolidate all household recipes in a searchable, filterable database with rich metadata (nutrition, prep time, tags).
2. **Streamline Meal Planning:** Enable drag-and-drop meal planning across a calendar with multiple recipes per meal slot, including future planning and "Out of Kitchen" options.
3. **Automate Grocery Lists:** Generate a grocery list from the meal plan, summing ingredient quantities and grouping by store department.
4. **Reduce Decision Fatigue:** Provide an ingredient-based menu recommendation feature to suggest recipes based on available ingredients.
5. **Organize Recipes:** Allow users to create Cookbooks (collections) for themed grouping (e.g., "Holiday Recipes," "High Protein Meals").

### 1.3 Scope

#### In Scope (MVP - Phase 1)
- ✅ User authentication (email/password, JWT-based sessions)
- ✅ Recipe CRUD (Create, Read, Update, Delete) with full metadata
- ✅ Recipe image upload (≤5MB) or URL reference
- ✅ Recipe search, filter (by tags, prep time, nutrition), and sort (alphabetical, date added)
- ✅ Nutrition data integration via USDA FoodData Central API
- ✅ Meal planning calendar with future date support and "Out of Kitchen" meal option
- ✅ 5 meal slots per day: Breakfast, Lunch, Snack, Dinner, Dessert
- ✅ Grocery list generation for a date range (aggregates ingredients, groups by store department)
- ✅ Grocery list interaction (mark "Bought" or "Have on Hand," add manual items)
- ✅ Ingredient-based menu recommendation (partial match, ranked results)
- ✅ Cookbook creation and recipe assignment (many-to-many relationship)
- ✅ Design token system (forest green primary, warm off-white background, classic typography)
- ✅ Test-Driven Development (TDD) workflow for all backend and frontend logic

#### Out of Scope (Future Phases)
- ❌ Multi-household/multi-tenancy
- ❌ Public recipe sharing (social features)
- ❌ Native mobile apps (iOS/Android)
- ❌ Integration with grocery delivery services (Instacart, Amazon Fresh)
- ❌ Automatic recipe import from arbitrary URLs (web scraping)
- ❌ Calorie/macro tracking per user
- ❌ Dietary restriction filters (e.g., "gluten-free," "vegan") — can be added as tags in MVP

---

## 2. Stakeholders

| Name | Role | Responsibility |
|------|------|----------------|
| BHooton | Project Owner / Developer | Defines requirements, builds system, validates functionality |
| Household Members | End Users | Test usability, provide feedback on meal planning workflow |

---

## 3. Business Requirements

### 3.1 Functional Requirements

#### FR-1: User Management
- **FR-1.1:** The system must support user registration with email and password.
- **FR-1.2:** Passwords must be hashed using bcrypt before storage.
- **FR-1.3:** The system must authenticate users via JWT tokens (HTTP-only cookies, 24-hour expiration).
- **FR-1.4:** All users in the household have equal permissions (no role-based access control in MVP).
- **FR-1.5:** Users must be able to log out, invalidating their session token.

#### FR-2: Recipe Management
- **FR-2.1:** Users must be able to create a recipe with the following fields:
  - Title (required, max 200 characters)
  - Image (optional: file upload ≤5MB OR URL reference)
  - Nutrition Details (optional: calories, protein, fiber, carbs, fats — sourced from USDA API)
  - Servings (required, integer)
  - Serving Size (required, e.g., "1 cup," "2 slices")
  - Prep Time (optional, minutes)
  - Cook Time (optional, minutes)
  - Total Time (auto-calculated: Prep + Cook)
  - Ingredients (required, list of `{ingredient_name, quantity, unit, nutrition_id}`)
  - Instructions (required, ordered list of text steps)
  - Storage/Freezing Notes (optional, text)
  - Tags (optional, array of strings, e.g., "vegetarian," "quick," "comfort food")
  - Source (optional, e.g., "Grandma's cookbook," URL)
  - Notes (optional, free-form text)
- **FR-2.2:** Users must be able to edit and delete recipes they created.
- **FR-2.3:** Recipes must be soft-deleted (flagged as `deleted_at` in DB, not hard-deleted).
- **FR-2.4:** The system must display a grid/list view of recipes with:
  - Search (by title, ingredient, tag)
  - Filter (by tags, prep time range, nutrition ranges)
  - Sort (alphabetical, date added, prep time)

#### FR-3: Nutrition Data Integration
- **FR-3.1:** The system must integrate with the USDA FoodData Central API to fetch nutrition data for ingredients.
- **FR-3.2:** When a user adds an ingredient, the system must search USDA API and allow the user to select a matching food item.
- **FR-3.3:** The system must cache USDA nutrition data in the database to reduce API calls.
- **FR-3.4:** Recipe-level nutrition must be auto-calculated by summing ingredient nutrition and dividing by servings.

#### FR-4: Meal Planning
- **FR-4.1:** The system must provide a meal planning calendar that supports any future date.
- **FR-4.2:** Each day must have 5 meal slots: Breakfast, Lunch, Snack, Dinner, Dessert.
- **FR-4.3:** Users must be able to add multiple recipes to a single meal slot.
- **FR-4.4:** When adding a recipe to a meal slot, users must specify the number of servings.
- **FR-4.5:** Users must be able to remove recipes from meal slots.
- **FR-4.6:** The meal plan is shared by all household users (single source of truth).
- **FR-4.7:** Users must be able to navigate to any future date to plan meals in advance (e.g., plan next month's holiday meals).
- **FR-4.8:** Users must be able to mark a meal slot as "Out of Kitchen" to indicate the household is eating out, skipping the meal, or otherwise not cooking.

#### FR-5: Grocery List Generation
- **FR-5.1:** Users must be able to generate a grocery list for a selected date range (e.g., "Next 7 days").
- **FR-5.2:** The system must aggregate all recipes in the meal plan for the selected range.
- **FR-5.3:** Meal slots marked as "Out of Kitchen" must be excluded from grocery list calculations.
- **FR-5.4:** The system must sum ingredient quantities across recipes (e.g., "2 eggs" + "3 eggs" = "5 eggs").
- **FR-5.5:** Ingredients must be grouped by grocery store department (e.g., "Produce," "Dairy," "Meat").
- **FR-5.6:** The system must maintain a YAML configuration file (`grocery_departments.yaml`) with a predefined list of departments.
- **FR-5.7:** Each ingredient in the system must be tagged with a department (enum value).
- **FR-5.8:** Users must be able to mark grocery list items as "Bought" or "Have on Hand."
- **FR-5.9:** Users must be able to manually add free-form items to the grocery list (e.g., "Paper towels").
- **FR-5.10:** Manual items must not be linked to recipes or meal plans.

#### FR-6: Ingredient-Based Menu Recommendation
- **FR-6.1:** Users must be able to input a list of ingredients they have on hand.
- **FR-6.2:** The system must return a ranked list of recipes, ordered by:
  - **Match Score** = (Number of matching ingredients) / (Total ingredients in recipe)
- **FR-6.3:** Recipes with partial matches must be included (e.g., if recipe needs 5 ingredients and user has 3, it shows up with a 60% match).
- **FR-6.4:** The results must be displayed in a "restaurant menu" style layout:
  - Classic French aesthetic (elegant serif typography, minimal decoration)
  - Recipe title, brief description, match score, missing ingredients
- **FR-6.5:** Users must be able to click a recipe in the menu to view full details.

#### FR-7: Cookbooks (Recipe Collections)
- **FR-7.1:** Users must be able to create Cookbooks with a name and optional description.
- **FR-7.2:** Cookbooks are public to all household users but track the creator (ownership).
- **FR-7.3:** A recipe can belong to multiple Cookbooks (many-to-many relationship).
- **FR-7.4:** Users must be able to add/remove recipes from Cookbooks.
- **FR-7.5:** Users must be able to delete Cookbooks (soft-delete, preserves recipes).

### 3.2 User Stories (Epics)

#### Epic 1: User Authentication
- **US-1.1:** As a user, I want to register an account with email and password so I can access the system.
- **US-1.2:** As a user, I want to log in with my credentials so I can view and manage recipes.
- **US-1.3:** As a user, I want to log out so my session is securely ended.

#### Epic 2: Recipe Management
- **US-2.1:** As a user, I want to create a recipe with all metadata so I can store my favorite dishes.
- **US-2.2:** As a user, I want to upload a recipe image or provide an image URL so I can visually identify recipes.
- **US-2.3:** As a user, I want to edit a recipe so I can update ingredients or instructions.
- **US-2.4:** As a user, I want to delete a recipe so I can remove dishes I no longer cook.
- **US-2.5:** As a user, I want to search and filter recipes so I can quickly find what I'm looking for.

#### Epic 3: Nutrition Data
- **US-3.1:** As a user, I want to search for ingredient nutrition data so I can see accurate calorie/macro information.
- **US-3.2:** As a user, I want the system to auto-calculate recipe nutrition so I don't do math manually.

#### Epic 4: Meal Planning
- **US-4.1:** As a user, I want to view a meal planning calendar so I can plan the week's meals.
- **US-4.2:** As a user, I want to add recipes to specific meal slots so I can schedule when to cook each dish.
- **US-4.3:** As a user, I want to specify serving counts when adding recipes so the grocery list scales correctly.
- **US-4.4:** As a user, I want to navigate to future dates on the calendar so I can plan meals weeks or months in advance.
- **US-4.5:** As a user, I want to mark a meal slot as "Out of Kitchen" so the system knows I'm not cooking that meal and excludes it from grocery lists.

#### Epic 5: Grocery List
- **US-5.1:** As a user, I want to generate a grocery list for a date range so I know what to buy.
- **US-5.2:** As a user, I want ingredients grouped by store department so I can shop efficiently.
- **US-5.3:** As a user, I want to mark items as "Bought" or "Have on Hand" so I track my shopping progress.
- **US-5.4:** As a user, I want to add manual items to the grocery list so I can include non-recipe purchases.

#### Epic 6: Ingredient Recommendation
- **US-6.1:** As a user, I want to input ingredients I have so the system suggests recipes I can make.
- **US-6.2:** As a user, I want to see a ranked list with match scores so I know which recipes are easiest to prepare.
- **US-6.3:** As a user, I want the menu styled elegantly so the experience feels refined.

#### Epic 7: Cookbooks
- **US-7.1:** As a user, I want to create a Cookbook so I can group related recipes.
- **US-7.2:** As a user, I want to add recipes to multiple Cookbooks so I can organize by theme.
- **US-7.3:** As a user, I want to view all recipes in a Cookbook so I can browse themed collections.

---

## 4. Non-Functional Requirements (NFRs)

### 4.1 Performance
- **NFR-4.1.1:** API endpoints must respond in < 500ms under normal load (2-6 concurrent users).
- **NFR-4.1.2:** Recipe search must return results in < 300ms for a database of up to 500 recipes.
- **NFR-4.1.3:** Grocery list generation must complete in < 2 seconds for a 7-day meal plan with 20 recipes.
- **NFR-4.1.4:** Image uploads must complete in < 5 seconds for files up to 5MB.

### 4.2 Security
- **NFR-4.2.1:** All passwords must be hashed with bcrypt (cost factor ≥12).
- **NFR-4.2.2:** JWT tokens must be stored in HTTP-only cookies (not localStorage).
- **NFR-4.2.3:** All API requests (except `/auth/login` and `/auth/register`) must require a valid JWT.
- **NFR-4.2.4:** Uploaded images must be scanned for file type (allow only `.jpg`, `.jpeg`, `.png`, `.webp`).
- **NFR-4.2.5:** All database queries must use parameterized statements (sqlc) to prevent SQL injection.
- **NFR-4.2.6:** HTTPS must be enforced in production (Caddy auto-TLS).

### 4.3 Scalability
- **NFR-4.3.1:** The database schema must support up to 1,000 recipes without performance degradation.
- **NFR-4.3.2:** The system must handle 6 concurrent users without resource contention.

### 4.4 Availability
- **NFR-4.4.1:** The system must maintain 99% uptime (excluding planned maintenance).
- **NFR-4.4.2:** Database backups must be automated (daily, retained for 7 days).

### 4.5 Usability
- **NFR-4.5.1:** The UI must be responsive (mobile, tablet, desktop).
- **NFR-4.5.2:** The design must use the defined color palette (forest green primary, warm off-white background, black accents).
- **NFR-4.5.3:** Typography must evoke a classic, elegant aesthetic (serif fonts for headings, sans-serif for body).

### 4.6 Development Process (Test-Driven Development)
- **NFR-4.6.1:** All backend functions (handlers, services, repositories) must have unit tests written BEFORE implementation.
- **NFR-4.6.2:** All frontend components and hooks must have tests written BEFORE implementation.
- **NFR-4.6.3:** Tests must validate the expected behavior as defined in Acceptance Criteria.
- **NFR-4.6.4:** The development workflow must follow Red-Green-Refactor:
  1. **Red:** Write a failing test that defines the desired behavior.
  2. **Green:** Write the minimum code necessary to make the test pass.
  3. **Refactor:** Improve code quality while keeping tests green.
- **NFR-4.6.5:** All PRs must include evidence that tests were written first (e.g., test file timestamps, commit history).
- **NFR-4.6.6:** Code coverage targets:
  - Backend: Minimum 80% coverage for service and repository layers.
  - Frontend: Minimum 70% coverage for components and hooks.

---

## 5. Assumptions, Risks, and Constraints

### 5.1 Assumptions
- **A-1:** Users have stable internet access to interact with the web application.
- **A-2:** The USDA FoodData Central API remains free and available (no rate limits exceeded).
- **A-3:** Users will manually curate recipe data (no automated web scraping in MVP).
- **A-4:** The household size remains 2-6 users (no need for horizontal scaling).
- **A-5:** Image uploads are stored on the server filesystem (no cloud object storage in MVP).
- **A-6:** Developers are proficient in Test-Driven Development or willing to learn the discipline.

### 5.2 Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| USDA API downtime or deprecation | Medium | High | Cache nutrition data locally; fallback to manual entry |
| Users exceed 5MB image size | Low | Low | Client-side validation; server-side rejection with clear error |
| Database query performance degrades with 500+ recipes | Low | Medium | Index on `title`, `tags`, `created_at`; optimize search queries |
| Users confuse "Bought" vs. "Have on Hand" | Medium | Low | Clear UI labels; tooltips explaining difference |
| TDD discipline not followed consistently | Medium | Medium | Enforce via PR reviews; automated coverage checks in CI |
| Future meal planning creates UI complexity | Low | Medium | Implement infinite scroll or pagination for calendar navigation |

### 5.3 Constraints
- **C-1:** Must use the mandated tech stack (Go, React, PostgreSQL, no exceptions).
- **C-2:** Must host on a single VPS (no cloud services like AWS S3, RDS).
- **C-3:** Budget limited to hosting costs (~$10-20/month).
- **C-4:** No third-party recipe import in MVP (manual entry only).
- **C-5:** All code must be test-driven (tests written before implementation).

---

## 6. Success Metrics

### 6.1 Key Performance Indicators (KPIs)

| Metric | Target | Measurement Method |
|--------|--------|-------------------|
| **Weekly Meal Plan Completion Rate** | 80% within 3 months | % of weeks where a full 7-day plan is created |
| **Recipe Database Size** | 100 recipes within 6 months | Count of active recipes in DB |
| **Grocery List Usage** | 90% of meal plans generate a grocery list | % of meal plans with associated grocery list |
| **Average Recipe Search Time** | < 10 seconds | Time from search input to result click |
| **User Satisfaction** | 4/5 stars | Post-use survey (informal, household feedback) |
| **Code Coverage** | Backend ≥80%, Frontend ≥70% | Automated coverage reports in CI pipeline |
| **"Out of Kitchen" Adoption** | 50% of meal plans use this feature | % of meal slots marked "Out of Kitchen" |

---

## 7. Document Sign-off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Project Owner | BHooton | ✅ Approved | 2025-11-09 |

---

**Document Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-09 | GhostDev | Initial BRD creation and approval |
