import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Document, DocumentUploadRequest } from '../types/document.types'
import { documentService } from '../services/documentService'
import { logger } from '@/utils/logger'
import { isDesignMode, getMockDocuments } from '@/utils/designMode'

export const useDocumentsStore = defineStore('documents', () => {
  // State
  const documents = ref<Document[]>([])
  const currentDocument = ref<Document | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const currentPage = ref(1)
  const totalPages = ref(0)

  // Actions
  async function uploadDocument(data: DocumentUploadRequest) {
    logger.logStoreAction('documentsStore', 'uploadDocument', {
      fileName: data.file.name,
      fileSize: data.file.size,
    })
    loading.value = true
    error.value = null

    try {
      const document = await documentService.upload(data)
      documents.value.unshift(document)
      total.value++
      logger.info('Document uploaded successfully', 'documentsStore', {
        documentId: document.id,
      })
      return document
    } catch (err: any) {
      error.value = err.message || 'Failed to upload document'
      logger.logStoreError('documentsStore', 'uploadDocument', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDocuments(page = 1, limit = 20) {
    logger.logStoreAction('documentsStore', 'fetchDocuments', { page, limit })
    loading.value = true
    error.value = null

    try {
      // Design mode: return mock documents
      if (isDesignMode()) {
        const mockDocs = getMockDocuments() as any[]
        documents.value = mockDocs
        total.value = mockDocs.length
        currentPage.value = 1
        totalPages.value = 1
        logger.info('Design mode: Mock documents returned', 'documentsStore')
        loading.value = false
        return { documents: mockDocs, total: mockDocs.length, page: 1, page_size: limit }
      }

      const response = await documentService.list(page, limit)
      documents.value = response.documents
      total.value = response.total
      currentPage.value = response.page
      totalPages.value = Math.ceil(response.total / limit)
      logger.debug('Documents fetched successfully', 'documentsStore', {
        count: response.documents.length,
        total: response.total,
      })
      return response
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch documents'
      logger.logStoreError('documentsStore', 'fetchDocuments', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDocument(id: string) {
    logger.logStoreAction('documentsStore', 'fetchDocument', { id })
    loading.value = true
    error.value = null

    try {
      const document = await documentService.getById(id)
      currentDocument.value = document
      logger.debug('Document fetched successfully', 'documentsStore', { id })
      return document
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch document'
      logger.logStoreError('documentsStore', 'fetchDocument', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteDocument(id: string) {
    logger.logStoreAction('documentsStore', 'deleteDocument', { id })
    loading.value = true
    error.value = null

    try {
      await documentService.delete(id)
      documents.value = documents.value.filter((doc) => doc.id !== id)
      total.value--
      if (currentDocument.value?.id === id) {
        currentDocument.value = null
      }
      logger.info('Document deleted successfully', 'documentsStore', { id })
    } catch (err: any) {
      error.value = err.message || 'Failed to delete document'
      logger.logStoreError('documentsStore', 'deleteDocument', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function parseDocument(id: string) {
    logger.logStoreAction('documentsStore', 'parseDocument', { id })
    loading.value = true
    error.value = null

    try {
      const result = await documentService.parse(id)

      // Update document in list
      const index = documents.value.findIndex((doc) => doc.id === id)
      if (index !== -1) {
        documents.value[index].status = result.status as any
        documents.value[index].parsed_text = result.parsed_text
      }

      // Update current document if it's the one being parsed
      if (currentDocument.value?.id === id) {
        currentDocument.value.status = result.status as any
        currentDocument.value.parsed_text = result.parsed_text
      }

      logger.info('Document parsed successfully', 'documentsStore', { id })
      return result
    } catch (err: any) {
      error.value = err.message || 'Failed to parse document'
      logger.logStoreError('documentsStore', 'parseDocument', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    // State
    documents,
    currentDocument,
    loading,
    error,
    total,
    currentPage,
    totalPages,
    // Actions
    uploadDocument,
    fetchDocuments,
    fetchDocument,
    deleteDocument,
    parseDocument,
    clearError,
  }
})
