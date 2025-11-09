# Pinecone API Documentation

Base URL: `http://localhost:8080/api/v1`

All endpoints except `/auth/register` and `/auth/login` require authentication via JWT cookie.

## Authentication

### Register User
**POST** `/auth/register`

Create a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123",
  "name": "John Doe"
}
```

**Response (201):**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-11-09T12:00:00Z"
    }
  },
  "meta": {
    "timestamp": "2025-11-09T12:00:00Z"
  }
}
```

### Login
**POST** `/auth/login`

Login and receive JWT cookie.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-11-09T12:00:00Z"
    }
  },
  "meta": {
    "timestamp": "2025-11-09T12:00:00Z"
  }
}
```

Sets HTTP-only cookie: `token=<jwt>; HttpOnly; Secure; SameSite=Strict; Max-Age=86400`

### Logout
**POST** `/auth/logout` 🔒

Logout and clear JWT cookie.

**Response (204):** No content

---

## Recipes

### List Recipes
**GET** `/recipes?search={query}&tags={tags}&sort={field}&limit={n}&offset={n}` 🔒

Get paginated list of recipes.

**Query Parameters:**
- `search` (optional): Search in title
- `tags` (optional): Filter by tags (comma-separated)
- `sort` (optional): Sort field (created_at, title)
- `limit` (optional): Items per page (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "created_by_user_id": "uuid",
      "title": "Chocolate Chip Cookies",
      "image_url": "https://example.com/image.jpg",
      "servings": 24,
      "serving_size": "1 cookie",
      "prep_time_minutes": 15,
      "cook_time_minutes": 12,
      "total_time_minutes": 27,
      "ingredients": [...],
      "instructions": [...],
      "tags": [...],
      "created_at": "2025-11-09T12:00:00Z",
      "updated_at": "2025-11-09T12:00:00Z"
    }
  ],
  "meta": {
    "total": 100,
    "limit": 20,
    "offset": 0
  }
}
```

### Get Recipe
**GET** `/recipes/{id}` 🔒

Get single recipe by ID.

**Response (200):** Same as list item above

### Create Recipe
**POST** `/recipes` 🔒

Create a new recipe.

**Request:**
```json
{
  "title": "Chocolate Chip Cookies",
  "servings": 24,
  "serving_size": "1 cookie",
  "prep_time_minutes": 15,
  "cook_time_minutes": 12,
  "ingredients": [
    {
      "ingredient_name": "flour",
      "quantity": 2.5,
      "unit": "cups",
      "department": "pantry",
      "order_index": 0
    }
  ],
  "instructions": [
    {
      "step_number": 1,
      "instruction": "Preheat oven to 350°F"
    }
  ],
  "tags": ["dessert", "baking"]
}
```

**Response (201):** Created recipe object

### Update Recipe
**PUT** `/recipes/{id}` 🔒

Update existing recipe (full replacement).

**Request:** Same as Create Recipe

**Response (200):** Updated recipe object

### Delete Recipe
**DELETE** `/recipes/{id}` 🔒

Soft delete recipe.

**Response (204):** No content

### Upload Image
**POST** `/recipes/upload-image` 🔒

Upload recipe image.

**Request:** `multipart/form-data` with `image` field

**Constraints:**
- Max size: 5MB
- Allowed types: jpg, png, webp

**Response (200):**
```json
{
  "data": {
    "image_url": "/uploads/uuid.jpg"
  }
}
```

---

## Nutrition

### Search Nutrition
**GET** `/nutrition/search?query={query}` 🔒

Search for nutrition data.

**Response (200):**
```json
{
  "data": [
    {
      "fdc_id": "123456",
      "description": "Chicken breast, raw",
      "data_type": "sr_legacy_food",
      "calories": 120,
      "protein_g": 22.5,
      "carbs_g": 0,
      "fiber_g": 0,
      "fat_g": 2.6
    }
  ],
  "meta": {
    "total": 10
  }
}
```

---

## Meal Plans

### Get Meal Plans by Date Range
**GET** `/meal-plans?start_date={YYYY-MM-DD}&end_date={YYYY-MM-DD}` 🔒

