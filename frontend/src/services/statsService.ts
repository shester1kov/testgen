import api from './api'

export interface DashboardStats {
  documents_count: number
  tests_count: number
  questions_count: number
}

export const statsService = {
  async getDashboardStats(): Promise<DashboardStats> {
    const response = await api.get<DashboardStats>('/stats/dashboard')
    return response as DashboardStats
  },
}

export default statsService
