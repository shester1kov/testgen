import { describe, it, expect, beforeEach, vi } from 'vitest'
import { authService } from './authService'
import api from './api'

vi.mock('./api')

describe('Auth Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('should call API with credentials and return auth response', async () => {
      const mockResponse = {
        user: {
          id: '123',
          email: 'test@example.com',
          full_name: 'Test User',
          role: 'student',
        },
        token: 'mock-token',
      }

      // api.post returns response.data directly due to interceptor
      vi.mocked(api.post).mockResolvedValue(mockResponse as any)

      const result = await authService.login({
        email: 'test@example.com',
        password: 'password',
      })

      expect(api.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@example.com',
        password: 'password',
      })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('register', () => {
    it('should call API without role field and return auth response', async () => {
      const mockResponse = {
        user: {
          id: '123',
          email: 'new@example.com',
          full_name: 'New User',
          role: 'student', // Should be automatically assigned by backend
        },
        token: 'mock-token',
      }

      // api.post returns response.data directly due to interceptor
      vi.mocked(api.post).mockResolvedValue(mockResponse as any)

      const result = await authService.register({
        email: 'new@example.com',
        password: 'password',
        full_name: 'New User',
      })

      expect(api.post).toHaveBeenCalledWith('/auth/register', {
        email: 'new@example.com',
        password: 'password',
        full_name: 'New User',
      })
      expect(result).toEqual(mockResponse)
      expect(result.user.role).toBe('student')
    })
  })

  describe('logout', () => {
    it('should call logout API endpoint', async () => {
      vi.mocked(api.post).mockResolvedValue(undefined as any)

      await authService.logout()

      expect(api.post).toHaveBeenCalledWith('/auth/logout')
    })
  })

  describe('getMe', () => {
    it('should fetch current user data', async () => {
      const mockUser = {
        id: '123',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'student',
      }

      // api.get returns response.data directly due to interceptor
      vi.mocked(api.get).mockResolvedValue(mockUser as any)

      const result = await authService.getMe()

      expect(api.get).toHaveBeenCalledWith('/auth/me')
      expect(result).toEqual(mockUser)
    })
  })
})
