import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'

export function useNutritionSearch(query: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['nutrition', query],
    queryFn: () => api.searchNutrition(query),
    enabled: enabled && query.length >= 2, // Only search if query is at least 2 chars
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
  })
}
