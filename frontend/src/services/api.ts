import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types/api.types'
import { logger } from '@/utils/logger'

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
    // Log outgoing request
    logger.logRequest(
      config.method?.toUpperCase() || 'GET',
      config.url || '',
      config.data
    )

    // Cookies are automatically sent with withCredentials: true
    // No need to manually add Authorization header for cookie-based auth
    return config
  },
  (error: AxiosError) => {
    logger.logError('REQUEST', error.config?.url || '', error)
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  response => {
    // Log successful response
    logger.logResponse(
      response.config.method?.toUpperCase() || 'GET',
      response.config.url || '',
      response.status,
      response.data
    )
    return response.data
  },
  (error: AxiosError<any>) => {
    // Log error response
    if (error.response) {
      logger.logResponse(
        error.config?.method?.toUpperCase() || 'GET',
        error.config?.url || '',
        error.response.status,
        error.response.data
      )
    } else {
      logger.logError(
        error.config?.method?.toUpperCase() || 'GET',
        error.config?.url || '',
        error
      )
    }

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
        logger.warn('Unauthorized - redirecting to login', 'AUTH')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }

      return Promise.reject(apiError)
    } else if (error.request) {
      // Network error
      const networkError: ApiError = {
        message: 'Network error. Please check your connection.',
        status: 0,
      }
      logger.error('Network error', 'HTTP', error)
      return Promise.reject(networkError)
    } else {
      const unexpectedError: ApiError = {
        message: error.message || 'An unexpected error occurred',
      }
      logger.error('Unexpected error', 'HTTP', error)
      return Promise.reject(unexpectedError)
    }
  }
)

export default api
