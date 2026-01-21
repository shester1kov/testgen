import api from './api'
import type {
  Test,
  TestGenerationRequest,
  TestExportRequest,
  MoodleSyncRequest,
} from '@/features/tests/types/test.types'
import type { PaginatedResponse } from '@/types/api.types'
import type { FormatInfo, ExportFormat } from '@/features/tests/types/export.types'

export const testService = {
  async createTest(data: Partial<Test>): Promise<Test> {
    const response = await api.post('/tests', data)
    return response as unknown as Test
  },

  async getTests(page = 1, limit = 10): Promise<PaginatedResponse<Test>> {
    const response: any = await api.get('/tests', {
      params: { page, page_size: limit },
    })
    // Backend returns { tests: [...], total: ..., page: ..., page_size: ... }
    // Transform to { data: [...], total: ... } format
    return {
      data: response.tests || [],
      total: response.total || 0,
      page: response.page,
      page_size: response.page_size,
    }
  },

  async getTest(id: string): Promise<Test> {
    const response = await api.get(`/tests/${id}`)
    return response as unknown as Test
  },

  async deleteTest(id: string): Promise<void> {
    await api.delete(`/tests/${id}`)
  },

  async generateTest(data: TestGenerationRequest): Promise<Test> {
    const response = await api.post('/tests/generate', data)
    return response as unknown as Test
  },

  async exportTest(data: TestExportRequest): Promise<Blob> {
    const response = await api.get(`/moodle/tests/${data.test_id}/export`, {
      responseType: 'blob',
    })
    return response as unknown as Blob
  },

  async syncToMoodle(data: MoodleSyncRequest): Promise<Test> {
    const response = await api.post(`/moodle/tests/${data.test_id}/sync`, data)
    return response as unknown as Test
  },

  async getExportFormats(): Promise<FormatInfo[]> {
    const response: any = await api.get('/tests/export/formats')
    return response as FormatInfo[]
  },

  async exportToFormat(testId: string, format: ExportFormat): Promise<void> {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
    const url = `${API_BASE_URL}/tests/${testId}/export?format=${format}`

    const response = await fetch(url, {
      method: 'GET',
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Export failed: ${response.statusText}`)
    }

    // Get filename from Content-Disposition header
    const disposition = response.headers.get('Content-Disposition')
    let filename = `test.${format}`
    if (disposition && disposition.includes('filename=')) {
      const filenameMatch = disposition.match(/filename=(.+)/)
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1].replace(/['"]/g, '')
      }
    }

    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  },
}

export default testService
