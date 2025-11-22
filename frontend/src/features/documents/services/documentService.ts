import api from '@/services/api'
import type {
  Document,
  DocumentUploadRequest,
  DocumentParseRequest,
} from '../types/document.types'

const DOCUMENTS_BASE_URL = '/documents'

export const documentService = {
  /**
   * Upload a new document
   */
  async upload(request: DocumentUploadRequest): Promise<Document> {
    const formData = new FormData()
    formData.append('file', request.file)
    if (request.title) {
      formData.append('title', request.title)
    }

    const response = await api.post<Document>(DOCUMENTS_BASE_URL, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response
  },

  /**
   * Get list of documents with pagination
   */
  async list(page = 1, pageSize = 20): Promise<{
    documents: Document[]
    total: number
    page: number
    page_size: number
  }> {
    const response = await api.get<{
      documents: Document[]
      total: number
      page: number
      page_size: number
    }>(DOCUMENTS_BASE_URL, {
      params: { page, page_size: pageSize },
    })
    return response
  },

  /**
   * Get document by ID
   */
  async getById(id: string): Promise<Document> {
    const response = await api.get<Document>(`${DOCUMENTS_BASE_URL}/${id}`)
    return response
  },

  /**
   * Delete document
   */
  async delete(id: string): Promise<void> {
    await api.delete(`${DOCUMENTS_BASE_URL}/${id}`)
  },

  /**
   * Parse document to extract text
   */
  async parse(id: string): Promise<{
    id: string
    parsed_text: string
    status: string
    text_preview: string
  }> {
    const response = await api.post<{
      id: string
      parsed_text: string
      status: string
      text_preview: string
    }>(`${DOCUMENTS_BASE_URL}/${id}/parse`)
    return response
  },
}
