import { describe, it, expect, vi, beforeEach } from 'vitest'
import { documentService } from '../documentService'
import api from '@/services/api'

vi.mock('@/services/api', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('documentService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('upload', () => {
    it('should upload document successfully', async () => {
      const mockFile = new File(['test content'], 'test.pdf', { type: 'application/pdf' })
      const mockResponse = {
        id: '123',
        title: 'test.pdf',
        file_name: 'test.pdf',
        file_type: 'pdf',
        file_size: 1024,
        status: 'uploaded',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const result = await documentService.upload({ file: mockFile, title: 'test.pdf' })

      expect(api.post).toHaveBeenCalledWith(
        '/documents',
        expect.any(FormData),
        expect.objectContaining({
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
      )
      expect(result).toEqual(mockResponse)
    })

    it('should upload document without title', async () => {
      const mockFile = new File(['test content'], 'test.pdf', { type: 'application/pdf' })
      const mockResponse = {
        id: '123',
        title: 'test.pdf',
        file_name: 'test.pdf',
        file_type: 'pdf',
        file_size: 1024,
        status: 'uploaded',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      await documentService.upload({ file: mockFile, title: '' })

      const formData = vi.mocked(api.post).mock.calls[0][1] as FormData
      expect(formData.has('file')).toBe(true)
      expect(formData.has('title')).toBe(false)
    })

    it('should handle upload error', async () => {
      const mockFile = new File(['test content'], 'test.pdf', { type: 'application/pdf' })
      const mockError = new Error('Upload failed')

      vi.mocked(api.post).mockRejectedValue(mockError)

      await expect(documentService.upload({ file: mockFile, title: 'test.pdf' })).rejects.toThrow(
        'Upload failed'
      )
    })
  })

  describe('list', () => {
    it('should fetch documents list with default pagination', async () => {
      const mockResponse = {
        documents: [
          {
            id: '1',
            title: 'Doc 1',
            file_name: 'doc1.pdf',
            file_type: 'pdf',
            file_size: 1024,
            status: 'uploaded',
            created_at: '2024-01-01T00:00:00Z',
          },
        ],
        total: 1,
        page: 1,
        page_size: 20,
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)

      const result = await documentService.list()

      expect(api.get).toHaveBeenCalledWith('/documents', {
        params: { page: 1, page_size: 20 },
      })
      expect(result).toEqual(mockResponse)
    })

    it('should fetch documents list with custom pagination', async () => {
      const mockResponse = {
        documents: [],
        total: 0,
        page: 2,
        page_size: 10,
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)

      await documentService.list(2, 10)

      expect(api.get).toHaveBeenCalledWith('/documents', {
        params: { page: 2, page_size: 10 },
      })
    })

    it('should handle list error', async () => {
      vi.mocked(api.get).mockRejectedValue(new Error('Failed to fetch'))

      await expect(documentService.list()).rejects.toThrow('Failed to fetch')
    })
  })

  describe('getById', () => {
    it('should fetch document by ID successfully', async () => {
      const mockDocument = {
        id: '123',
        title: 'Test Doc',
        file_name: 'test.pdf',
        file_type: 'pdf',
        file_size: 1024,
        status: 'uploaded',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockDocument)

      const result = await documentService.getById('123')

      expect(api.get).toHaveBeenCalledWith('/documents/123')
      expect(result).toEqual(mockDocument)
    })

    it('should handle getById error (document not found)', async () => {
      vi.mocked(api.get).mockRejectedValue(new Error('Document not found'))

      await expect(documentService.getById('invalid-id')).rejects.toThrow('Document not found')
    })
  })

  describe('delete', () => {
    it('should delete document successfully', async () => {
      vi.mocked(api.delete).mockResolvedValue(undefined)

      await documentService.delete('123')

      expect(api.delete).toHaveBeenCalledWith('/documents/123')
    })

    it('should handle delete error', async () => {
      vi.mocked(api.delete).mockRejectedValue(new Error('Delete failed'))

      await expect(documentService.delete('123')).rejects.toThrow('Delete failed')
    })
  })

  describe('parse', () => {
    it('should parse document successfully', async () => {
      const mockResponse = {
        id: '123',
        parsed_text: 'This is parsed text content',
        status: 'parsed',
        text_preview: 'This is parsed...',
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const result = await documentService.parse('123')

      expect(api.post).toHaveBeenCalledWith('/documents/123/parse')
      expect(result).toEqual(mockResponse)
    })

    it('should handle parse error', async () => {
      vi.mocked(api.post).mockRejectedValue(new Error('Parsing failed'))

      await expect(documentService.parse('123')).rejects.toThrow('Parsing failed')
    })
  })

  // Negative tests
  describe('negative scenarios', () => {
    it('should handle network error during upload', async () => {
      const mockFile = new File(['test'], 'test.pdf')
      vi.mocked(api.post).mockRejectedValue(new Error('Network error'))

      await expect(documentService.upload({ file: mockFile, title: 'test' })).rejects.toThrow(
        'Network error'
      )
    })

    it('should handle empty response from list', async () => {
      vi.mocked(api.get).mockResolvedValue({
        documents: [],
        total: 0,
        page: 1,
        page_size: 20,
      })

      const result = await documentService.list()

      expect(result.documents).toEqual([])
      expect(result.total).toBe(0)
    })

    it('should handle unauthorized error', async () => {
      vi.mocked(api.get).mockRejectedValue(new Error('Unauthorized'))

      await expect(documentService.getById('123')).rejects.toThrow('Unauthorized')
    })

    it('should handle large file upload', async () => {
      const largeContent = 'a'.repeat(50 * 1024 * 1024) // 50MB
      const mockFile = new File([largeContent], 'large.pdf')
      const mockResponse = {
        id: '123',
        title: 'large.pdf',
        file_name: 'large.pdf',
        file_type: 'pdf',
        file_size: 50 * 1024 * 1024,
        status: 'uploaded',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const result = await documentService.upload({ file: mockFile, title: 'large.pdf' })

      expect(result.file_size).toBe(50 * 1024 * 1024)
    })
  })
})
