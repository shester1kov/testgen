import api from './api'
import type { Document, DocumentUploadRequest } from '@/features/documents/types/document.types'
import type { ApiResponse, PaginatedResponse } from '@/types/api.types'

export const documentService = {
  async uploadDocument(data: DocumentUploadRequest): Promise<Document> {
    const formData = new FormData()
    formData.append('title', data.title)
    formData.append('file', data.file)

    const response = await api.post<ApiResponse<Document>>('/documents', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data as Document
  },

  async getDocuments(page = 1, limit = 10): Promise<PaginatedResponse<Document>> {
    const response = await api.get<ApiResponse<PaginatedResponse<Document>>>('/documents', {
      params: { page, limit },
    })
    return response.data as PaginatedResponse<Document>
  },

  async getDocument(id: string): Promise<Document> {
    const response = await api.get<ApiResponse<Document>>(`/documents/${id}`)
    return response.data as Document
  },

  async deleteDocument(id: string): Promise<void> {
    await api.delete(`/documents/${id}`)
  },

  async parseDocument(id: string): Promise<Document> {
    const response = await api.post<ApiResponse<Document>>(`/documents/${id}/parse`)
    return response.data as Document
  },
}
