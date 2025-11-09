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
}

export const api = new ApiClient(API_BASE_URL)
