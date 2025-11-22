import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { logger, LogLevel } from '../logger'

describe('Logger Utility', () => {
  // Mock console methods
  const consoleSpy = {
    debug: vi.spyOn(console, 'debug').mockImplementation(() => {}),
    info: vi.spyOn(console, 'info').mockImplementation(() => {}),
    warn: vi.spyOn(console, 'warn').mockImplementation(() => {}),
    error: vi.spyOn(console, 'error').mockImplementation(() => {}),
    log: vi.spyOn(console, 'log').mockImplementation(() => {}),
  }

  beforeEach(() => {
    // Clear all mocks before each test
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('POSITIVE: Basic logging', () => {
    it('should log debug messages', () => {
      logger.debug('Debug message', 'TestContext')

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('[DEBUG]')
      expect(callArgs[0]).toContain('[TestContext]')
      expect(callArgs[0]).toContain('Debug message')
    })

    it('should log info messages', () => {
      logger.info('Info message')

      expect(consoleSpy.info).toHaveBeenCalled()
      const callArgs = consoleSpy.info.mock.calls[0]
      expect(callArgs[0]).toContain('[INFO]')
      expect(callArgs[0]).toContain('Info message')
    })

    it('should log warn messages', () => {
      logger.warn('Warning message')

      expect(consoleSpy.warn).toHaveBeenCalled()
      const callArgs = consoleSpy.warn.mock.calls[0]
      expect(callArgs[0]).toContain('[WARN]')
      expect(callArgs[0]).toContain('Warning message')
    })

    it('should log error messages with error object', () => {
      const testError = new Error('Test error')
      logger.error('Error message', 'ErrorContext', testError)

      expect(consoleSpy.error).toHaveBeenCalled()
      const callArgs = consoleSpy.error.mock.calls[0]
      expect(callArgs[0]).toContain('[ERROR]')
      expect(callArgs[0]).toContain('[ErrorContext]')
      expect(callArgs[0]).toContain('Error message')
      expect(callArgs[2]).toBe(testError)
    })
  })

  describe('NEGATIVE: Empty or null inputs', () => {
    it('should handle empty message', () => {
      logger.info('')

      expect(consoleSpy.info).toHaveBeenCalled()
      const callArgs = consoleSpy.info.mock.calls[0]
      expect(callArgs[0]).toContain('[INFO]')
    })

    it('should handle undefined context', () => {
      logger.debug('Message with undefined context', undefined)

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('[DEBUG]')
      expect(callArgs[0]).toContain('Message with undefined context')
    })

    it('should handle null data', () => {
      logger.info('Message with null data', 'Context', null)

      expect(consoleSpy.info).toHaveBeenCalled()
    })

    it('should handle undefined error', () => {
      logger.error('Error without error object', 'Context', undefined)

      expect(consoleSpy.error).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Invalid input types', () => {
    it('should handle numeric message', () => {
      // @ts-expect-error Testing invalid input
      logger.info(12345)

      expect(consoleSpy.info).toHaveBeenCalled()
    })

    it('should handle object as message', () => {
      // @ts-expect-error Testing invalid input
      logger.info({ message: 'test' })

      expect(consoleSpy.info).toHaveBeenCalled()
    })

    it('should handle array as context', () => {
      // @ts-expect-error Testing invalid input
      logger.warn('Message', ['invalid', 'context'])

      expect(consoleSpy.warn).toHaveBeenCalled()
    })
  })

  describe('POSITIVE: HTTP request logging', () => {
    it('should log HTTP requests', () => {
      logger.logRequest('GET', '/api/users', { page: 1 })

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('GET /api/users')
      expect(callArgs[0]).toContain('[HTTP]')
    })

    it('should log successful HTTP responses', () => {
      logger.logResponse('POST', '/api/login', 200, { token: 'abc' })

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('POST /api/login - 200')
    })

    it('should log HTTP errors', () => {
      const error = new Error('Network error')
      logger.logError('GET', '/api/data', error)

      expect(consoleSpy.error).toHaveBeenCalled()
      const callArgs = consoleSpy.error.mock.calls[0]
      expect(callArgs[0]).toContain('GET /api/data failed')
      expect(callArgs[0]).toContain('[HTTP]')
    })

    it('should log HTTP errors as warnings for 4xx status codes', () => {
      logger.logResponse('DELETE', '/api/users/123', 404)

      expect(consoleSpy.warn).toHaveBeenCalled()
      const callArgs = consoleSpy.warn.mock.calls[0]
      expect(callArgs[0]).toContain('DELETE /api/users/123 - 404')
    })
  })

  describe('NEGATIVE: HTTP logging with invalid inputs', () => {
    it('should handle empty method', () => {
      logger.logRequest('', '/api/endpoint')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle empty URL', () => {
      logger.logRequest('GET', '')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle negative status code', () => {
      logger.logResponse('GET', '/api/test', -1)

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle very large status code', () => {
      logger.logResponse('GET', '/api/test', 99999)

      // Large status codes (>= 400) are logged as warnings
      expect(consoleSpy.warn).toHaveBeenCalled()
    })

    it('should handle null error object', () => {
      // @ts-expect-error Testing invalid input
      logger.logError('GET', '/api/test', null)

      expect(consoleSpy.error).toHaveBeenCalled()
    })
  })

  describe('POSITIVE: Store action logging', () => {
    it('should log store actions', () => {
      logger.logStoreAction('authStore', 'login', { email: 'test@example.com' })

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('authStore.login')
      expect(callArgs[0]).toContain('[STORE]')
    })

    it('should log store errors', () => {
      const error = new Error('Store action failed')
      logger.logStoreError('userStore', 'fetchUsers', error)

      expect(consoleSpy.error).toHaveBeenCalled()
      const callArgs = consoleSpy.error.mock.calls[0]
      expect(callArgs[0]).toContain('userStore.fetchUsers failed')
      expect(callArgs[0]).toContain('[STORE]')
    })
  })

  describe('NEGATIVE: Store logging with invalid inputs', () => {
    it('should handle empty store name', () => {
      logger.logStoreAction('', 'action')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle empty action name', () => {
      logger.logStoreAction('store', '')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle special characters in store/action names', () => {
      logger.logStoreAction('store@#$', 'action!@#', { data: 'test' })

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle very long store/action names', () => {
      const longName = 'a'.repeat(1000)
      logger.logStoreAction(longName, longName)

      expect(consoleSpy.debug).toHaveBeenCalled()
    })
  })

  describe('POSITIVE: Component lifecycle logging', () => {
    it('should log component mount', () => {
      logger.logComponentMount('LoginForm')

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('Component mounted')
      expect(callArgs[0]).toContain('[LoginForm]')
    })

    it('should log component unmount', () => {
      logger.logComponentUnmount('LoginForm')

      expect(consoleSpy.debug).toHaveBeenCalled()
      const callArgs = consoleSpy.debug.mock.calls[0]
      expect(callArgs[0]).toContain('Component unmounted')
    })
  })

  describe('NEGATIVE: Component logging with invalid inputs', () => {
    it('should handle empty component name', () => {
      logger.logComponentMount('')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle numeric component name', () => {
      // @ts-expect-error Testing invalid input
      logger.logComponentMount(123)

      expect(consoleSpy.debug).toHaveBeenCalled()
    })

    it('should handle null component name', () => {
      // @ts-expect-error Testing invalid input
      logger.logComponentUnmount(null)

      expect(consoleSpy.debug).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Circular reference handling', () => {
    it('should handle circular references in data', () => {
      const circularObj: any = { name: 'test' }
      circularObj.self = circularObj

      // Should not throw error
      expect(() => {
        logger.info('Circular object', 'Test', circularObj)
      }).not.toThrow()

      expect(consoleSpy.info).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Very large data payloads', () => {
    it('should handle very large data objects', () => {
      const largeData = {
        items: Array(10000).fill({ id: 1, name: 'Item', data: 'x'.repeat(100) }),
      }

      expect(() => {
        logger.debug('Large payload', 'Test', largeData)
      }).not.toThrow()

      expect(consoleSpy.debug).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Special characters and unicode', () => {
    it('should handle unicode characters', () => {
      logger.info('Message with unicode: 你好世界 🎉')

      expect(consoleSpy.info).toHaveBeenCalled()
      const callArgs = consoleSpy.info.mock.calls[0]
      expect(callArgs[0]).toContain('你好世界 🎉')
    })

    it('should handle emojis', () => {
      logger.warn('Warning with emojis 🚨⚠️💥')

      expect(consoleSpy.warn).toHaveBeenCalled()
    })

    it('should handle newlines in message', () => {
      logger.info('Line 1\nLine 2\nLine 3')

      expect(consoleSpy.info).toHaveBeenCalled()
    })

    it('should handle tabs in message', () => {
      logger.debug('Column1\tColumn2\tColumn3')

      expect(consoleSpy.debug).toHaveBeenCalled()
    })
  })

  describe('NEGATIVE: Error object variations', () => {
    it('should handle Error without message', () => {
      const error = new Error()
      logger.error('Error without message', 'Test', error)

      expect(consoleSpy.error).toHaveBeenCalled()
    })

    it('should handle Error with very long message', () => {
      const error = new Error('x'.repeat(10000))
      logger.error('Error with long message', 'Test', error)

      expect(consoleSpy.error).toHaveBeenCalled()
    })

    it('should handle custom error objects', () => {
      class CustomError extends Error {
        constructor(
          message: string,
          public code: string,
        ) {
          super(message)
        }
      }

      const error = new CustomError('Custom error', 'ERR_CUSTOM')
      logger.error('Custom error type', 'Test', error)

      expect(consoleSpy.error).toHaveBeenCalled()
    })

    it('should handle non-Error objects as errors', () => {
      const fakeError = { message: 'Not a real Error', code: 500 }
      // @ts-expect-error Testing invalid input
      logger.error('Fake error', 'Test', fakeError)

      expect(consoleSpy.error).toHaveBeenCalled()
    })
  })
})
