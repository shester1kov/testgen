import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import router from '../index'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/authService')

describe('Router Guards', () => {
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

  beforeEach(async () => {
    setActivePinia(createPinia())
    authStore = useAuthStore()
    vi.clearAllMocks()
    // Reset router to initial route
    await router.push('/')
    await router.isReady()
  })

  describe('Authentication', () => {
    // NEGATIVE TEST: Redirects unauthenticated users to login
    it('redirects to login when accessing protected route without auth', async () => {
      authStore.user = null

      await router.push('/dashboard')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Login')
    })

    // POSITIVE TEST: Allows authenticated users to access protected routes
    it('allows authenticated users to access protected routes', async () => {
      authStore.user = studentUser

      await router.push('/dashboard')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })

    // POSITIVE TEST: Redirects authenticated users away from login
    it('redirects authenticated users away from login page', async () => {
      authStore.user = studentUser

      await router.push('/login')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })
  })

  describe('Documents access (teacher/admin only)', () => {
    // POSITIVE TEST: Admin can access documents
    it('allows admin to access documents', async () => {
      authStore.user = adminUser

      await router.push('/documents')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Documents')
    })

    // POSITIVE TEST: Teacher can access documents
    it('allows teacher to access documents', async () => {
      authStore.user = teacherUser

      await router.push('/documents')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Documents')
    })

    // NEGATIVE TEST: Student cannot access documents
    it('redirects student away from documents', async () => {
      authStore.user = studentUser

      await router.push('/documents')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })

    // NEGATIVE TEST: Student cannot access document details
    it('redirects student away from document details', async () => {
      authStore.user = studentUser

      await router.push('/documents/123')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })
  })

  describe('Test creation access (teacher/admin only)', () => {
    // POSITIVE TEST: Admin can create tests
    it('allows admin to create tests', async () => {
      authStore.user = adminUser

      await router.push('/tests/create')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('CreateTest')
    })

    // POSITIVE TEST: Teacher can create tests
    it('allows teacher to create tests', async () => {
      authStore.user = teacherUser

      await router.push('/tests/create')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('CreateTest')
    })

    // NEGATIVE TEST: Student cannot create tests
    it('redirects student away from test creation', async () => {
      authStore.user = studentUser

      await router.push('/tests/create')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })

    // POSITIVE TEST: Admin can edit tests
    it('allows admin to edit tests', async () => {
      authStore.user = adminUser

      await router.push('/tests/123/edit')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('EditTest')
    })

    // POSITIVE TEST: Teacher can edit tests
    it('allows teacher to edit tests', async () => {
      authStore.user = teacherUser

      await router.push('/tests/123/edit')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('EditTest')
    })

    // NEGATIVE TEST: Student cannot edit tests
    it('redirects student away from test editing', async () => {
      authStore.user = studentUser

      await router.push('/tests/123/edit')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })
  })

  describe('Tests view access (all roles)', () => {
    // POSITIVE TEST: Admin can view tests
    it('allows admin to view tests', async () => {
      authStore.user = adminUser

      await router.push('/tests')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Tests')
    })

    // POSITIVE TEST: Teacher can view tests
    it('allows teacher to view tests', async () => {
      authStore.user = teacherUser

      await router.push('/tests')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Tests')
    })

    // POSITIVE TEST: Student can view tests
    it('allows student to view tests', async () => {
      authStore.user = studentUser

      await router.push('/tests')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Tests')
    })

    // POSITIVE TEST: All roles can view test details
    it('allows all roles to view test details', async () => {
      const users = [adminUser, teacherUser, studentUser]

      for (const user of users) {
        authStore.user = user

        await router.push('/tests/123')
        await router.isReady()

        expect(router.currentRoute.value.name).toBe('TestDetails')
      }
    })
  })

  describe('Users access (teacher/admin only)', () => {
    // POSITIVE TEST: Admin can access users
    it('allows admin to access users', async () => {
      authStore.user = adminUser

      await router.push('/users')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Users')
    })

    // POSITIVE TEST: Teacher can access users
    it('allows teacher to access users', async () => {
      authStore.user = teacherUser

      await router.push('/users')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Users')
    })

    // NEGATIVE TEST: Student cannot access users
    it('redirects student away from users', async () => {
      authStore.user = studentUser

      await router.push('/users')
      await router.isReady()

      expect(router.currentRoute.value.name).toBe('Dashboard')
    })
  })

  describe('Dashboard access (all authenticated)', () => {
    // POSITIVE TEST: All roles can access dashboard
    it('allows all authenticated users to access dashboard', async () => {
      const users = [adminUser, teacherUser, studentUser]

      for (const user of users) {
        authStore.user = user

        await router.push('/dashboard')
        await router.isReady()

        expect(router.currentRoute.value.name).toBe('Dashboard')
      }
    })
  })
})
