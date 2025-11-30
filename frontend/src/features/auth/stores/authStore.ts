import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '../types/auth.types'
import { authService } from '@/services/authService'
import { logger } from '@/utils/logger'
import { isDesignMode, logDesignModeStatus, getMockUser } from '@/utils/designMode'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => {
    // In design mode, always consider user as authenticated
    if (isDesignMode()) return true
    return !!user.value
  })
  const userRole = computed(() => {
    // In design mode, return admin role for full access
    if (isDesignMode()) return 'admin'
    return user.value?.role || null
  })
  const isAdmin = computed(() => {
    // In design mode, always admin
    if (isDesignMode()) return true
    return user.value?.role === 'admin'
  })
  const isTeacher = computed(() => {
    // In design mode, has teacher rights too
    if (isDesignMode()) return true
    return user.value?.role === 'teacher'
  })

  // Actions
  async function login(credentials: LoginRequest) {
    logger.logStoreAction('authStore', 'login', { email: credentials.email })
    loading.value = true
    error.value = null

    try {
      // Design mode: skip login, just resolve without setting user
      if (isDesignMode()) {
        logger.info('Design mode: Login bypassed', 'authStore')
        // Return empty response - user remains null
        return { user: null as any, token: '' }
      }

      const response = await authService.login(credentials)
      user.value = response.user

      // Save only user data to localStorage (token is in HTTP-only cookie)
      localStorage.setItem('user', JSON.stringify(response.user))

      logger.info('User logged in successfully', 'authStore', {
        userId: response.user.id,
        role: response.user.role,
      })
      return response
    } catch (err: any) {
      error.value = err.message || 'Login failed'
      logger.logStoreError('authStore', 'login', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    loading.value = true
    error.value = null

    try {
      // Design mode: skip registration, just resolve without setting user
      if (isDesignMode()) {
        logger.info('Design mode: Registration bypassed', 'authStore')
        // Return empty response - user remains null
        return { user: null as any, token: '' }
      }

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
    logger.logStoreAction('authStore', 'logout')
    loading.value = true
    error.value = null

    try {
      await authService.logout()
      logger.info('User logged out successfully', 'authStore')
    } catch (err: any) {
      logger.logStoreError('authStore', 'logout', err)
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
      // Design mode: skip fetching user
      if (isDesignMode()) {
        logger.info('Design mode: Fetch user bypassed', 'authStore')
        // Return null - user remains null
        return null as any
      }

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
    // Log design mode status
    logDesignModeStatus()

    // Design mode: set mock user for components that need it
    if (isDesignMode()) {
      user.value = getMockUser() as any
      logger.info('Design mode: Mock user set, free access to all pages', 'authStore', user.value)
      return
    }

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
