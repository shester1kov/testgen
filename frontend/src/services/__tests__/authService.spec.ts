import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { authService } from '../authService'
import api from '../api'
import type { ApiError } from '@/types/api.types'

// Mock the api module
vi.mock('../api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
  },
}))

describe('AuthService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('POSITIVE: Successful operations', () => {
    it('should login successfully', async () => {
      const mockResponse = {
        user: {
          id: '123',
          email: 'test@example.com',
          full_name: 'Test User',
          role: 'student',
        },
        token: 'mock-jwt-token',
      }

      vi.mocked(api.post).mockResolvedValueOnce(mockResponse)

      const result = await authService.login({
        email: 'test@example.com',
        password: 'password123',
      })

      expect(api.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      })
      expect(result).toEqual(mockResponse)
      expect(result.token).toBe('mock-jwt-token')
    })

    it('should register successfully', async () => {
      const mockResponse = {
        user: {
          id: '456',
          email: 'new@example.com',
          full_name: 'Test User',
          role: 'student' // Default role assigned by backend
        },
        token: 'mock-jwt-token',
      }

      vi.mocked(api.post).mockResolvedValueOnce(mockResponse)

      const result = await authService.register({
        email: 'new@example.com',
        password: 'securepass',
        full_name: 'Test User',
      })

      expect(api.post).toHaveBeenCalledWith('/auth/register', {
        email: 'new@example.com',
        password: 'securepass',
        full_name: 'Test User',
      })
      expect(result).toEqual(mockResponse)
      expect(result.user.role).toBe('student')
      expect(result.token).toBe('mock-jwt-token')
    })

    it('should logout successfully', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ message: 'Logged out' })

      await authService.logout()

      expect(api.post).toHaveBeenCalledWith('/auth/logout')
    })

    it('should fetch current user successfully', async () => {
      const mockUser = { id: '789', email: 'me@example.com', role: 'admin' }

      vi.mocked(api.get).mockResolvedValueOnce(mockUser)

      const result = await authService.getMe()

      expect(api.get).toHaveBeenCalledWith('/auth/me')
      expect(result).toEqual(mockUser)
    })
  })

  describe('NEGATIVE: Invalid credentials', () => {
    it('should handle login with empty email', async () => {
      const error: ApiError = {
        message: 'Invalid email',
        code: 'INVALID_INPUT',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: '',
          password: 'password',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle login with empty password', async () => {
      const error: ApiError = {
        message: 'Invalid password',
        code: 'INVALID_INPUT',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: 'test@example.com',
          password: '',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle login with malformed email', async () => {
      const error: ApiError = {
        message: 'Invalid email format',
        code: 'INVALID_EMAIL',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: 'not-an-email',
          password: 'password',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle login with wrong credentials', async () => {
      const error: ApiError = {
        message: 'Invalid credentials',
        code: 'UNAUTHORIZED',
        status: 401,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: 'test@example.com',
          password: 'wrongpass',
        }),
      ).rejects.toEqual(error)
    })
  })

  describe('NEGATIVE: Registration validation errors', () => {
    it('should handle duplicate email registration', async () => {
      const error: ApiError = {
        message: 'Email already exists',
        code: 'DUPLICATE_EMAIL',
        status: 409,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.register({
          email: 'existing@example.com',
          password: 'password123',
          full_name: 'Test User',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle weak password', async () => {
      const error: ApiError = {
        message: 'Password too weak',
        code: 'WEAK_PASSWORD',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.register({
          email: 'test@example.com',
          password: '123',
          full_name: 'Test',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle empty full name', async () => {
      const error: ApiError = {
        message: 'Full name is required',
        code: 'INVALID_INPUT',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.register({
          email: 'test@example.com',
          password: 'password123',
          full_name: '',
        }),
      ).rejects.toEqual(error)
    })

  })

  describe('NEGATIVE: Network errors', () => {
    it('should handle network timeout on login', async () => {
      const error: ApiError = {
        message: 'Network error. Please check your connection.',
        status: 0,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: 'test@example.com',
          password: 'password',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle server error (500)', async () => {
      const error: ApiError = {
        message: 'Internal server error',
        code: 'INTERNAL_ERROR',
        status: 500,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(authService.register({} as any)).rejects.toEqual(error)
    })

    it('should handle service unavailable (503)', async () => {
      const error: ApiError = {
        message: 'Service temporarily unavailable',
        code: 'SERVICE_UNAVAILABLE',
        status: 503,
      }

      vi.mocked(api.get).mockRejectedValueOnce(error)

      await expect(authService.getMe()).rejects.toEqual(error)
    })
  })

  describe('NEGATIVE: Unauthorized access', () => {
    it('should handle expired token when fetching user', async () => {
      const error: ApiError = {
        message: 'Token expired',
        code: 'TOKEN_EXPIRED',
        status: 401,
      }

      vi.mocked(api.get).mockRejectedValueOnce(error)

      await expect(authService.getMe()).rejects.toEqual(error)
    })

    it('should handle invalid token', async () => {
      const error: ApiError = {
        message: 'Invalid token',
        code: 'INVALID_TOKEN',
        status: 401,
      }

      vi.mocked(api.get).mockRejectedValueOnce(error)

      await expect(authService.getMe()).rejects.toEqual(error)
    })

    it('should handle missing token', async () => {
      const error: ApiError = {
        message: 'No authentication token provided',
        code: 'NO_TOKEN',
        status: 401,
      }

      vi.mocked(api.get).mockRejectedValueOnce(error)

      await expect(authService.getMe()).rejects.toEqual(error)
    })
  })

  describe('NEGATIVE: Special character handling', () => {
    it('should handle email with special characters', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      await authService.login({
        email: "test+tag@example.com",
        password: 'password',
      })

      expect(api.post).toHaveBeenCalledWith('/auth/login', {
        email: "test+tag@example.com",
        password: 'password',
      })
    })

    it('should handle password with special characters', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      await authService.login({
        email: 'test@example.com',
        password: 'p@ssw0rd!#$%^&*()',
      })

      expect(api.post).toHaveBeenCalled()
    })

    it('should handle full name with unicode', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      await authService.register({
        email: 'test@example.com',
        password: 'password',
        full_name: '张三 Müller',
      })

      expect(api.post).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: SQL injection attempts', () => {
    it('should handle SQL injection in email', async () => {
      const error: ApiError = {
        message: 'Invalid email format',
        code: 'INVALID_INPUT',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: "' OR '1'='1",
          password: 'password',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle SQL injection in password', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      // Backend should handle this safely
      await authService.login({
        email: 'test@example.com',
        password: "'; DROP TABLE users--",
      })

      expect(api.post).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: XSS attempts', () => {
    it('should handle XSS in full name', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      await authService.register({
        email: 'test@example.com',
        password: 'password',
        full_name: "<script>alert('XSS')</script>",
      })

      // Backend should sanitize this
      expect(api.post).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Very long inputs', () => {
    it('should handle very long email', async () => {
      const longEmail = 'a'.repeat(1000) + '@example.com'

      const error: ApiError = {
        message: 'Email too long',
        code: 'INVALID_INPUT',
        status: 400,
      }

      vi.mocked(api.post).mockRejectedValueOnce(error)

      await expect(
        authService.login({
          email: longEmail,
          password: 'password',
        }),
      ).rejects.toEqual(error)
    })

    it('should handle very long password', async () => {
      const longPassword = 'a'.repeat(10000)

      vi.mocked(api.post).mockResolvedValueOnce({ user: {} })

      await authService.register({
        email: 'test@example.com',
        password: longPassword,
        full_name: 'Test',
      })

      expect(api.post).toHaveBeenCalled()
    })
  })
})
