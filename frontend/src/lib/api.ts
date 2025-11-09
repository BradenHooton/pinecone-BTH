/**
 * API Client for Pinecone Backend
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

export interface User {
  id: string
  email: string
  name: string
  created_at: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  data: {
    user: User
  }
  meta: {
    timestamp: string
  }
}

export interface ErrorResponse {
  error: {
    message: string
  }
  meta: {
    timestamp: string
  }
}

export interface Recipe {
  id: string
  created_by_user_id: string
  title: string
  image_url?: string
  servings: number
  serving_size: string
  prep_time_minutes?: number
  cook_time_minutes?: number
  total_time_minutes: number
  storage_notes?: string
  source?: string
  notes?: string
  created_at: string
  updated_at: string
  ingredients?: RecipeIngredient[]
  instructions?: RecipeInstruction[]
  tags?: RecipeTag[]
}

export interface RecipeIngredient {
  id: string
  recipe_id: string
  nutrition_id?: string
  ingredient_name: string
  quantity: number
  unit: string
  department: string
  order_index: number
}

export interface RecipeInstruction {
  id: string
  recipe_id: string
  step_number: number
  instruction: string
}

export interface RecipeTag {
  id: string
  recipe_id: string
  tag_name: string
}

export interface CreateIngredientRequest {
  ingredient_name: string
  quantity: number
  unit: string
  department: string
  nutrition_id?: string
}

export interface CreateInstructionRequest {
  step_number: number
  instruction: string
}

export interface CreateRecipeRequest {
  title: string
  image_url?: string
  servings: number
  serving_size: string
  prep_time_minutes?: number
  cook_time_minutes?: number
  storage_notes?: string
  source?: string
  notes?: string
  ingredients: CreateIngredientRequest[]
  instructions: CreateInstructionRequest[]
  tags?: string[]
}

export interface UpdateRecipeRequest extends CreateRecipeRequest {}

export interface RecipeListResponse {
  data: Recipe[]
  meta: {
    total: number
    limit: number
    offset: number
  }
}

export interface RecipeResponse {
  data: Recipe
}

export interface RecipeSearchParams {
  search?: string
  tags?: string
  sort?: string
  limit?: number
  offset?: number
}

export interface NutritionSearchResult {
  fdc_id: string
  description: string
  data_type: string
  calories?: number
  protein_g?: number
  carbs_g?: number
  fiber_g?: number
  fat_g?: number
}

export interface NutritionSearchResponse {
  data: NutritionSearchResult[]
  meta: {
    total: number
  }
}

export type MealType = 'breakfast' | 'lunch' | 'snack' | 'dinner' | 'dessert'

export interface MealPlanRecipe {
  id: string
  meal_plan_id: string
  meal_type: MealType
  recipe_id?: string
  recipe?: Recipe
  servings?: number
  out_of_kitchen: boolean
  order_index: number
}

export interface MealPlan {
  id: string
  plan_date: string
  created_at: string
  updated_at: string
  meals: MealPlanRecipe[]
}

export interface CreateMealPlanRecipeRequest {
  meal_type: MealType
  recipe_id?: string
  servings?: number
  out_of_kitchen: boolean
}

export interface UpdateMealPlanRequest {
  meals: CreateMealPlanRecipeRequest[]
}

export interface MealPlanResponse {
  data: MealPlan
}

export interface MealPlanListResponse {
  data: MealPlan[]
  meta: {
    start_date: string
    end_date: string
  }
}

export type GroceryDepartment =
  | 'produce'
  | 'meat'
  | 'seafood'
  | 'dairy'
  | 'bakery'
  | 'frozen'
  | 'pantry'
  | 'spices'
  | 'beverages'
  | 'other'

export type GroceryItemStatus = 'pending' | 'bought' | 'have_on_hand'

export interface GroceryListItem {
  id: string
  grocery_list_id: string
  item_name: string
  quantity?: number
  unit?: string
  department: GroceryDepartment
  status: GroceryItemStatus
  is_manual: boolean
  source_recipe_id?: string
  source_recipe?: Recipe
}

export interface GroceryList {
  id: string
  created_by_user_id: string
  start_date: string
  end_date: string
  items: GroceryListItem[]
  created_at: string
  updated_at: string
}

export interface CreateGroceryListRequest {
  start_date: string
  end_date: string
}

export interface CreateManualItemRequest {
  item_name: string
  quantity?: number
  unit?: string
  department?: GroceryDepartment
}

export interface UpdateItemStatusRequest {
  status: GroceryItemStatus
}

export interface GroceryListResponse {
  data: GroceryList
}

export interface GroceryListListResponse {
  data: GroceryList[]
  meta: {
    total: number
  }
}

export interface RecommendRecipesRequest {
  ingredients: string[]
}

export interface RecipeRecommendation {
  recipe: Recipe
  match_score: number
  matched_ingredients: string[]
  missing_ingredients: string[]
}

export interface RecommendRecipesResponse {
  data: RecipeRecommendation[]
  meta: {
    provided_ingredients: string[]
    total_recipes_found: number
  }
}

export interface Cookbook {
  id: string
  created_by_user_id: string
  name: string
  description?: string
  recipe_count: number
  recipes?: Recipe[]
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface CreateCookbookRequest {
  name: string
  description?: string
}

export interface UpdateCookbookRequest {
  name: string
  description?: string
}

export interface CookbookResponse {
  data: Cookbook
}

export interface CookbookListResponse {
  data: Cookbook[]
  meta: {
    total: number
  }
}

class ApiClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    const response = await fetch(url, {
      ...options,
      credentials: 'include', // Send cookies
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      const error: ErrorResponse = await response.json()
      throw new Error(error.error.message || 'Request failed')
    }

    return response.json()
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async logout(): Promise<void> {
    await this.request('/auth/logout', {
      method: 'POST',
    })
  }

  // Recipe endpoints
  async getRecipes(params?: RecipeSearchParams): Promise<RecipeListResponse> {
    const queryParams = new URLSearchParams()
    if (params?.search) queryParams.append('search', params.search)
    if (params?.tags) queryParams.append('tags', params.tags)
    if (params?.sort) queryParams.append('sort', params.sort)
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.offset) queryParams.append('offset', params.offset.toString())

    const query = queryParams.toString()
    const endpoint = query ? `/recipes?${query}` : '/recipes'

    return this.request<RecipeListResponse>(endpoint)
  }

  async getRecipeById(id: string): Promise<RecipeResponse> {
    return this.request<RecipeResponse>(`/recipes/${id}`)
  }

  async createRecipe(data: CreateRecipeRequest): Promise<RecipeResponse> {
    return this.request<RecipeResponse>('/recipes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateRecipe(id: string, data: UpdateRecipeRequest): Promise<RecipeResponse> {
    return this.request<RecipeResponse>(`/recipes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteRecipe(id: string): Promise<void> {
    await this.request(`/recipes/${id}`, {
      method: 'DELETE',
    })
  }

  async uploadRecipeImage(file: File): Promise<{ data: { image_url: string } }> {
    const formData = new FormData()
    formData.append('image', file)

    const url = `${this.baseURL}/recipes/upload-image`
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    })

    if (!response.ok) {
      const error: ErrorResponse = await response.json()
      throw new Error(error.error.message || 'Upload failed')
    }

    return response.json()
  }

  // Nutrition endpoints
  async searchNutrition(query: string): Promise<NutritionSearchResponse> {
    const encodedQuery = encodeURIComponent(query)
    return this.request<NutritionSearchResponse>(`/nutrition/search?query=${encodedQuery}`)
  }

  // Meal plan endpoints
  async getMealPlanByDate(date: string): Promise<MealPlanResponse> {
    return this.request<MealPlanResponse>(`/meal-plans/date?date=${date}`)
  }

  async getMealPlansByDateRange(startDate: string, endDate: string): Promise<MealPlanListResponse> {
    return this.request<MealPlanListResponse>(`/meal-plans?start_date=${startDate}&end_date=${endDate}`)
  }

  async updateMealPlan(date: string, data: UpdateMealPlanRequest): Promise<MealPlanResponse> {
    return this.request<MealPlanResponse>(`/meal-plans/date?date=${date}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // Grocery list endpoints
  async getGroceryLists(limit?: number, offset?: number): Promise<GroceryListListResponse> {
    const queryParams = new URLSearchParams()
    if (limit) queryParams.append('limit', limit.toString())
    if (offset) queryParams.append('offset', offset.toString())

    const query = queryParams.toString()
    const endpoint = query ? `/grocery-lists?${query}` : '/grocery-lists'

    return this.request<GroceryListListResponse>(endpoint)
  }

  async getGroceryListById(id: string): Promise<GroceryListResponse> {
    return this.request<GroceryListResponse>(`/grocery-lists/${id}`)
  }

  async createGroceryList(data: CreateGroceryListRequest): Promise<GroceryListResponse> {
    return this.request<GroceryListResponse>('/grocery-lists', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async deleteGroceryList(id: string): Promise<void> {
    await this.request(`/grocery-lists/${id}`, {
      method: 'DELETE',
    })
  }

  async addManualItem(listId: string, data: CreateManualItemRequest): Promise<{ data: GroceryListItem }> {
    return this.request<{ data: GroceryListItem }>(`/grocery-lists/${listId}/items`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateItemStatus(itemId: string, data: UpdateItemStatusRequest): Promise<void> {
    await this.request(`/grocery-lists/items/${itemId}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  // Menu recommendation endpoints
  async recommendRecipes(data: RecommendRecipesRequest): Promise<RecommendRecipesResponse> {
    return this.request<RecommendRecipesResponse>('/menu/recommend', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Cookbook endpoints
  async getCookbooks(limit?: number, offset?: number): Promise<CookbookListResponse> {
    const queryParams = new URLSearchParams()
    if (limit) queryParams.append('limit', limit.toString())
    if (offset) queryParams.append('offset', offset.toString())

    const query = queryParams.toString()
    const endpoint = query ? `/cookbooks?${query}` : '/cookbooks'

    return this.request<CookbookListResponse>(endpoint)
  }

  async getCookbookById(id: string): Promise<CookbookResponse> {
    return this.request<CookbookResponse>(`/cookbooks/${id}`)
  }

  async createCookbook(data: CreateCookbookRequest): Promise<CookbookResponse> {
    return this.request<CookbookResponse>('/cookbooks', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateCookbook(id: string, data: UpdateCookbookRequest): Promise<CookbookResponse> {
    return this.request<CookbookResponse>(`/cookbooks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteCookbook(id: string): Promise<void> {
    await this.request(`/cookbooks/${id}`, {
      method: 'DELETE',
    })
  }

  async addRecipeToCookbook(cookbookId: string, recipeId: string): Promise<void> {
    await this.request(`/cookbooks/${cookbookId}/recipes/${recipeId}`, {
      method: 'POST',
    })
  }

  async removeRecipeFromCookbook(cookbookId: string, recipeId: string): Promise<void> {
    await this.request(`/cookbooks/${cookbookId}/recipes/${recipeId}`, {
      method: 'DELETE',
    })
  }
}

export const api = new ApiClient(API_BASE_URL)
