import { describe, it, expect, beforeEach, vi } from 'vitest'
import userService from '../userService'
import api from '../api'

vi.mock('../api')

describe('userService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('listUsers', () => {
    // POSITIVE TEST: Successfully fetch users list
    it('fetches users list with pagination', async () => {
      const mockResponse = {
        users: [
          {
            id: '1',
            email: 'user1@test.com',
            full_name: 'User One',
            role: 'student',
            created_at: '2024-01-01T00:00:00Z',
          },
          {
            id: '2',
            email: 'user2@test.com',
            full_name: 'User Two',
            role: 'teacher',
            created_at: '2024-01-02T00:00:00Z',
          },
        ],
        total: 2,
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)

      const result = await userService.listUsers(20, 0)

      expect(api.get).toHaveBeenCalledWith('/users', {
        params: { limit: 20, offset: 0 },
      })
      expect(result).toEqual(mockResponse)
      expect(result.users).toHaveLength(2)
      expect(result.total).toBe(2)
    })

    // POSITIVE TEST: Fetch with custom pagination
    it('fetches users with custom limit and offset', async () => {
      const mockResponse = {
        users: [],
        total: 50,
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)

      await userService.listUsers(10, 20)

      expect(api.get).toHaveBeenCalledWith('/users', {
        params: { limit: 10, offset: 20 },
      })
    })

    // NEGATIVE TEST: API returns error
    it('throws error when API request fails', async () => {
      const mockError = new Error('Network error')
      vi.mocked(api.get).mockRejectedValue(mockError)

      await expect(userService.listUsers()).rejects.toThrow('Network error')
    })

    // NEGATIVE TEST: Empty users list
    it('handles empty users list', async () => {
      const mockResponse = {
        users: [],
        total: 0,
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)

      const result = await userService.listUsers()

      expect(result.users).toEqual([])
      expect(result.total).toBe(0)
    })
  })

  describe('updateUserRole', () => {
    // POSITIVE TEST: Successfully update user role to admin
    it('updates user role to admin', async () => {
      const mockUser = {
        id: '123',
        email: 'user@test.com',
        full_name: 'Test User',
        role: 'admin',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.put).mockResolvedValue(mockUser)

      const result = await userService.updateUserRole('123', 'admin')

      expect(api.put).toHaveBeenCalledWith('/users/123/role', {
        role_name: 'admin',
      })
      expect(result).toEqual(mockUser)
      expect(result.role).toBe('admin')
    })

    // POSITIVE TEST: Update role to teacher
    it('updates user role to teacher', async () => {
      const mockUser = {
        id: '456',
        email: 'teacher@test.com',
        full_name: 'Teacher User',
        role: 'teacher',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.put).mockResolvedValue(mockUser)

      const result = await userService.updateUserRole('456', 'teacher')

      expect(api.put).toHaveBeenCalledWith('/users/456/role', {
        role_name: 'teacher',
      })
      expect(result.role).toBe('teacher')
    })

    // POSITIVE TEST: Update role to student
    it('updates user role to student', async () => {
      const mockUser = {
        id: '789',
        email: 'student@test.com',
        full_name: 'Student User',
        role: 'student',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.put).mockResolvedValue(mockUser)

      const result = await userService.updateUserRole('789', 'student')

      expect(result.role).toBe('student')
    })

    // NEGATIVE TEST: User not found
    it('throws error when user not found', async () => {
      const mockError = {
        response: {
          status: 404,
          data: { message: 'User not found' },
        },
      }
      vi.mocked(api.put).mockRejectedValue(mockError)

      await expect(userService.updateUserRole('999', 'admin')).rejects.toEqual(mockError)
    })

    // NEGATIVE TEST: Unauthorized access (403)
    it('throws error when user lacks permissions', async () => {
      const mockError = {
        response: {
          status: 403,
          data: { message: 'Access denied' },
        },
      }
      vi.mocked(api.put).mockRejectedValue(mockError)

      await expect(userService.updateUserRole('123', 'admin')).rejects.toEqual(mockError)
    })

    // NEGATIVE TEST: Network error
    it('throws error on network failure', async () => {
      const mockError = new Error('Network error')
      vi.mocked(api.put).mockRejectedValue(mockError)

      await expect(userService.updateUserRole('123', 'admin')).rejects.toThrow('Network error')
    })

    // NEGATIVE TEST: Invalid user ID
    it('handles invalid user ID', async () => {
      const mockError = {
        response: {
          status: 400,
          data: { message: 'Invalid user ID' },
        },
      }
      vi.mocked(api.put).mockRejectedValue(mockError)

      await expect(userService.updateUserRole('invalid-id', 'admin')).rejects.toEqual(mockError)
    })

    // NEGATIVE TEST: Server error (500)
    it('handles server errors', async () => {
      const mockError = {
        response: {
          status: 500,
          data: { message: 'Internal server error' },
        },
      }
      vi.mocked(api.put).mockRejectedValue(mockError)

      await expect(userService.updateUserRole('123', 'admin')).rejects.toEqual(mockError)
    })
  })
})
