import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import userService from '@/services/userService'
import type { User } from '@/features/auth/types/auth.types'
import logger from '@/utils/logger'
import { isDesignMode, getMockUsers } from '@/utils/designMode'

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

  /**
   * Fetch users list
   */
  async function fetchUsers(page: number = 1) {
    loading.value = true
    error.value = null
    logger.logStoreAction('usersStore', 'fetchUsers', { page })

    try {
      const offset = (page - 1) * pageSize.value
      const response = await userService.listUsers(pageSize.value, offset)

      users.value = response.users
      total.value = response.total
      currentPage.value = page

      logger.info('Users fetched successfully', 'usersStore', { count: users.value.length })
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to fetch users'
      logger.logStoreError('usersStore', 'fetchUsers', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update user role
   */
  async function updateUserRole(userId: string, roleName: 'admin' | 'teacher' | 'student') {
    loading.value = true
    error.value = null
    logger.logStoreAction('usersStore', 'updateUserRole', { userId, roleName })

    try {
      const updatedUser = await userService.updateUserRole(userId, roleName)

      // Update user in local state
      const index = users.value.findIndex(u => u.id === userId)
      if (index !== -1) {
        users.value[index] = updatedUser
      }

      logger.info('User role updated successfully', 'usersStore', { userId, roleName })
      return updatedUser
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to update user role'
      logger.logStoreError('usersStore', 'updateUserRole', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Clear error message
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    users,
    total,
    currentPage,
    pageSize,
    totalPages,
    loading,
    error,

    // Actions
    fetchUsers,
    updateUserRole,
    clearError
  }
})
