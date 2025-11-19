import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './authStore'
import { authService } from '@/services/authService'
import type { AuthResponse, User } from '../types/auth.types'

// Mock authService
vi.mock('@/services/authService', () => ({
  authService: {
    login: vi.fn(),
    register: vi.fn(),
    logout: vi.fn(),
    getMe: vi.fn(),
  },
}))

describe('AuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('should login user successfully', async () => {
      const mockResponse: AuthResponse = {
        user: {
          id: '1',
          email: 'test@example.com',
          full_name: 'Test User',
          role: 'teacher',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        token: 'test-token',
      }

      vi.mocked(authService.login).mockResolvedValue(mockResponse)

      const store = useAuthStore()
      await store.login({ email: 'test@example.com', password: 'password' })

      expect(store.user).toEqual(mockResponse.user)
      expect(store.token).toBe('test-token')
      expect(store.isAuthenticated).toBe(true)
      expect(localStorage.getItem('auth_token')).toBe('test-token')
    })

    it('should handle login error', async () => {
      vi.mocked(authService.login).mockRejectedValue(new Error('Invalid credentials'))

      const store = useAuthStore()
      await expect(store.login({ email: 'test@example.com', password: 'wrong' })).rejects.toThrow()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('logout', () => {
    it('should logout user successfully', async () => {
      const store = useAuthStore()
      store.user = {
        id: '1',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'teacher',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      store.token = 'test-token'

      await store.logout()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
      expect(localStorage.getItem('auth_token')).toBeNull()
    })
  })

  describe('initializeAuth', () => {
    it('should initialize auth from localStorage', () => {
      const mockUser: User = {
        id: '1',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'teacher',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }

      localStorage.setItem('auth_token', 'test-token')
      localStorage.setItem('user', JSON.stringify(mockUser))

      const store = useAuthStore()
      store.initializeAuth()

      expect(store.token).toBe('test-token')
      expect(store.user).toEqual(mockUser)
      expect(store.isAuthenticated).toBe(true)
    })
  })

  describe('computed properties', () => {
    it('should check user role correctly', () => {
      const store = useAuthStore()
      store.user = {
        id: '1',
        email: 'admin@example.com',
        full_name: 'Admin User',
        role: 'admin',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }

      expect(store.userRole).toBe('admin')
      expect(store.isAdmin).toBe(true)
      expect(store.isTeacher).toBe(false)
    })
  })
})
