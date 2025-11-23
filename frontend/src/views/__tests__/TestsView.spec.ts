import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import TestsView from '../TestsView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')

describe('TestsView', () => {
  let wrapper: VueWrapper
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

  describe('Admin role', () => {
    // POSITIVE TEST: Admin sees create test button
    it('shows create test button for admin', async () => {
      authStore.user = adminUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      const createButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Create Test'))
      expect(createButtons.length).toBeGreaterThan(0)
    })

    // POSITIVE TEST: Admin sees generate test button
    it('shows generate test button for admin', async () => {
      authStore.user = adminUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate Test')
    })

    // POSITIVE TEST: Admin sees correct description
    it('shows management description for admin', async () => {
      authStore.user = adminUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate and manage your test questions')
    })

    // POSITIVE TEST: Admin sees correct empty state
    it('shows create-focused empty state for admin', async () => {
      authStore.user = adminUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('No tests yet')
      expect(wrapper.text()).toContain('Create your first test from uploaded documents')
    })
  })

  describe('Teacher role', () => {
    // POSITIVE TEST: Teacher sees create test button
    it('shows create test button for teacher', async () => {
      authStore.user = teacherUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      const createButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Create Test'))
      expect(createButtons.length).toBeGreaterThan(0)
    })

    // POSITIVE TEST: Teacher sees generate test button
    it('shows generate test button for teacher', async () => {
      authStore.user = teacherUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate Test')
    })

    // POSITIVE TEST: Teacher sees correct description
    it('shows management description for teacher', async () => {
      authStore.user = teacherUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('Generate and manage your test questions')
    })
  })

  describe('Student role', () => {
    // NEGATIVE TEST: Student does not see create test button
    it('hides create test button for student', async () => {
      authStore.user = studentUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      const createButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Create Test'))
      expect(createButtons.length).toBe(0)
    })

    // NEGATIVE TEST: Student does not see generate test button
    it('hides generate test button for student', async () => {
      authStore.user = studentUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).not.toContain('Generate Test')
    })

    // POSITIVE TEST: Student sees view-only description
    it('shows view-only description for student', async () => {
      authStore.user = studentUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('View your assigned tests')
      expect(wrapper.text()).not.toContain('Generate and manage')
    })

    // POSITIVE TEST: Student sees assigned tests empty state
    it('shows assigned tests empty state for student', async () => {
      authStore.user = studentUser

      wrapper = mount(TestsView)
      await wrapper.vm.$nextTick()

      expect(wrapper.text()).toContain('No assigned tests')
      expect(wrapper.text()).toContain('You have no tests assigned yet. Please contact your teacher.')
    })

    // NEGATIVE TEST: Student does not see create-focused messages
    it('does not show create-focused messages for student', async () => {
      authStore.user = studentUser

      wrapper = mount(TestsView)
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
        wrapper = mount(TestsView)
        await wrapper.vm.$nextTick()

        const vm = wrapper.vm as any
        expect(vm.isTeacherOrAdmin).toBe(expected)
      }
    })
  })
})
