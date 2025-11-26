import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import QuestionEditModal from '../QuestionEditModal.vue'
import type { Question } from '@/features/tests/types/test.types'

// Mock logger
vi.mock('@/utils/logger', () => ({
  default: {
    info: vi.fn(),
    error: vi.fn(),
  },
}))

const mockQuestion: Question = {
  id: 'question-1',
  test_id: 'test-1',
  question_text: 'What is Vue.js?',
  question_type: 'single_choice',
  difficulty: 'medium',
  points: 2.0,
  order_num: 1,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  answers: [
    {
      id: 'answer-1',
      question_id: 'question-1',
      answer_text: 'A JavaScript framework',
      is_correct: true,
      order_num: 0,
      created_at: '2024-01-01T00:00:00Z',
    },
    {
      id: 'answer-2',
      question_id: 'question-1',
      answer_text: 'A CSS library',
      is_correct: false,
      order_num: 1,
      created_at: '2024-01-01T00:00:00Z',
    },
  ],
}

describe('QuestionEditModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render when show is true', () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    expect(wrapper.find('h2').text()).toBe('Edit Question')
    expect(wrapper.find('textarea').element.value).toBe('What is Vue.js?')
  })

  it('should not render when show is false', () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: false,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    expect(wrapper.find('h2').exists()).toBe(false)
  })

  it('should populate form with question data', () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const textarea = wrapper.find('textarea')
    const selects = wrapper.findAll('select')
    const pointsInput = wrapper.find('input[type="number"]')

    expect(textarea.element.value).toBe('What is Vue.js?')
    expect(selects[0].element.value).toBe('single_choice')
    expect(selects[1].element.value).toBe('medium')
    expect(pointsInput.element.value).toBe('2')
  })

  it('should render answers with correct checkbox states', () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    const textInputs = wrapper.findAll('input[type="text"]')

    expect(checkboxes).toHaveLength(2)
    expect(checkboxes[0].element.checked).toBe(true) // First answer is correct
    expect(checkboxes[1].element.checked).toBe(false) // Second answer is incorrect
    expect(textInputs[0].element.value).toBe('A JavaScript framework')
    expect(textInputs[1].element.value).toBe('A CSS library')
  })

  it('should add a new answer when clicking Add Answer button', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const addButton = wrapper.findAll('button').find(btn => btn.text().includes('Add Answer'))
    expect(addButton).toBeDefined()

    await addButton?.trigger('click')

    const textInputs = wrapper.findAll('input[type="text"]')
    expect(textInputs).toHaveLength(3) // Should have 3 answers now
  })

  it('should remove an answer when clicking remove button', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: {
          ...mockQuestion,
          answers: [
            ...mockQuestion.answers,
            {
              id: 'answer-3',
              question_id: 'question-1',
              answer_text: 'A database',
              is_correct: false,
              order_num: 2,
              created_at: '2024-01-01T00:00:00Z',
            },
          ],
        },
        testId: 'test-1',
      },
    })

    let textInputs = wrapper.findAll('input[type="text"]')
    expect(textInputs).toHaveLength(3)

    // Find and click the last remove button
    const removeButtons = wrapper.findAll('button').filter(btn => {
      const svg = btn.find('svg')
      return svg.exists() && btn.element.type === 'button' && !btn.text()
    })

    await removeButtons[removeButtons.length - 1].trigger('click')

    textInputs = wrapper.findAll('input[type="text"]')
    expect(textInputs).toHaveLength(2)
  })

  it('should not allow removing answers when only 2 remain', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const textInputs = wrapper.findAll('input[type="text"]')
    expect(textInputs).toHaveLength(2)

    const removeButtons = wrapper.findAll('button').filter(btn => {
      const svg = btn.find('svg')
      return svg.exists() && btn.element.type === 'button' && !btn.text()
    })

    // Both remove buttons should be disabled
    removeButtons.forEach(btn => {
      expect(btn.element.disabled).toBe(true)
    })
  })

  it('should disable submit when form is invalid (empty question text)', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const textarea = wrapper.find('textarea')
    await textarea.setValue('')

    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Save Changes'))
    expect(submitButton?.element.disabled).toBe(true)
  })

  it('should disable submit when form is invalid (no correct answer)', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    // Uncheck all checkboxes
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    for (const checkbox of checkboxes) {
      await checkbox.setValue(false)
    }

    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Save Changes'))
    expect(submitButton?.element.disabled).toBe(true)
  })

  it('should emit close event when clicking cancel', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const cancelButton = wrapper.findAll('button').find(btn => btn.text() === 'Cancel')
    await cancelButton?.trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit close event when clicking backdrop', async () => {
    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const backdrop = wrapper.find('.fixed.inset-0.bg-black\\/60')
    await backdrop.trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should submit form with valid data', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ ...mockQuestion, question_text: 'Updated question' }),
      } as Response)
    )

    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const textarea = wrapper.find('textarea')
    await textarea.setValue('Updated question')

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/tests/test-1/questions/question-1'),
      expect.objectContaining({
        method: 'PUT',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
    )

    expect(wrapper.emitted('saved')).toBeTruthy()
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should display error message on failed submission', async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        json: () => Promise.resolve({ error: 'Failed to update question' }),
        statusText: 'Internal Server Error',
      } as Response)
    )

    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Failed to update question')
  })

  it('should show loading state during submission', async () => {
    global.fetch = vi.fn(() => new Promise(resolve => setTimeout(resolve, 1000)))

    const wrapper = mount(QuestionEditModal, {
      props: {
        show: true,
        question: mockQuestion,
        testId: 'test-1',
      },
    })

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    // Should show loading text
    const submitButton = wrapper.findAll('button').find(btn => btn.text().includes('Saving'))
    expect(submitButton).toBeDefined()
    expect(submitButton?.element.disabled).toBe(true)
  })
})
