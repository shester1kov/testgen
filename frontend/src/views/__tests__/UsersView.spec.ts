import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import UsersView from '../UsersView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import UserList from '@/features/users/components/UserList.vue'
import type { User } from '@/features/auth/types/auth.types'
import { UserRole } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')
vi.mock('@/services/userService')

describe('UsersView', () => {
  let authStore: ReturnType<typeof useAuthStore>

  const adminUser: User = {
    id: '1',
    email: 'admin@test.com',
    full_name: 'Admin User',
    role: UserRole.ADMIN,
  }

  const teacherUser: User = {
    id: '2',
    email: 'teacher@test.com',
    full_name: 'Teacher User',
    role: UserRole.TEACHER,
  }

  const studentUser: User = {
    id: '3',
    email: 'student@test.com',
    full_name: 'Student User',
    role: UserRole.STUDENT,
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

      expect(wrapper.text()).toContain('Управление пользователями')
      expect(wrapper.text()).toContain('Управляйте ролями и правами пользователей')
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

      expect(wrapper.text()).toContain('Управление пользователями')
      expect(wrapper.text()).toContain('Просматривайте студентов и назначайте им тесты')
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

      expect(wrapper.text()).not.toContain('Управление ролями пользователей и правами доступа')
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

      expect(wrapper.text()).toContain('Управляйте ролями и правами пользователей')
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

      expect(wrapper.text()).toContain('Просматривайте студентов и назначайте им тесты')
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

      expect(wrapper.text()).toContain('Просматривайте студентов и назначайте им тесты')
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
      expect(wrapper.text()).toContain('Просматривайте студентов и назначайте им тесты')
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
      expect(heading.text()).toBe('Управление пользователями')
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
      expect(description.text()).toBe('Управляйте ролями и правами пользователей')
    })
  })
})
