import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  api,
  CreateGroceryListRequest,
  CreateManualItemRequest,
  UpdateItemStatusRequest,
} from '../lib/api'

// Get all grocery lists
export function useGroceryLists(limit?: number, offset?: number) {
  return useQuery({
    queryKey: ['groceryLists', limit, offset],
    queryFn: () => api.getGroceryLists(limit, offset),
  })
}

// Get single grocery list by ID
export function useGroceryList(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['groceryList', id],
    queryFn: () => api.getGroceryListById(id),
    enabled,
  })
}

// Create grocery list
export function useCreateGroceryList() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateGroceryListRequest) => api.createGroceryList(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groceryLists'] })
    },
  })
}

// Delete grocery list
export function useDeleteGroceryList() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => api.deleteGroceryList(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groceryLists'] })
    },
  })
}

// Add manual item
export function useAddManualItem() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ listId, data }: { listId: string; data: CreateManualItemRequest }) =>
      api.addManualItem(listId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['groceryList', variables.listId] })
      queryClient.invalidateQueries({ queryKey: ['groceryLists'] })
    },
  })
}

// Update item status
export function useUpdateItemStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ itemId, data }: { itemId: string; data: UpdateItemStatusRequest }) =>
      api.updateItemStatus(itemId, data),
    onSuccess: () => {
      // Invalidate all grocery list queries since we don't know which list the item belongs to
      queryClient.invalidateQueries({ queryKey: ['groceryList'] })
      queryClient.invalidateQueries({ queryKey: ['groceryLists'] })
    },
  })
}
