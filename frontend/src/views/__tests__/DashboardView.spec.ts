import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../DashboardView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')

describe('DashboardView', () => {
  let authStore: ReturnType<typeof useAuthStore>
  let router: ReturnType<typeof createRouter>

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

    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/documents', component: { template: '<div>Documents</div>' } },
        { path: '/tests', component: { template: '<div>Tests</div>' } },
        { path: '/tests/create', component: { template: '<div>Create Test</div>' } },
      ],
    })

    vi.clearAllMocks()
  })

  describe('Admin User', () => {
    it('displays admin welcome message', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Welcome to TestGen - AI-Powered Test Generation')
    })

    it('shows all three stat cards including Documents', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Total Documents')
      expect(wrapper.text()).toContain('Generated Tests')
      expect(wrapper.text()).toContain('Total Questions')
    })

    it('shows Upload Document and Create Test actions', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Upload Document')
      expect(wrapper.text()).toContain('Add new learning material')
      expect(wrapper.text()).toContain('Create Test')
      expect(wrapper.text()).toContain('Generate new test questions')
    })

    it('links to /documents for Upload Document', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      const uploadLink = wrapper.find('a[href="/documents"]')
      expect(uploadLink.exists()).toBe(true)
    })

    it('links to /tests/create for Create Test', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      const createLink = wrapper.find('a[href="/tests/create"]')
      expect(createLink.exists()).toBe(true)
    })
  })

  describe('Teacher User', () => {
    it('displays teacher welcome message', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Welcome to TestGen - AI-Powered Test Generation')
    })

    it('shows all three stat cards including Documents', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Total Documents')
      expect(wrapper.text()).toContain('Generated Tests')
      expect(wrapper.text()).toContain('Total Questions')
    })

    it('shows Upload Document and Create Test actions', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Upload Document')
      expect(wrapper.text()).toContain('Create Test')
    })
  })

  describe('Student User', () => {
    it('displays student welcome message', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Welcome to TestGen - View and take your assigned tests')
    })

    it('does NOT show Documents card', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Total Documents')
    })

    it('shows Assigned Tests instead of Generated Tests', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Assigned Tests')
      expect(wrapper.text()).not.toContain('Generated Tests')
    })

    it('shows Average Score instead of Total Questions', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Average Score')
      expect(wrapper.text()).toContain('0%')
      expect(wrapper.text()).not.toContain('Total Questions')
    })

    it('does NOT show Upload Document action', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Upload Document')
      expect(wrapper.text()).not.toContain('Add new learning material')
    })

    it('does NOT show Create Test action', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Create Test')
      expect(wrapper.text()).not.toContain('Generate new test questions')
    })

    it('shows View Tests action linking to /tests', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('View Tests')
      expect(wrapper.text()).toContain('See your assigned tests')

      const viewTestsLink = wrapper.find('a[href="/tests"]')
      expect(viewTestsLink.exists()).toBe(true)
    })

    it('shows Practice Mode as coming soon (disabled)', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Practice Mode')
      expect(wrapper.text()).toContain('Coming soon')

      // Practice mode should be a div, not a link
      const practiceDiv = wrapper.find('div.opacity-50.cursor-not-allowed')
      expect(practiceDiv.exists()).toBe(true)
      expect(practiceDiv.text()).toContain('Practice Mode')
    })
  })

  describe('Computed Properties', () => {
    it('isTeacherOrAdmin is true for admin', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Admin should see teacher/admin content
      expect(wrapper.text()).toContain('Total Documents')
      expect(wrapper.text()).toContain('Upload Document')
    })

    it('isTeacherOrAdmin is true for teacher', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Teacher should see teacher/admin content
      expect(wrapper.text()).toContain('Total Documents')
      expect(wrapper.text()).toContain('Upload Document')
    })

    it('isTeacherOrAdmin is false for student', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Student should NOT see teacher/admin content
      expect(wrapper.text()).not.toContain('Total Documents')
      expect(wrapper.text()).not.toContain('Upload Document')
    })
  })
})
