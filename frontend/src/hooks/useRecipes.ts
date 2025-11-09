import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  api,
  CreateRecipeRequest,
  RecipeSearchParams,
  UpdateRecipeRequest,
} from '../lib/api'

export function useRecipes(params?: RecipeSearchParams) {
  return useQuery({
    queryKey: ['recipes', params],
    queryFn: () => api.getRecipes(params),
  })
}

export function useRecipe(id: string | undefined) {
  return useQuery({
    queryKey: ['recipe', id],
    queryFn: () => api.getRecipeById(id!),
    enabled: !!id,
  })
}

export function useCreateRecipe() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateRecipeRequest) => api.createRecipe(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['recipes'] })
    },
  })
}

export function useUpdateRecipe() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateRecipeRequest }) =>
      api.updateRecipe(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['recipes'] })
      queryClient.invalidateQueries({ queryKey: ['recipe', variables.id] })
    },
  })
}

export function useDeleteRecipe() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => api.deleteRecipe(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['recipes'] })
    },
  })
}

export function useUploadRecipeImage() {
  return useMutation({
    mutationFn: (file: File) => api.uploadRecipeImage(file),
  })
}
