import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDocumentsStore } from '../documentsStore'
import { documentService } from '../../services/documentService'
import type { Document } from '../../types/document.types'

vi.mock('../../services/documentService')
vi.mock('@/utils/logger', () => ({
  logger: {
    logStoreAction: vi.fn(),
    logStoreError: vi.fn(),
    info: vi.fn(),
    debug: vi.fn(),
  },
}))

describe('documentsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  const mockDocument: Document = {
    id: '123',
    user_id: 'user-1',
    title: 'Test Document',
    file_name: 'test.pdf',
    file_path: '/uploads/test.pdf',
    file_type: 'pdf',
    file_size: 1024,
    status: 'uploaded',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }

  describe('uploadDocument', () => {
    it('should upload document successfully', async () => {
      const store = useDocumentsStore()
      const mockFile = new File(['test'], 'test.pdf')

      vi.mocked(documentService.upload).mockResolvedValue(mockDocument)

      const result = await store.uploadDocument({ file: mockFile, title: 'Test Document' })

      expect(documentService.upload).toHaveBeenCalledWith({
        file: mockFile,
        title: 'Test Document',
      })
      expect(result).toEqual(mockDocument)
      expect(store.documents).toHaveLength(1)
      expect(store.documents[0]).toEqual(mockDocument)
      expect(store.total).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should handle upload error', async () => {
      const store = useDocumentsStore()
      const mockFile = new File(['test'], 'test.pdf')
      const mockError = new Error('Upload failed')

      vi.mocked(documentService.upload).mockRejectedValue(mockError)

      await expect(
        store.uploadDocument({ file: mockFile, title: 'Test' })
      ).rejects.toThrow('Upload failed')

      expect(store.error).toBe('Upload failed')
      expect(store.documents).toHaveLength(0)
      expect(store.loading).toBe(false)
    })

    it('should add uploaded document to beginning of list', async () => {
      const store = useDocumentsStore()
      const existingDoc = { ...mockDocument, id: '1' }
      const newDoc = { ...mockDocument, id: '2' }
      store.documents = [existingDoc]

      vi.mocked(documentService.upload).mockResolvedValue(newDoc)

      await store.uploadDocument({
        file: new File(['test'], 'new.pdf'),
        title: 'New Doc',
      })

      expect(store.documents[0].id).toBe('2')
      expect(store.documents[1].id).toBe('1')
    })
  })

  describe('fetchDocuments', () => {
    it('should fetch documents successfully', async () => {
      const store = useDocumentsStore()
      const mockResponse = {
        documents: [mockDocument],
        total: 1,
        page: 1,
        page_size: 20,
      }

      vi.mocked(documentService.list).mockResolvedValue(mockResponse)

      await store.fetchDocuments(1)

      expect(documentService.list).toHaveBeenCalledWith(1, 20)
      expect(store.documents).toEqual([mockDocument])
      expect(store.total).toBe(1)
      expect(store.currentPage).toBe(1)
      expect(store.totalPages).toBe(1)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should calculate total pages correctly', async () => {
      const store = useDocumentsStore()
      const mockResponse = {
        documents: Array(20).fill(mockDocument),
        total: 45,
        page: 1,
        page_size: 20,
      }

      vi.mocked(documentService.list).mockResolvedValue(mockResponse)

      await store.fetchDocuments(1)

      expect(store.totalPages).toBe(3) // 45 / 20 = 2.25 -> 3
    })

    it('should handle fetch error', async () => {
      const store = useDocumentsStore()

      vi.mocked(documentService.list).mockRejectedValue(new Error('Fetch failed'))

      await expect(store.fetchDocuments(1)).rejects.toThrow('Fetch failed')

      expect(store.error).toBe('Fetch failed')
      expect(store.loading).toBe(false)
    })

    it('should handle empty documents list', async () => {
      const store = useDocumentsStore()
      const mockResponse = {
        documents: [],
        total: 0,
        page: 1,
        page_size: 20,
      }

      vi.mocked(documentService.list).mockResolvedValue(mockResponse)

      await store.fetchDocuments(1)

      expect(store.documents).toEqual([])
      expect(store.total).toBe(0)
    })
  })

  describe('fetchDocument', () => {
    it('should fetch single document successfully', async () => {
      const store = useDocumentsStore()

      vi.mocked(documentService.getById).mockResolvedValue(mockDocument)

      const result = await store.fetchDocument('123')

      expect(documentService.getById).toHaveBeenCalledWith('123')
      expect(result).toEqual(mockDocument)
      expect(store.currentDocument).toEqual(mockDocument)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should handle fetch document error', async () => {
      const store = useDocumentsStore()

      vi.mocked(documentService.getById).mockRejectedValue(new Error('Document not found'))

      await expect(store.fetchDocument('invalid-id')).rejects.toThrow('Document not found')

      expect(store.error).toBe('Document not found')
      expect(store.currentDocument).toBeNull()
    })
  })

  describe('deleteDocument', () => {
    it('should delete document successfully', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]
      store.total = 1

      vi.mocked(documentService.delete).mockResolvedValue(undefined)

      await store.deleteDocument('123')

      expect(documentService.delete).toHaveBeenCalledWith('123')
      expect(store.documents).toHaveLength(0)
      expect(store.total).toBe(0)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should clear current document if deleted', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]
      store.currentDocument = mockDocument
      store.total = 1

      vi.mocked(documentService.delete).mockResolvedValue(undefined)

      await store.deleteDocument('123')

      expect(store.currentDocument).toBeNull()
    })

    it('should not clear current document if different document deleted', async () => {
      const store = useDocumentsStore()
      const doc1 = { ...mockDocument, id: '1' }
      const doc2 = { ...mockDocument, id: '2' }
      store.documents = [doc1, doc2]
      store.currentDocument = doc1
      store.total = 2

      vi.mocked(documentService.delete).mockResolvedValue(undefined)

      await store.deleteDocument('2')

      expect(store.currentDocument).toEqual(doc1)
      expect(store.documents).toHaveLength(1)
      expect(store.documents[0].id).toBe('1')
    })

    it('should handle delete error', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]

      vi.mocked(documentService.delete).mockRejectedValue(new Error('Delete failed'))

      await expect(store.deleteDocument('123')).rejects.toThrow('Delete failed')

      expect(store.error).toBe('Delete failed')
      expect(store.documents).toHaveLength(1)
    })
  })

  describe('parseDocument', () => {
    it('should parse document successfully', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]
      const mockParseResult = {
        id: '123',
        parsed_text: 'Parsed content',
        status: 'parsed',
        text_preview: 'Parsed...',
      }

      vi.mocked(documentService.parse).mockResolvedValue(mockParseResult)

      const result = await store.parseDocument('123')

      expect(documentService.parse).toHaveBeenCalledWith('123')
      expect(result).toEqual(mockParseResult)
      expect(store.documents[0].status).toBe('parsed')
      expect(store.documents[0].parsed_text).toBe('Parsed content')
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should update current document if parsing current', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]
      store.currentDocument = mockDocument
      const mockParseResult = {
        id: '123',
        parsed_text: 'Parsed content',
        status: 'parsed',
        text_preview: 'Parsed...',
      }

      vi.mocked(documentService.parse).mockResolvedValue(mockParseResult)

      await store.parseDocument('123')

      expect(store.currentDocument?.status).toBe('parsed')
      expect(store.currentDocument?.parsed_text).toBe('Parsed content')
    })

    it('should handle parse error', async () => {
      const store = useDocumentsStore()
      store.documents = [mockDocument]

      vi.mocked(documentService.parse).mockRejectedValue(new Error('Parsing failed'))

      await expect(store.parseDocument('123')).rejects.toThrow('Parsing failed')

      expect(store.error).toBe('Parsing failed')
    })
  })

  describe('clearError', () => {
    it('should clear error', () => {
      const store = useDocumentsStore()
      store.error = 'Some error'

      store.clearError()

      expect(store.error).toBeNull()
    })
  })

  // Negative tests
  describe('negative scenarios', () => {
    it('should handle network timeout', async () => {
      const store = useDocumentsStore()
      vi.mocked(documentService.list).mockRejectedValue(new Error('Network timeout'))

      await expect(store.fetchDocuments(1)).rejects.toThrow('Network timeout')
    })

    it('should handle unauthorized access', async () => {
      const store = useDocumentsStore()
      vi.mocked(documentService.list).mockRejectedValue(new Error('Unauthorized'))

      await expect(store.fetchDocuments(1)).rejects.toThrow('Unauthorized')
    })

    it('should handle concurrent uploads', async () => {
      const store = useDocumentsStore()
      const file1 = new File(['test1'], 'test1.pdf')
      const file2 = new File(['test2'], 'test2.pdf')
      const doc1 = { ...mockDocument, id: '1' }
      const doc2 = { ...mockDocument, id: '2' }

      vi.mocked(documentService.upload)
        .mockResolvedValueOnce(doc1)
        .mockResolvedValueOnce(doc2)

      await Promise.all([
        store.uploadDocument({ file: file1, title: 'Test 1' }),
        store.uploadDocument({ file: file2, title: 'Test 2' }),
      ])

      expect(store.documents).toHaveLength(2)
      expect(store.total).toBe(2)
    })

    it('should handle parsing non-existent document', async () => {
      const store = useDocumentsStore()
      store.documents = []

      const mockParseResult = {
        id: 'non-existent',
        parsed_text: 'Parsed content',
        status: 'parsed',
        text_preview: 'Parsed...',
      }

      vi.mocked(documentService.parse).mockResolvedValue(mockParseResult)

      await store.parseDocument('non-existent')

      // Should not crash, just not update any document in list
      expect(store.documents).toHaveLength(0)
    })
  })
})
