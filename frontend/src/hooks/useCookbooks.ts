import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  api,
  CreateCookbookRequest,
  UpdateCookbookRequest,
} from '../lib/api'

// Get all cookbooks
export function useCookbooks(limit?: number, offset?: number) {
  return useQuery({
    queryKey: ['cookbooks', limit, offset],
    queryFn: () => api.getCookbooks(limit, offset),
  })
}

// Get single cookbook by ID
export function useCookbook(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['cookbook', id],
    queryFn: () => api.getCookbookById(id),
    enabled,
  })
}

// Create cookbook
export function useCreateCookbook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateCookbookRequest) => api.createCookbook(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cookbooks'] })
    },
  })
}

// Update cookbook
export function useUpdateCookbook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCookbookRequest }) =>
      api.updateCookbook(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['cookbook', variables.id] })
      queryClient.invalidateQueries({ queryKey: ['cookbooks'] })
    },
  })
}

// Delete cookbook
export function useDeleteCookbook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => api.deleteCookbook(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cookbooks'] })
    },
  })
}

// Add recipe to cookbook
export function useAddRecipeToCookbook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ cookbookId, recipeId }: { cookbookId: string; recipeId: string }) =>
      api.addRecipeToCookbook(cookbookId, recipeId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['cookbook', variables.cookbookId] })
      queryClient.invalidateQueries({ queryKey: ['cookbooks'] })
    },
  })
}

// Remove recipe from cookbook
export function useRemoveRecipeFromCookbook() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ cookbookId, recipeId }: { cookbookId: string; recipeId: string }) =>
      api.removeRecipeFromCookbook(cookbookId, recipeId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['cookbook', variables.cookbookId] })
      queryClient.invalidateQueries({ queryKey: ['cookbooks'] })
    },
  })
}
