import api from './api'
import type {
  Test,
  TestGenerationRequest,
  TestExportRequest,
  MoodleSyncRequest,
} from '@/features/tests/types/test.types'
import type { ApiResponse, PaginatedResponse } from '@/types/api.types'

export const testService = {
  async createTest(data: Partial<Test>): Promise<Test> {
    const response = await api.post<ApiResponse<Test>>('/tests', data)
    return response.data as Test
  },

  async getTests(page = 1, limit = 10): Promise<PaginatedResponse<Test>> {
    const response = await api.get<ApiResponse<PaginatedResponse<Test>>>('/tests', {
      params: { page, limit },
    })
    return response.data as PaginatedResponse<Test>
  },

  async getTest(id: string): Promise<Test> {
    const response = await api.get<ApiResponse<Test>>(`/tests/${id}`)
    return response.data as Test
  },

  async updateTest(id: string, data: Partial<Test>): Promise<Test> {
    const response = await api.put<ApiResponse<Test>>(`/tests/${id}`, data)
    return response.data as Test
  },

  async deleteTest(id: string): Promise<void> {
    await api.delete(`/tests/${id}`)
  },

  async generateTest(data: TestGenerationRequest): Promise<Test> {
    const response = await api.post<ApiResponse<Test>>('/tests/generate', data)
    return response.data as Test
  },

  async exportTest(data: TestExportRequest): Promise<Blob> {
    const response = await api.post(`/tests/${data.test_id}/export`, data, {
      responseType: 'blob',
    })
    return response as unknown as Blob
  },

  async syncToMoodle(data: MoodleSyncRequest): Promise<Test> {
    const response = await api.post<ApiResponse<Test>>(`/tests/${data.test_id}/sync`, data)
    return response.data as Test
  },
}
