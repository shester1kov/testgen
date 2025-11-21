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

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  const mockUser: User = {
    id: '123',
    email: 'test@example.com',
    full_name: 'Test User',
    role: 'student',
  }

  const mockAuthResponse: AuthResponse = {
    user: mockUser,
    token: 'mock-token',
  }

  describe('login', () => {
    it('should login successfully and store user data', async () => {
      vi.mocked(authService.login).mockResolvedValue(mockAuthResponse)

      const store = useAuthStore()
      await store.login({ email: 'test@example.com', password: 'password' })

      expect(store.user).toEqual(mockUser)
      expect(store.isAuthenticated).toBe(true)
      expect(localStorage.setItem).toHaveBeenCalledWith('user', JSON.stringify(mockUser))
    })

    it('should handle login error', async () => {
      const error = new Error('Invalid credentials')
      vi.mocked(authService.login).mockRejectedValue(error)

      const store = useAuthStore()
      await expect(store.login({ email: 'test@example.com', password: 'wrong' })).rejects.toThrow()

      expect(store.user).toBeNull()
      expect(store.error).toBe('Invalid credentials')
    })
  })

  describe('register', () => {
    it('should register successfully without role selection', async () => {
      vi.mocked(authService.register).mockResolvedValue(mockAuthResponse)

      const store = useAuthStore()
      await store.register({
        email: 'new@example.com',
        password: 'password',
        full_name: 'New User',
      })

      expect(store.user).toEqual(mockUser)
      expect(store.isAuthenticated).toBe(true)
      expect(authService.register).toHaveBeenCalledWith({
        email: 'new@example.com',
        password: 'password',
        full_name: 'New User',
      })
    })

    it('should handle registration error', async () => {
      const error = new Error('Email already exists')
      vi.mocked(authService.register).mockRejectedValue(error)

      const store = useAuthStore()
      await expect(store.register({
        email: 'test@example.com',
        password: 'password',
        full_name: 'Test User',
      })).rejects.toThrow()

      expect(store.user).toBeNull()
      expect(store.error).toBe('Email already exists')
    })
  })

  describe('logout', () => {
    it('should logout and clear user data', async () => {
      vi.mocked(authService.logout).mockResolvedValue()

      const store = useAuthStore()
      store.user = mockUser

      await store.logout()

      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
      expect(localStorage.removeItem).toHaveBeenCalledWith('user')
    })

    it('should clear user data even if API call fails', async () => {
      vi.mocked(authService.logout).mockRejectedValue(new Error('Network error'))

      const store = useAuthStore()
      store.user = mockUser

      await store.logout()

      expect(store.user).toBeNull()
      expect(localStorage.removeItem).toHaveBeenCalledWith('user')
    })
  })

  describe('fetchUser', () => {
    it('should fetch current user data', async () => {
      vi.mocked(authService.getMe).mockResolvedValue(mockUser)

      const store = useAuthStore()
      await store.fetchUser()

      expect(store.user).toEqual(mockUser)
      expect(localStorage.setItem).toHaveBeenCalledWith('user', JSON.stringify(mockUser))
    })

    it('should clear user on fetch error', async () => {
      vi.mocked(authService.getMe).mockRejectedValue(new Error('Unauthorized'))

      const store = useAuthStore()
      store.user = mockUser

      await expect(store.fetchUser()).rejects.toThrow()

      expect(store.user).toBeNull()
      expect(localStorage.removeItem).toHaveBeenCalledWith('user')
    })
  })

  describe('computed properties', () => {
    it('should correctly compute isAdmin', () => {
      const store = useAuthStore()
      store.user = { ...mockUser, role: 'admin' }

      expect(store.isAdmin).toBe(true)
      expect(store.isTeacher).toBe(false)
    })

    it('should correctly compute isTeacher', () => {
      const store = useAuthStore()
      store.user = { ...mockUser, role: 'teacher' }

      expect(store.isTeacher).toBe(true)
      expect(store.isAdmin).toBe(false)
    })

    it('should correctly compute userRole', () => {
      const store = useAuthStore()
      store.user = mockUser

      expect(store.userRole).toBe('student')
    })
  })

  describe('initializeAuth', () => {
    it('should restore user from localStorage', () => {
      vi.mocked(localStorage.getItem).mockReturnValue(JSON.stringify(mockUser))
      vi.mocked(authService.getMe).mockResolvedValue(mockUser)

      const store = useAuthStore()
      store.initializeAuth()

      expect(store.user).toEqual(mockUser)
    })

    it('should not restore user if localStorage is empty', () => {
      vi.mocked(localStorage.getItem).mockReturnValue(null)

      const store = useAuthStore()
      store.initializeAuth()

      expect(store.user).toBeNull()
    })
  })
})
