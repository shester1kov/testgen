import api from './api'
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '@/features/auth/types/auth.types'

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    // api.post returns data directly due to response interceptor
    const response = await api.post('/auth/login', credentials)
    return response as AuthResponse
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    // api.post returns data directly due to response interceptor
    const response = await api.post('/auth/register', data)
    return response as AuthResponse
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout')
  },

  async getMe(): Promise<User> {
    // api.get returns data directly due to response interceptor
    const response = await api.get('/auth/me')
    return response as User
  },
}
