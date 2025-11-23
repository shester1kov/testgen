import api from './api'
import type { User } from '@/features/auth/types/auth.types'

export interface UserListResponse {
  users: User[]
  total: number
  limit: number
  offset: number
}

export interface UpdateUserRoleRequest {
  role_name: 'admin' | 'teacher' | 'student'
}

const userService = {
  /**
   * Get list of all users (admin/teacher only)
   */
  async listUsers(limit: number = 20, offset: number = 0): Promise<UserListResponse> {
    const response = await api.get<UserListResponse>('/users', {
      params: { limit, offset }
    })
    return response
  },

  /**
   * Update user role (admin only)
   */
  async updateUserRole(userId: string, roleName: 'admin' | 'teacher' | 'student'): Promise<User> {
    const response = await api.put<User>(`/users/${userId}/role`, {
      role_name: roleName
    })
    return response
  }
}

export default userService
