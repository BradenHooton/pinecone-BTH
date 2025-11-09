import { useMutation } from '@tanstack/react-query'
import { api, RecommendRecipesRequest } from '../lib/api'

// Get recipe recommendations based on ingredients
export function useRecommendRecipes() {
  return useMutation({
    mutationFn: (data: RecommendRecipesRequest) => api.recommendRecipes(data),
  })
}
