import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import UsersView from '../UsersView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import UserList from '@/features/users/components/UserList.vue'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')
vi.mock('@/services/userService')

describe('UsersView', () => {
  let authStore: ReturnType<typeof useAuthStore>

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
    vi.clearAllMocks()
  })

  describe('Admin User', () => {
    it('displays admin-specific description', () => {
      authStore.user = adminUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).toContain('Users Management')
      expect(wrapper.text()).toContain('Manage user roles and permissions')
    })

    it('renders UserList component', () => {
      authStore.user = adminUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.findComponent(UserList).exists()).toBe(true)
    })
  })

  describe('Teacher User', () => {
    it('displays teacher-specific description', () => {
      authStore.user = teacherUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).toContain('Users Management')
      expect(wrapper.text()).toContain('View students and assign tests to them')
    })

    it('does NOT show admin description', () => {
      authStore.user = teacherUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).not.toContain('Manage user roles and permissions')
    })

    it('renders UserList component', () => {
      authStore.user = teacherUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.findComponent(UserList).exists()).toBe(true)
    })
  })

  describe('Computed Property isAdmin', () => {
    it('returns true for admin user', () => {
      authStore.user = adminUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).toContain('Manage user roles and permissions')
    })

    it('returns false for teacher user', () => {
      authStore.user = teacherUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).toContain('View students and assign tests to them')
    })

    it('returns false for student user', () => {
      authStore.user = studentUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      expect(wrapper.text()).toContain('View students and assign tests to them')
    })

    it('returns false when user is null', () => {
      authStore.user = null

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      // Should default to non-admin description
      expect(wrapper.text()).toContain('View students and assign tests to them')
    })
  })

  describe('Layout', () => {
    it('has proper heading structure', () => {
      authStore.user = adminUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      const heading = wrapper.find('h1')
      expect(heading.exists()).toBe(true)
      expect(heading.text()).toBe('Users Management')
      expect(heading.classes()).toContain('text-3xl')
      expect(heading.classes()).toContain('font-bold')
    })

    it('has description paragraph', () => {
      authStore.user = adminUser

      const wrapper = mount(UsersView, {
        global: {
          stubs: {
            UserList: true,
          },
        },
      })

      const description = wrapper.find('p.text-text-secondary')
      expect(description.exists()).toBe(true)
      expect(description.text()).toBe('Manage user roles and permissions')
    })
  })
})
