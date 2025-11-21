import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '../types/auth.types'
import { authService } from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isTeacher = computed(() => user.value?.role === 'teacher')

  // Actions
  async function login(credentials: LoginRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.login(credentials)
      user.value = response.user

      // Save only user data to localStorage (token is in HTTP-only cookie)
      localStorage.setItem('user', JSON.stringify(response.user))

      return response
    } catch (err: any) {
      error.value = err.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authService.register(data)
      user.value = response.user

      // Save only user data to localStorage (token is in HTTP-only cookie)
      localStorage.setItem('user', JSON.stringify(response.user))

      return response
    } catch (err: any) {
      error.value = err.message || 'Registration failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    error.value = null

    try {
      await authService.logout()
    } catch (err: any) {
      console.error('Logout error:', err)
    } finally {
      // Clear state regardless of API call result
      user.value = null
      localStorage.removeItem('user')
      loading.value = false
    }
  }

  async function fetchUser() {
    loading.value = true
    error.value = null

    try {
      const userData = await authService.getMe()
      user.value = userData
      localStorage.setItem('user', JSON.stringify(userData))
      return userData
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch user'
      // Clear user if token is invalid
      user.value = null
      localStorage.removeItem('user')
      throw err
    } finally {
      loading.value = false
    }
  }

  function initializeAuth() {
    const storedUser = localStorage.getItem('user')

    if (storedUser) {
      user.value = JSON.parse(storedUser)
      // Verify token is still valid by fetching user data
      fetchUser().catch(() => {
        // Token is invalid, clear user
        user.value = null
        localStorage.removeItem('user')
      })
    }
  }

  return {
    // State
    user,
    loading,
    error,
    // Getters
    isAuthenticated,
    userRole,
    isAdmin,
    isTeacher,
    // Actions
    login,
    register,
    logout,
    fetchUser,
    initializeAuth,
  }
})