Get meal plans for date range.

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "plan_date": "2025-11-09",
      "created_at": "2025-11-09T12:00:00Z",
      "updated_at": "2025-11-09T12:00:00Z",
      "meals": [...]
    }
  ],
  "meta": {
    "start_date": "2025-11-09",
    "end_date": "2025-11-15"
  }
}
```

### Get Meal Plan by Date
**GET** `/meal-plans/date?date={YYYY-MM-DD}` 🔒

Get meal plan for specific date.

**Response (200):** Single meal plan object

### Update Meal Plan
**PUT** `/meal-plans/date?date={YYYY-MM-DD}` 🔒

Update or create meal plan for specific date.

**Request:**
```json
{
  "meals": [
    {
      "meal_type": "breakfast",
      "recipe_id": "uuid",
      "servings": 4,
      "out_of_kitchen": false
    },
    {
      "meal_type": "lunch",
      "out_of_kitchen": true
    }
  ]
}
```

**Response (200):** Updated meal plan object

---

## Grocery Lists

### List Grocery Lists
**GET** `/grocery-lists?limit={n}&offset={n}` 🔒

Get paginated list of grocery lists.

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "created_by_user_id": "uuid",
      "start_date": "2025-11-09",
      "end_date": "2025-11-15",
      "items": [...],
      "created_at": "2025-11-09T12:00:00Z",
      "updated_at": "2025-11-09T12:00:00Z"
    }
  ],
  "meta": {
    "total": 5
  }
}
```

### Create Grocery List
**POST** `/grocery-lists` 🔒

Generate grocery list from meal plan date range.

**Request:**
```json
{
  "start_date": "2025-11-09",
  "end_date": "2025-11-15"
}
```

**Response (201):** Created grocery list with aggregated items

### Get Grocery List
**GET** `/grocery-lists/{id}` 🔒

Get single grocery list.

**Response (200):** Grocery list object

### Delete Grocery List
**DELETE** `/grocery-lists/{id}` 🔒

Delete grocery list.

**Response (204):** No content

### Add Manual Item
**POST** `/grocery-lists/{id}/items` 🔒

Add manual item to grocery list.

**Request:**
```json
{
  "item_name": "Paper towels",
  "quantity": 2,
  "unit": "rolls",
  "department": "other"
}
```

**Response (201):** Created item object

### Update Item Status
**PATCH** `/grocery-lists/items/{item_id}` 🔒

Update grocery item status.

**Request:**
```json
{
  "status": "bought"
}
```

**Statuses:** `pending`, `bought`, `have_on_hand`

**Response (204):** No content

---

## Menu Recommendations

### Recommend Recipes
**POST** `/menu/recommend` 🔒

Get recipe recommendations based on ingredients.

**Request:**
```json
{
  "ingredients": ["chicken", "rice", "tomatoes"]
}
```

**Response (200):**
```json
{
  "data": [
    {
      "recipe": {...},
      "match_score": 85.5,
      "matched_ingredients": ["chicken", "rice", "tomatoes"],
      "missing_ingredients": ["onion", "garlic"]
    }
  ],
  "meta": {
    "provided_ingredients": ["chicken", "rice", "tomatoes"],
    "total_recipes_found": 5
  }
}
```

---

## Cookbooks

### List Cookbooks
**GET** `/cookbooks?limit={n}&offset={n}` 🔒

Get paginated list of cookbooks.

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "created_by_user_id": "uuid",
      "name": "Holiday Favorites",
      "description": "Recipes for special occasions",
      "recipe_count": 12,
      "created_at": "2025-11-09T12:00:00Z",
      "updated_at": "2025-11-09T12:00:00Z"
    }
  ],
  "meta": {
    "total": 3
  }
}
```

### Get Cookbook
**GET** `/cookbooks/{id}` 🔒

Get cookbook with all recipes.

**Response (200):** Cookbook object with full recipe details

### Create Cookbook
**POST** `/cookbooks` 🔒

Create new cookbook.

**Request:**
```json
{
  "name": "Holiday Favorites",
  "description": "Recipes for special occasions"
}
```

**Response (201):** Created cookbook object

### Update Cookbook
**PUT** `/cookbooks/{id}` 🔒

Update cookbook details.

**Request:** Same as Create

**Response (200):** Updated cookbook object

### Delete Cookbook
**DELETE** `/cookbooks/{id}` 🔒

Soft delete cookbook.

**Response (204):** No content

### Add Recipe to Cookbook
**POST** `/cookbooks/{cookbook_id}/recipes/{recipe_id}` 🔒

Add recipe to cookbook.

**Response (204):** No content

### Remove Recipe from Cookbook
**DELETE** `/cookbooks/{cookbook_id}/recipes/{recipe_id}` 🔒

Remove recipe from cookbook.

**Response (204):** No content

---

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "message": "Human-readable error message"
  }
}
```

### HTTP Status Codes
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Not authenticated
- `403 Forbidden` - Not authorized
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Server error

---

🔒 = Requires authentication (JWT cookie)
