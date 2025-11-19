import { beforeAll, afterEach, afterAll, vi } from 'vitest'
import { config } from '@vue/test-utils'

// Mock window.matchMedia
beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  })
})

// Cleanup after each test
afterEach(() => {
  vi.clearAllMocks()
})

// Global teardown
afterAll(() => {
  vi.clearAllTimers()
})

// Configure Vue Test Utils
config.global.mocks = {
  $t: (key: string) => key, // Mock i18n if needed
}
