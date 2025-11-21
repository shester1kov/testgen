import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types/api.types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// Create axios instance
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 seconds
  withCredentials: true, // Enable sending cookies
})

// Request interceptor
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Cookies are automatically sent with withCredentials: true
    // No need to manually add Authorization header for cookie-based auth
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  response => {
    return response.data
  },
  (error: AxiosError<any>) => {
    // Handle errors matching new backend DTO structure
    if (error.response) {
      const errorData = error.response.data

      // New backend error structure: { error: { code: string, message: string } }
      const apiError: ApiError = {
        message: errorData?.error?.message || errorData?.message || 'An error occurred',
        code: errorData?.error?.code,
        status: error.response.status,
      }

      // Handle 401 Unauthorized
      if (error.response.status === 401) {
        localStorage.removeItem('user')
        window.location.href = '/login'
      }

      return Promise.reject(apiError)
    } else if (error.request) {
      // Network error
      return Promise.reject({
        message: 'Network error. Please check your connection.',
        status: 0,
      } as ApiError)
    } else {
      return Promise.reject({
        message: error.message || 'An unexpected error occurred',
      } as ApiError)
    }
  }
)

export default api
