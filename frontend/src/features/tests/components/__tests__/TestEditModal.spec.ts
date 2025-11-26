import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import TestEditModal from '../TestEditModal.vue'
import type { Test } from '@/features/tests/types/test.types'
import { TestStatus } from '@/features/tests/types/test.types'

// Mock logger
vi.mock('@/utils/logger', () => ({
  default: {
    info: vi.fn(),
    error: vi.fn(),
  },
}))

const mockTest: Test = {
  id: 'test-1',
  user_id: 'user-1',
  title: 'JavaScript Basics Test',
  description: 'A test covering JavaScript fundamentals',
  total_questions: 5,
  status: TestStatus.DRAFT,
  moodle_synced: false,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  questions: [],
}

describe('TestEditModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render when show is true', () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    expect(wrapper.find('h2').text()).toBe('Edit Test')
    expect(wrapper.find('input[type="text"]').element.value).toBe('JavaScript Basics Test')
  })

  it('should not render when show is false', () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: false,
        test: mockTest,
      },
    })

    expect(wrapper.find('h2').exists()).toBe(false)
  })

  it('should populate form with test data', () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    const descriptionTextarea = wrapper.find('textarea')

    expect(titleInput.element.value).toBe('JavaScript Basics Test')
    expect(descriptionTextarea.element.value).toBe('A test covering JavaScript fundamentals')
  })

  it('should handle test without description', () => {
    const testWithoutDesc = { ...mockTest, description: undefined }
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: testWithoutDesc,
      },
    })

    const descriptionTextarea = wrapper.find('textarea')
    expect(descriptionTextarea.element.value).toBe('')
  })

  it('should update form when test prop changes', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const newTest: Test = {
      ...mockTest,
      id: 'test-2',
      title: 'Updated Test Title',
      description: 'Updated description',
    }

    await wrapper.setProps({ test: newTest })

    const titleInput = wrapper.find('input[type="text"]')
    const descriptionTextarea = wrapper.find('textarea')

    expect(titleInput.element.value).toBe('Updated Test Title')
    expect(descriptionTextarea.element.value).toBe('Updated description')
  })

  it('should disable submit when title is too short', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    await titleInput.setValue('ab') // Less than 3 characters

    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Save Changes'))
    expect(submitButton?.element.disabled).toBe(true)
  })

  it('should disable submit when title is empty', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    await titleInput.setValue('')

    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Save Changes'))
    expect(submitButton?.element.disabled).toBe(true)
  })

  it('should enable submit when title is valid', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    await titleInput.setValue('Valid Test Title')

    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Save Changes'))
    expect(submitButton?.element.disabled).toBe(false)
  })

  it('should emit close event when clicking cancel', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const cancelButton = wrapper.findAll('button').find(btn => btn.text() === 'Cancel')
    await cancelButton?.trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit close event when clicking backdrop', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const backdrop = wrapper.find('.fixed.inset-0.bg-black\\/60')
    await backdrop.trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit close event when clicking X button', async () => {
    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const closeButton = wrapper.findAll('button').find(btn => {
      const svg = btn.find('svg')
      return svg.exists() && svg.find('path[d*="M6 18L18 6M6 6l12 12"]').exists()
    })

    await closeButton?.trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should submit form with valid data', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () =>
          Promise.resolve({
            ...mockTest,
            title: 'Updated Test Title',
            description: 'Updated description',
          }),
      } as Response)
    )

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    const descriptionTextarea = wrapper.find('textarea')

    await titleInput.setValue('Updated Test Title')
    await descriptionTextarea.setValue('Updated description')

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/tests/test-1'),
      expect.objectContaining({
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: 'Updated Test Title',
          description: 'Updated description',
        }),
      })
    )

    expect(wrapper.emitted('saved')).toBeTruthy()
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should display error message on failed submission', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        json: () => Promise.resolve({ error: 'Failed to update test' }),
        statusText: 'Internal Server Error',
      } as Response)
    )

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Failed to update test')
  })

  it('should show loading state during submission', async () => {
    global.fetch = vi.fn(() => new Promise(resolve => setTimeout(resolve, 1000)))

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Should show loading text
    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Saving'))
    expect(submitButton).toBeDefined()
    expect(submitButton?.element.disabled).toBe(true)
  })

  it('should not allow closing modal during submission', async () => {
    global.fetch = vi.fn(() => new Promise(resolve => setTimeout(resolve, 1000)))

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Try to click cancel during loading
    const cancelButton = wrapper.findAll('button').find(btn => btn.text() === 'Cancel')
    expect(cancelButton?.element.disabled).toBe(true)
  })

  it('should handle network errors gracefully', async () => {
    global.fetch = vi.fn(() => Promise.reject(new Error('Network error')))

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Network error')
  })

  it('should trim whitespace from title', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ ...mockTest }),
      } as Response)
    )

    const wrapper = mount(TestEditModal, {
      props: {
        show: true,
        test: mockTest,
      },
    })

    const titleInput = wrapper.find('input[type="text"]')
    await titleInput.setValue('   Trimmed Title   ')

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    // Form is valid if trimmed length >= 3
    expect(wrapper.emitted('saved')).toBeTruthy()
  })
})
