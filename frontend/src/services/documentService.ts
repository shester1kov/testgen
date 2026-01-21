import api from './api'
import type { Document, DocumentUploadRequest } from '@/features/documents/types/document.types'

// Backend response for document list
interface DocumentListResponse {
  documents: Document[]
  total: number
  page: number
  page_size: number
}

export const documentService = {
  async uploadDocument(data: DocumentUploadRequest): Promise<Document> {
    const formData = new FormData()
    formData.append('title', data.title)
    formData.append('file', data.file)

    const response = await api.post('/documents', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response as unknown as Document
  },

  async getDocuments(page = 1, limit = 10): Promise<{ documents: Document[]; total: number; page: number; page_size: number }> {
    const response = await api.get('/documents', {
      params: { page, limit },
    })
    return response as unknown as DocumentListResponse
  },

  async getDocument(id: string): Promise<Document> {
    const response = await api.get(`/documents/${id}`)
    return response as unknown as Document
  },

  async deleteDocument(id: string): Promise<void> {
    await api.delete(`/documents/${id}`)
  },

  async parseDocument(id: string): Promise<Document> {
    const response = await api.post(`/documents/${id}/parse`)
    return response as unknown as Document
  },
}
