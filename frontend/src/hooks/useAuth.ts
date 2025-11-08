/**
 * useAuth Hook
 * Provides authentication actions and state
 */

import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { api, RegisterRequest, LoginRequest } from '../lib/api'
import { useAuthStore } from '../store/authStore'

export function useAuth() {
  const { user, isAuthenticated, setUser, logout: logoutStore } = useAuthStore()
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  const registerMutation = useMutation({
    mutationFn: (data: RegisterRequest) => api.register(data),
    onSuccess: (response) => {
      setUser(response.data.user)
      navigate({ to: '/' })
    },
  })

  const loginMutation = useMutation({
    mutationFn: (data: LoginRequest) => api.login(data),
    onSuccess: (response) => {
      setUser(response.data.user)
      navigate({ to: '/' })
    },
  })

  const logoutMutation = useMutation({
    mutationFn: () => api.logout(),
    onSuccess: () => {
      logoutStore()
      queryClient.clear()
      navigate({ to: '/login' })
    },
  })

  return {
    user,
    isAuthenticated,
    register: registerMutation.mutate,
    login: loginMutation.mutate,
    logout: logoutMutation.mutate,
    isLoading:
      registerMutation.isPending ||
      loginMutation.isPending ||
      logoutMutation.isPending,
    error:
      registerMutation.error ||
      loginMutation.error ||
      logoutMutation.error,
  }
}
