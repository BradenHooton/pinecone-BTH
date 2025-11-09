import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, UpdateMealPlanRequest } from '../lib/api'

export function useMealPlanByDate(date: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['mealPlan', date],
    queryFn: () => api.getMealPlanByDate(date),
    enabled,
  })
}

export function useMealPlansByDateRange(startDate: string, endDate: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['mealPlans', startDate, endDate],
    queryFn: () => api.getMealPlansByDateRange(startDate, endDate),
    enabled,
  })
}

export function useUpdateMealPlan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ date, data }: { date: string; data: UpdateMealPlanRequest }) =>
      api.updateMealPlan(date, data),
    onSuccess: (_, variables) => {
      // Invalidate both the specific date and date range queries
      queryClient.invalidateQueries({ queryKey: ['mealPlan', variables.date] })
      queryClient.invalidateQueries({ queryKey: ['mealPlans'] })
    },
  })
}
