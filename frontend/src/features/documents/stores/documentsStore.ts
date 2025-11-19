import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Document, DocumentUploadRequest } from '../types/document.types'
import { documentService } from '@/services/documentService'

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
    loading.value = true
    error.value = null

    try {
      const document = await documentService.uploadDocument(data)
      documents.value.unshift(document)
      return document
    } catch (err: any) {
      error.value = err.message || 'Failed to upload document'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDocuments(page = 1, limit = 10) {
    loading.value = true
    error.value = null

    try {
      const response = await documentService.getDocuments(page, limit)
      documents.value = response.data
      total.value = response.total
      currentPage.value = response.page
      totalPages.value = response.totalPages
      return response
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch documents'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDocument(id: string) {
    loading.value = true
    error.value = null

    try {
      const document = await documentService.getDocument(id)
      currentDocument.value = document
      return document
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch document'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteDocument(id: string) {
    loading.value = true
    error.value = null

    try {
      await documentService.deleteDocument(id)
      documents.value = documents.value.filter(doc => doc.id !== id)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete document'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function parseDocument(id: string) {
    loading.value = true
    error.value = null

    try {
      const document = await documentService.parseDocument(id)
      // Update document in list
      const index = documents.value.findIndex(doc => doc.id === id)
      if (index !== -1) {
        documents.value[index] = document
      }
      if (currentDocument.value?.id === id) {
        currentDocument.value = document
      }
      return document
    } catch (err: any) {
      error.value = err.message || 'Failed to parse document'
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
