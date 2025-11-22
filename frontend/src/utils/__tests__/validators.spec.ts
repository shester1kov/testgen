import { describe, it, expect } from 'vitest'
import { validateFileSize, validateFileType, formatFileSize } from '../validators'

describe('Validators', () => {
  describe('validateFileSize', () => {
    it('should accept files under max size', () => {
      const file = new File(['content'], 'test.pdf', { type: 'application/pdf' })
      Object.defineProperty(file, 'size', { value: 1024 * 1024 }) // 1MB

      expect(validateFileSize(file)).toBe(true)
    })

    it('should reject files over max size', () => {
      const file = new File(['content'], 'test.pdf', { type: 'application/pdf' })
      Object.defineProperty(file, 'size', { value: 100 * 1024 * 1024 }) // 100MB

      expect(validateFileSize(file)).toBe(false)
    })
  })

  describe('validateFileType', () => {
    it('should accept supported file types', () => {
      const validTypes = ['test.pdf', 'test.docx', 'test.pptx', 'test.txt']

      validTypes.forEach(filename => {
        const file = new File(['content'], filename)
        expect(validateFileType(file)).toBe(true)
      })
    })

    it('should reject unsupported file types', () => {
      const invalidTypes = ['test.exe', 'test.jpg', 'test.zip']

      invalidTypes.forEach(filename => {
        const file = new File(['content'], filename)
        expect(validateFileType(file)).toBe(false)
      })
    })
  })

  describe('formatFileSize', () => {
    it('should format bytes correctly', () => {
      expect(formatFileSize(0)).toBe('0 Bytes')
      expect(formatFileSize(1024)).toBe('1 KB')
      expect(formatFileSize(1024 * 1024)).toBe('1 MB')
      expect(formatFileSize(1536 * 1024)).toBe('1.5 MB')
    })
  })
})
