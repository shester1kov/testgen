import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import TestsView from '../TestsView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { useTestsStore } from '@/features/tests/stores/testsStore'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')
vi.mock('@/services/testService')

const mockRouterPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockRouterPush,
  }),
}))

describe('TestsView', () => {
  let wrapper: VueWrapper
  let authStore: ReturnType<typeof useAuthStore>
  let testsStore: ReturnType<typeof useTestsStore>

  const adminUser: User = {
    id: '1',
    email: 'admin@test.com',
    full_name: 'Admin User',
    role: 'admin',
    created_at: '2024-01-01T00:00:00Z',
  }

  const teacherUser: User = {
    id: '2',
    email: 'teacher@test.com',
    full_name: 'Teacher User',
    role: 'teacher',
    created_at: '2024-01-02T00:00:00Z',
  }

  const studentUser: User = {
    id: '3',
    email: 'student@test.com',
    full_name: 'Student User',
    role: 'student',
    created_at: '2024-01-03T00:00:00Z',
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    authStore = useAuthStore()
    testsStore = useTestsStore()
    vi.clearAllMocks()
    mockRouterPush.mockClear()

    // Reset tests store to initial state
    testsStore.loading = false
    testsStore.error = null
    testsStore.tests = []

    // Mock fetchTests to prevent onMounted lifecycle from failing
    vi.spyOn(testsStore, 'fetchTests').mockResolvedValue({
      data: [],
      total: 0,
      page: 1,
      page_size: 100,
    } as any)
  })

  function mountComponent() {
    return mount(TestsView)
  }

  describe('Admin role', () => {
    // POSITIVE TEST: Admin sees create test button
    it('shows create test button for admin', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const buttons = wrapper.findAll('button')
      const hasCreateButton = buttons.some(btn => btn.text().includes('Generate Test'))
      expect(hasCreateButton).toBe(true)
    })

    // POSITIVE TEST: Admin sees generate test button
    it('shows generate test button for admin', async () => {
      authStore.user = adminUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate Test')
    })

    // POSITIVE TEST: Admin sees correct description
    it('shows management description for admin', async () => {
      authStore.user = adminUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate and manage your test questions')
    })

    // POSITIVE TEST: Admin sees correct empty state
    it('shows create-focused empty state for admin', async () => {
      authStore.user = adminUser
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('No tests yet')
    })
  })

  describe('Teacher role', () => {
    // POSITIVE TEST: Teacher sees create test button
    it('shows create test button for teacher', async () => {
      authStore.user = teacherUser
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const buttons = wrapper.findAll('button')
      const hasGenerateButton = buttons.some(btn => btn.text().includes('Generate Test'))
      expect(hasGenerateButton).toBe(true)
    })

    // POSITIVE TEST: Teacher sees generate test button
    it('shows generate test button for teacher', async () => {
      authStore.user = teacherUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate Test')
    })

    // POSITIVE TEST: Teacher sees correct description
    it('shows management description for teacher', async () => {
      authStore.user = teacherUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate and manage your test questions')
    })
  })

  describe('Student role', () => {
    // NEGATIVE TEST: Student does not see create test button
    it('hides create test button for student', async () => {
      authStore.user = studentUser
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const buttons = wrapper.findAll('button')
      const hasGenerateButton = buttons.some(btn => btn.text().includes('Generate Test'))
      expect(hasGenerateButton).toBe(false)
    })

    // NEGATIVE TEST: Student does not see generate test button
    it('hides generate test button for student', async () => {
      authStore.user = studentUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('Generate Test')
    })

    // POSITIVE TEST: Student sees view-only description
    it('shows view-only description for student', async () => {
      authStore.user = studentUser

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('View your assigned tests')
      expect(wrapper.text()).not.toContain('Generate and manage')
    })

    // POSITIVE TEST: Student sees assigned tests empty state
    it('shows assigned tests empty state for student', async () => {
      authStore.user = studentUser
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('No assigned tests')
    })

    // NEGATIVE TEST: Student does not see create-focused messages
    it('does not show create-focused messages for student', async () => {
      authStore.user = studentUser
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('No tests yet')
      expect(wrapper.text()).not.toContain('Create your first test')
    })
  })

  describe('Role-based UI elements', () => {
    // POSITIVE TEST: isTeacherOrAdmin computed works correctly
    it('correctly calculates isTeacherOrAdmin for different roles', async () => {
      const roles = [
        { user: adminUser, expected: true },
        { user: teacherUser, expected: true },
        { user: studentUser, expected: false },
      ]

      for (const { user, expected } of roles) {
        authStore.user = user
        wrapper = mountComponent()
        await wrapper.vm.$nextTick()

        const vm = wrapper.vm as any
        expect(vm.isTeacherOrAdmin).toBe(expected)
      }
    })
  })

  describe('Loading and Error States', () => {
    // POSITIVE TEST: Shows loading state while fetching tests
    it('shows loading state while fetching tests', async () => {
      authStore.user = adminUser
      testsStore.loading = true

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Loading tests...')
      const loadingIcon = wrapper.find('svg')
      expect(loadingIcon.exists()).toBe(true)
    })

    // POSITIVE TEST: Shows error state when fetch fails
    it('shows error state when fetch fails', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = 'Network error: Failed to fetch'

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Error loading tests')
      expect(wrapper.text()).toContain('Network error: Failed to fetch')
    })

    // POSITIVE TEST: Allows retry on error
    it('allows retry on error with Try Again button', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = 'Failed to fetch tests'
      const fetchTestsSpy = vi.spyOn(testsStore, 'fetchTests').mockResolvedValue({} as any)

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const tryAgainButton = wrapper.findAll('button').find(btn => btn.text().includes('Try Again'))
      expect(tryAgainButton).toBeTruthy()

      await tryAgainButton!.trigger('click')
      await wrapper.vm.$nextTick()

      expect(fetchTestsSpy).toHaveBeenCalled()
    })

    // NEGATIVE TEST: Does not show error state when no error
    it('does not show error state when no error exists', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('Error loading tests')
    })
  })

  describe('Tests Display', () => {
    // POSITIVE TEST: Displays test cards with all information
    it('displays test cards with all information', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '1',
          title: 'Math Test',
          description: 'Basic algebra questions',
          total_questions: 10,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
        {
          id: '2',
          title: 'Science Quiz',
          description: 'Chemistry and physics',
          total_questions: 5,
          status: 'published',
          moodle_synced: true,
          created_at: '2024-01-20T14:30:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Math Test')
      expect(wrapper.text()).toContain('Basic algebra questions')
      expect(wrapper.text()).toContain('10 questions')

      expect(wrapper.text()).toContain('Science Quiz')
      expect(wrapper.text()).toContain('Chemistry and physics')
      expect(wrapper.text()).toContain('5 questions')
    })

    // POSITIVE TEST: Shows test status badges with correct colors
    it('shows test status badges with correct styling', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '1',
          title: 'Draft Test',
          total_questions: 5,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
        {
          id: '2',
          title: 'Published Test',
          total_questions: 10,
          status: 'published',
          moodle_synced: false,
          created_at: '2024-01-20T14:30:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const statusBadges = wrapper.findAll('span[class*="rounded-full"]')
      expect(statusBadges.length).toBeGreaterThanOrEqual(2)
      expect(wrapper.text()).toContain('draft')
      expect(wrapper.text()).toContain('published')
    })

    // POSITIVE TEST: Shows moodle sync indicator when synced
    it('shows moodle sync indicator when test is synced', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '1',
          title: 'Synced Test',
          total_questions: 10,
          status: 'published',
          moodle_synced: true,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Synced')
    })

    // NEGATIVE TEST: Does not show moodle sync indicator when not synced
    it('does not show moodle sync indicator when test is not synced', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '1',
          title: 'Unsynced Test',
          total_questions: 10,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('Synced')
    })

    // POSITIVE TEST: Formats dates correctly
    it('formats dates with relative time', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null

      const now = new Date()
      const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)

      testsStore.tests = [
        {
          id: '1',
          title: 'Recent Test',
          total_questions: 5,
          status: 'draft',
          moodle_synced: false,
          created_at: yesterday.toISOString(),
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Created')
      // Should show relative time like "1d ago"
      expect(wrapper.text()).toMatch(/Created.*ago|Created.*\d{1,2}\/\d{1,2}\/\d{4}/)
    })
  })

  describe('Test Actions', () => {
    // POSITIVE TEST: Navigates to test details on card click
    it('navigates to test details on card click', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '123',
          title: 'Math Test',
          total_questions: 10,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const testCard = wrapper.find('.card-cyber')
      await testCard.trigger('click')

      expect(mockRouterPush).toHaveBeenCalledWith('/tests/123')
    })

    // POSITIVE TEST: Deletes test with confirmation for admin
    it('deletes test with confirmation for admin', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '456',
          title: 'Test to Delete',
          total_questions: 5,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      const deleteTestSpy = vi.spyOn(testsStore, 'deleteTest').mockResolvedValue(undefined)
      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true)

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const deleteButton = wrapper.find('button[title="Delete test"]')
      expect(deleteButton.exists()).toBe(true)

      await deleteButton.trigger('click')
      await wrapper.vm.$nextTick()

      expect(confirmSpy).toHaveBeenCalledWith('Are you sure you want to delete this test?')
      expect(deleteTestSpy).toHaveBeenCalledWith('456')

      confirmSpy.mockRestore()
    })

    // NEGATIVE TEST: Cancels delete when user declines
    it('cancels delete when user declines confirmation', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '789',
          title: 'Test Not Deleted',
          total_questions: 5,
          status: 'draft',
          moodle_synced: false,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      const deleteTestSpy = vi.spyOn(testsStore, 'deleteTest')
      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false)

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const deleteButton = wrapper.find('button[title="Delete test"]')
      await deleteButton.trigger('click')

      expect(confirmSpy).toHaveBeenCalled()
      expect(deleteTestSpy).not.toHaveBeenCalled()

      confirmSpy.mockRestore()
    })

    // NEGATIVE TEST: Student cannot see delete button
    it('student cannot see delete button', async () => {
      authStore.user = studentUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = [
        {
          id: '1',
          title: 'Test for Student',
          total_questions: 5,
          status: 'published',
          moodle_synced: true,
          created_at: '2024-01-15T10:00:00Z',
        },
      ] as any

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const deleteButton = wrapper.find('button[title="Delete test"]')
      expect(deleteButton.exists()).toBe(false)
    })

    // POSITIVE TEST: Navigates to create test page
    it('navigates to create test page when Generate Test button clicked', async () => {
      authStore.user = adminUser
      testsStore.loading = false
      testsStore.error = null
      testsStore.tests = []

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      const generateButton = wrapper.findAll('button').find(btn => btn.text().includes('Generate Test'))
      expect(generateButton).toBeTruthy()

      await generateButton!.trigger('click')

      expect(mockRouterPush).toHaveBeenCalledWith('/tests/create')
    })
  })

  describe('Component Lifecycle', () => {
    // POSITIVE TEST: Fetches tests on mount
    it('fetches tests automatically on component mount', async () => {
      authStore.user = adminUser
      const fetchTestsSpy = vi.spyOn(testsStore, 'fetchTests').mockResolvedValue({} as any)

      wrapper = mountComponent()
      await wrapper.vm.$nextTick()

      expect(fetchTestsSpy).toHaveBeenCalledWith(1, 100)
    })
  })
})
