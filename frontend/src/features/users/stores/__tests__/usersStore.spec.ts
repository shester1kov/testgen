import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUsersStore } from '../usersStore'
import userService from '@/services/userService'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/userService')
vi.mock('@/utils/logger', () => ({
  default: {
    logStoreAction: vi.fn(),
    logStoreError: vi.fn(),
    info: vi.fn(),
  },
}))

describe('usersStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchUsers', () => {
    // POSITIVE TEST: Successfully fetch users
    it('fetches users successfully', async () => {
      const mockUsers: User[] = [
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
      ]

      vi.mocked(userService.listUsers).mockResolvedValue({
        users: mockUsers,
        total: 2,
      })

      const store = useUsersStore()
      await store.fetchUsers(1)

      expect(store.users).toEqual(mockUsers)
      expect(store.total).toBe(2)
      expect(store.currentPage).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    // POSITIVE TEST: Fetch with pagination
    it('fetches users with correct pagination', async () => {
      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [],
        total: 100,
      })

      const store = useUsersStore()
      await store.fetchUsers(3)

      const expectedOffset = (3 - 1) * 20 // page 3, default pageSize 20
      expect(userService.listUsers).toHaveBeenCalledWith(20, expectedOffset)
      expect(store.currentPage).toBe(3)
    })

    // POSITIVE TEST: Calculate total pages correctly
    it('calculates total pages correctly', async () => {
      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [],
        total: 45,
      })

      const store = useUsersStore()
      await store.fetchUsers(1)

      expect(store.totalPages).toBe(3) // 45 users / 20 per page = 3 pages
    })

    // NEGATIVE TEST: API error handling
    it('handles API errors gracefully', async () => {
      const mockError = {
        response: {
          data: { message: 'Failed to fetch users' },
        },
      }
      vi.mocked(userService.listUsers).mockRejectedValue(mockError)

      const store = useUsersStore()

      await expect(store.fetchUsers(1)).rejects.toEqual(mockError)
      expect(store.error).toBe('Failed to fetch users')
      expect(store.loading).toBe(false)
    })

    // NEGATIVE TEST: Network error
    it('handles network errors', async () => {
      const mockError = new Error('Network error')
      vi.mocked(userService.listUsers).mockRejectedValue(mockError)

      const store = useUsersStore()

      await expect(store.fetchUsers(1)).rejects.toThrow('Network error')
      expect(store.error).toBe('Failed to fetch users')
      expect(store.loading).toBe(false)
    })

    // POSITIVE TEST: Loading state management
    it('sets loading state correctly', async () => {
      vi.mocked(userService.listUsers).mockImplementation(
        () =>
          new Promise(resolve => {
            setTimeout(
              () =>
                resolve({
                  users: [],
                  total: 0,
                }),
              100
            )
          })
      )

      const store = useUsersStore()
      const promise = store.fetchUsers(1)

      expect(store.loading).toBe(true)

      await promise

      expect(store.loading).toBe(false)
    })

    // POSITIVE TEST: Empty users list
    it('handles empty users list', async () => {
      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [],
        total: 0,
      })

      const store = useUsersStore()
      await store.fetchUsers(1)

      expect(store.users).toEqual([])
      expect(store.total).toBe(0)
      expect(store.totalPages).toBe(0)
    })
  })

  describe('updateUserRole', () => {
    // POSITIVE TEST: Successfully update user role
    it('updates user role successfully', async () => {
      const updatedUser: User = {
        id: '1',
        email: 'user@test.com',
        full_name: 'Test User',
        role: 'admin',
        created_at: '2024-01-01T00:00:00Z',
      }

      const initialUser: User = {
        ...updatedUser,
        role: 'student',
      }

      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [initialUser],
        total: 1,
      })

      vi.mocked(userService.updateUserRole).mockResolvedValue(updatedUser)

      const store = useUsersStore()
      await store.fetchUsers(1)

      expect(store.users[0].role).toBe('student')

      const result = await store.updateUserRole('1', 'admin')

      expect(result).toEqual(updatedUser)
      expect(store.users[0].role).toBe('admin')
      expect(store.error).toBeNull()
      expect(store.loading).toBe(false)
    })

    // POSITIVE TEST: Update role to teacher
    it('updates user role to teacher', async () => {
      const updatedUser: User = {
        id: '2',
        email: 'teacher@test.com',
        full_name: 'Teacher User',
        role: 'teacher',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [{ ...updatedUser, role: 'student' }],
        total: 1,
      })

      vi.mocked(userService.updateUserRole).mockResolvedValue(updatedUser)

      const store = useUsersStore()
      await store.fetchUsers(1)
      await store.updateUserRole('2', 'teacher')

      expect(store.users[0].role).toBe('teacher')
    })

    // NEGATIVE TEST: Update fails with error
    it('handles update errors gracefully', async () => {
      const mockError = {
        response: {
          data: { message: 'Access denied' },
        },
      }
      vi.mocked(userService.updateUserRole).mockRejectedValue(mockError)

      const store = useUsersStore()

      await expect(store.updateUserRole('1', 'admin')).rejects.toEqual(mockError)
      expect(store.error).toBe('Access denied')
      expect(store.loading).toBe(false)
    })

    // NEGATIVE TEST: User not found in local state
    it('handles user not found in local state', async () => {
      const updatedUser: User = {
        id: '999',
        email: 'user@test.com',
        full_name: 'Test User',
        role: 'admin',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(userService.listUsers).mockResolvedValue({
        users: [],
        total: 0,
      })

      vi.mocked(userService.updateUserRole).mockResolvedValue(updatedUser)

      const store = useUsersStore()
      await store.fetchUsers(1)

      // Should not throw, but user won't be in local state
      await store.updateUserRole('999', 'admin')

      expect(store.users.find(u => u.id === '999')).toBeUndefined()
    })

    // NEGATIVE TEST: Network error during update
    it('handles network errors during update', async () => {
      const mockError = new Error('Network error')
      vi.mocked(userService.updateUserRole).mockRejectedValue(mockError)

      const store = useUsersStore()

      await expect(store.updateUserRole('1', 'admin')).rejects.toThrow('Network error')
      expect(store.error).toBe('Failed to update user role')
      expect(store.loading).toBe(false)
    })

    // POSITIVE TEST: Loading state during update
    it('sets loading state correctly during update', async () => {
      vi.mocked(userService.updateUserRole).mockImplementation(
        () =>
          new Promise(resolve => {
            setTimeout(
              () =>
                resolve({
                  id: '1',
                  email: 'user@test.com',
                  full_name: 'Test User',
                  role: 'admin',
                  created_at: '2024-01-01T00:00:00Z',
                }),
              100
            )
          })
      )

      const store = useUsersStore()
      const promise = store.updateUserRole('1', 'admin')

      expect(store.loading).toBe(true)

      await promise

      expect(store.loading).toBe(false)
    })
  })

  describe('clearError', () => {
    // POSITIVE TEST: Clear error
    it('clears error message', async () => {
      const mockError = {
        response: {
          data: { message: 'Some error' },
        },
      }
      vi.mocked(userService.listUsers).mockRejectedValue(mockError)

      const store = useUsersStore()

      await expect(store.fetchUsers(1)).rejects.toEqual(mockError)
      expect(store.error).toBe('Some error')

      store.clearError()

      expect(store.error).toBeNull()
    })
  })

  describe('totalPages computed', () => {
    // POSITIVE TEST: Calculate total pages for various totals
    it('calculates total pages correctly for different totals', async () => {
      const testCases = [
        { total: 0, expectedPages: 0 },
        { total: 1, expectedPages: 1 },
        { total: 20, expectedPages: 1 },
        { total: 21, expectedPages: 2 },
        { total: 40, expectedPages: 2 },
        { total: 41, expectedPages: 3 },
        { total: 100, expectedPages: 5 },
      ]

      for (const testCase of testCases) {
        vi.mocked(userService.listUsers).mockResolvedValue({
          users: [],
          total: testCase.total,
        })

        const store = useUsersStore()
        await store.fetchUsers(1)

        expect(store.totalPages).toBe(testCase.expectedPages)
      }
    })
  })
})
