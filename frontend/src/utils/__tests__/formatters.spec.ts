import { describe, it, expect } from 'vitest'
import { formatDate, formatRelativeTime, truncateText } from '../formatters'

describe('Formatters', () => {
  describe('formatDate', () => {
    it('should format date correctly', () => {
      const date = new Date('2024-01-15T10:30:00')
      const formatted = formatDate(date)

      expect(formatted).toContain('January')
      expect(formatted).toContain('15')
      expect(formatted).toContain('2024')
    })
  })

  describe('formatRelativeTime', () => {
    it('should return "just now" for recent dates', () => {
      const now = new Date()
      expect(formatRelativeTime(now)).toBe('just now')
    })

    it('should return minutes ago', () => {
      const date = new Date(Date.now() - 5 * 60 * 1000) // 5 minutes ago
      expect(formatRelativeTime(date)).toBe('5 minutes ago')
    })

    it('should return hours ago', () => {
      const date = new Date(Date.now() - 2 * 60 * 60 * 1000) // 2 hours ago
      expect(formatRelativeTime(date)).toBe('2 hours ago')
    })
  })

  describe('truncateText', () => {
    it('should truncate long text', () => {
      const text = 'This is a very long text that should be truncated'
      expect(truncateText(text, 20)).toBe('This is a very long ...')
    })

    it('should not truncate short text', () => {
      const text = 'Short text'
      expect(truncateText(text, 20)).toBe('Short text')
    })
  })
})
