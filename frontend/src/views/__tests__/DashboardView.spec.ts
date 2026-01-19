import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../DashboardView.vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'
import { UserRole } from '@/features/auth/types/auth.types'
import * as statsService from '@/services/statsService'
import type { DashboardStats } from '@/services/statsService'

vi.mock('@/services/authService')
vi.mock('@/services/statsService', () => ({
  statsService: {
    getDashboardStats: vi.fn(),
  },
}))

describe('DashboardView', () => {
  let authStore: ReturnType<typeof useAuthStore>
  let router: ReturnType<typeof createRouter>

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

    // Default mock for stats service
    const defaultStats: DashboardStats = {
      documents_count: 0,
      tests_count: 0,
      questions_count: 0,
    }
    vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(defaultStats)
  })

  describe('Admin User', () => {
    it('displays admin welcome message', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Добро пожаловать в TestGen - Генерация тестов на основе ИИ')
    })

    it('shows all three stat cards including Documents', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Всего документов')
      expect(wrapper.text()).toContain('Созданных тестов')
      expect(wrapper.text()).toContain('Всего вопросов')
    })

    it('shows Upload Document and Create Test actions', () => {
      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Загрузить документ')
      expect(wrapper.text()).toContain('Добавить новый учебный материал')
      expect(wrapper.text()).toContain('Создать тест')
      expect(wrapper.text()).toContain('Сгенерировать новые вопросы')
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

      expect(wrapper.text()).toContain('Добро пожаловать в TestGen - Генерация тестов на основе ИИ')
    })

    it('shows all three stat cards including Documents', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Всего документов')
      expect(wrapper.text()).toContain('Созданных тестов')
      expect(wrapper.text()).toContain('Всего вопросов')
    })

    it('shows Upload Document and Create Test actions', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Загрузить документ')
      expect(wrapper.text()).toContain('Создать тест')
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

      expect(wrapper.text()).toContain('Добро пожаловать в TestGen - Просматривайте и проходите назначенные тесты')
    })

    it('does NOT show Documents card', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Всего документов')
    })

    it('shows Assigned Tests instead of Generated Tests', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Назначенных тестов')
      expect(wrapper.text()).not.toContain('Созданных тестов')
    })

    it('shows Average Score instead of Total Questions', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Средний балл')
      expect(wrapper.text()).toContain('0%')
      expect(wrapper.text()).not.toContain('Всего вопросов')
    })

    it('does NOT show Upload Document action', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Загрузить документ')
      expect(wrapper.text()).not.toContain('Добавить новый учебный материал')
    })

    it('does NOT show Create Test action', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).not.toContain('Создать тест')
      expect(wrapper.text()).not.toContain('Сгенерировать новые вопросы')
    })

    it('shows View Tests action linking to /tests', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      expect(wrapper.text()).toContain('Просмотр тестов')
      expect(wrapper.text()).toContain('Посмотреть назначенные тесты')

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

      expect(wrapper.text()).toContain('Режим практики')
      expect(wrapper.text()).toContain('Скоро')

      // Practice mode should be a div, not a link
      const practiceDiv = wrapper.find('div.opacity-50.cursor-not-allowed')
      expect(practiceDiv.exists()).toBe(true)
      expect(practiceDiv.text()).toContain('Режим практики')
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
      expect(wrapper.text()).toContain('Всего документов')
      expect(wrapper.text()).toContain('Загрузить документ')
    })

    it('isTeacherOrAdmin is true for teacher', () => {
      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Teacher should see teacher/admin content
      expect(wrapper.text()).toContain('Всего документов')
      expect(wrapper.text()).toContain('Загрузить документ')
    })

    it('isTeacherOrAdmin is false for student', () => {
      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Student should NOT see teacher/admin content
      expect(wrapper.text()).not.toContain('Всего документов')
      expect(wrapper.text()).not.toContain('Загрузить документ')
    })
  })

  describe('Statistics Loading', () => {
    it('fetches and displays real statistics on mount', async () => {
      const mockStats: DashboardStats = {
        documents_count: 12,
        tests_count: 7,
        questions_count: 84,
      }
      vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(mockStats)

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      expect(statsService.statsService.getDashboardStats).toHaveBeenCalledTimes(1)
      expect(wrapper.text()).toContain('12')
      expect(wrapper.text()).toContain('7')
      expect(wrapper.text()).toContain('84')
    })

    it('displays loading state while fetching statistics', async () => {
      const mockStats: DashboardStats = {
        documents_count: 5,
        tests_count: 10,
        questions_count: 50,
      }

      // Create a promise we can control
      let resolveStats: (value: DashboardStats) => void
      const statsPromise = new Promise<DashboardStats>((resolve) => {
        resolveStats = resolve
      })
      vi.mocked(statsService.statsService.getDashboardStats).mockReturnValue(statsPromise)

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })

      // Wait a tick to let the component mount
      await wrapper.vm.$nextTick()

      // Should show loading state
      expect(wrapper.text()).toContain('Загрузка статистики...')

      // Resolve the promise
      resolveStats!(mockStats)
      await flushPromises()

      // Should no longer show loading
      expect(wrapper.text()).not.toContain('Загрузка статистики...')
      expect(wrapper.text()).toContain('5')
      expect(wrapper.text()).toContain('10')
      expect(wrapper.text()).toContain('50')
    })

    it('displays zero values when no data exists', async () => {
      const mockStats: DashboardStats = {
        documents_count: 0,
        tests_count: 0,
        questions_count: 0,
      }
      vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(mockStats)

      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      // Student should see "Assigned Tests" card with 0
      expect(wrapper.text()).toContain('Назначенных тестов')
      expect(wrapper.text()).toContain('0')
    })

    it('displays error message when stats loading fails', async () => {
      const errorMessage = 'Failed to load statistics'
      vi.mocked(statsService.statsService.getDashboardStats).mockRejectedValue(
        new Error(errorMessage)
      )

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      expect(wrapper.text()).toContain(errorMessage)
    })

    it('displays "Try Again" button when error occurs', async () => {
      vi.mocked(statsService.statsService.getDashboardStats).mockRejectedValue(
        new Error('Network error')
      )

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      const tryAgainButton = wrapper.find('button')
      expect(tryAgainButton.exists()).toBe(true)
      expect(tryAgainButton.text()).toContain('Попробовать снова')
    })

    it('retries loading stats when "Try Again" button is clicked', async () => {
      const mockStats: DashboardStats = {
        documents_count: 5,
        tests_count: 10,
        questions_count: 50,
      }

      // First call fails
      vi.mocked(statsService.statsService.getDashboardStats)
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce(mockStats)

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      // Should show error (checking for partial message)
      expect(wrapper.text()).toContain('Ошибка загрузки статистики')
      expect(statsService.statsService.getDashboardStats).toHaveBeenCalledTimes(1)

      // Click "Try Again"
      const tryAgainButton = wrapper.find('button')
      await tryAgainButton.trigger('click')
      await flushPromises()

      // Should have retried and succeeded
      expect(statsService.statsService.getDashboardStats).toHaveBeenCalledTimes(2)
      expect(wrapper.text()).not.toContain('Ошибка загрузки статистики')
      expect(wrapper.text()).toContain('5')
      expect(wrapper.text()).toContain('10')
      expect(wrapper.text()).toContain('50')
    })

    it('displays different stats for admin (all users data)', async () => {
      const mockStats: DashboardStats = {
        documents_count: 100,
        tests_count: 50,
        questions_count: 500,
      }

      // Setup mock before mounting
      vi.mocked(statsService.statsService.getDashboardStats).mockClear()
      vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(mockStats)

      authStore.user = adminUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      // Admin should see all stats from all users
      expect(wrapper.text()).toContain('100')
      expect(wrapper.text()).toContain('50')
      expect(wrapper.text()).toContain('500')
    })

    it('displays teacher stats (only their own data)', async () => {
      const mockStats: DashboardStats = {
        documents_count: 8,
        tests_count: 5,
        questions_count: 40,
      }
      vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(mockStats)

      authStore.user = teacherUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      // Teacher should see only their own stats
      expect(wrapper.text()).toContain('8')
      expect(wrapper.text()).toContain('5')
      expect(wrapper.text()).toContain('40')
    })

    it('does not display question count for students', async () => {
      const mockStats: DashboardStats = {
        documents_count: 0,
        tests_count: 3,
        questions_count: 0,
      }
      vi.mocked(statsService.statsService.getDashboardStats).mockResolvedValue(mockStats)

      authStore.user = studentUser

      const wrapper = mount(DashboardView, {
        global: {
          plugins: [router],
        },
      })
      await flushPromises()

      // Student sees "Average Score" (0%) instead of question count
      expect(wrapper.text()).toContain('Средний балл')
      expect(wrapper.text()).toContain('0%')
      expect(wrapper.text()).not.toContain('Всего вопросов')
    })
  })
})
