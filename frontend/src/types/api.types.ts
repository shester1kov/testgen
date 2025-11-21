export interface ApiResponse<T = unknown> {
  data: T
  message?: string
  success?: boolean
}

// New backend error structure
export interface ApiError {
  message: string
  code?: string // Error code from backend (e.g., "INVALID_INPUT", "USER_NOT_FOUND")
  status?: number
}

export interface PaginationParams {
  page?: number
  limit?: number
  offset?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page?: number
  page_size?: number
  limit?: number
  offset?: number
}
