import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { AxiosError } from 'axios'

import { authService } from '@/services/authService'
import type { User, LoginInput, RegisterInput } from '@/types'

export const useAuthStore = defineStore(
  'auth',
  () => {
    const token = ref<string | null>(null)
    const user = ref<User | null>(null)
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    const isAuthenticated = computed(() => !!token.value && !!user.value)

    async function login(credentials: LoginInput): Promise<void> {
      isLoading.value = true
      error.value = null
      try {
        const response = await authService.login(credentials)
        token.value = response.token
        user.value = response.user
      } catch (err: unknown) {
        const axiosErr = err as AxiosError<{ message?: string }>
        error.value = axiosErr.response?.data?.message ?? 'Login failed. Please try again.'
        throw err
      } finally {
        isLoading.value = false
      }
    }

    async function register(payload: RegisterInput): Promise<void> {
      isLoading.value = true
      error.value = null
      try {
        const body = { name: payload.name, email: payload.email, password: payload.password }
        const response = await authService.register(body)
        token.value = response.token
        user.value = response.user
      } catch (err: unknown) {
        const axiosErr = err as AxiosError<{ message?: string }>
        error.value = axiosErr.response?.data?.message ?? 'Registration failed. Please try again.'
        throw err
      } finally {
        isLoading.value = false
      }
    }

    function logout(): void {
      authService.logout()
      token.value = null
      user.value = null
      error.value = null
    }

    function clearError(): void {
      error.value = null
    }

    return {
      token,
      user,
      isLoading,
      error,
      isAuthenticated,
      login,
      register,
      logout,
      clearError
    }
  },
  {
    persist: {
      pick: ['token', 'user']
    }
  }
)
