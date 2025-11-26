<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">{{ testsTitle }}</h1>
        <p class="text-text-secondary">
          {{ isTeacherOrAdmin ? 'Генерируйте и управляйте тестовыми вопросами' : 'Просматривайте назначенные тесты' }}
        </p>
      </div>
      <button v-if="isTeacherOrAdmin" @click="navigateToCreate" class="btn-neon">
        <svg class="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Сгенерировать тест
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="testsStore.loading" class="card-cyber text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-cyber-blue/20 flex items-center justify-center animate-pulse">
        <svg class="w-8 h-8 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </div>
      <p class="text-text-secondary">Загрузка тестов...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="testsStore.error" class="card-cyber border-red-500/20 bg-red-500/5 p-6">
      <div class="flex items-center gap-3 mb-4">
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-red-400">Ошибка загрузки тестов</h3>
      </div>
      <p class="text-red-300 mb-4">{{ testsStore.error }}</p>
      <button @click="loadTests" class="btn-neon">
        Попробовать снова
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="testsStore.tests.length === 0" class="card-cyber text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-cyber-blue/20 flex items-center justify-center">
        <svg class="w-8 h-8 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-text-primary mb-2">
        {{ isTeacherOrAdmin ? 'Тестов пока нет' : 'Нет назначенных тестов' }}
      </h3>
      <p class="text-text-muted">
        {{ isTeacherOrAdmin
          ? 'Создайте свой первый тест из загруженных документов, нажав кнопку "Сгенерировать тест" выше'
          : 'У вас пока нет назначенных тестов. Пожалуйста, свяжитесь с преподавателем.'
        }}
      </p>
    </div>

    <!-- Tests List -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="test in testsStore.tests"
        :key="test.id"
        class="card-cyber hover:border-neon-orange/50 transition-all cursor-pointer group"
        @click="viewTest(test.id)"
      >
        <div class="flex justify-between items-start mb-4">
          <h3 class="text-xl font-semibold text-text-primary group-hover:text-neon-orange transition-colors">
            {{ test.title }}
          </h3>
          <span
            :class="getStatusClass(test.status)"
            class="px-3 py-1 rounded-full text-xs font-medium"
          >
            {{ test.status }}
          </span>
        </div>

        <p v-if="test.description" class="text-text-muted text-sm mb-4 line-clamp-2">
          {{ test.description }}
        </p>

        <div v-if="test.user_name || test.user_email" class="flex items-center gap-2 mb-3">
          <span class="px-2 py-0.5 bg-cyber-blue/10 border border-cyber-blue/30 rounded text-xs text-cyber-blue">
            <span v-if="test.user_name">{{ test.user_name }}</span>
            <span v-if="test.user_email" class="text-text-muted ml-1">({{ test.user_email }})</span>
          </span>
        </div>

        <div class="flex items-center justify-between text-sm text-text-secondary">
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ test.total_questions }} questions</span>
            </div>

            <div v-if="test.moodle_synced" class="flex items-center gap-1 text-green-500">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <span>Synced</span>
            </div>
          </div>

          <button
            v-if="isTeacherOrAdmin"
            @click.stop="handleDelete(test.id)"
            class="p-2 hover:bg-red-500/20 rounded-lg transition-colors"
            title="Delete test"
          >
            <svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>

        <div class="mt-4 pt-4 border-t border-cyber-blue/20 text-xs text-text-muted">
          Created {{ formatDate(test.created_at) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { useTestsStore } from '@/features/tests/stores/testsStore'
import { TestStatus } from '@/features/tests/types/test.types'
import logger from '@/utils/logger'

const authStore = useAuthStore()
const testsStore = useTestsStore()
const router = useRouter()

const isTeacherOrAdmin = computed(() => {
  const role = authStore.user?.role
  return role === 'teacher' || role === 'admin'
})

const isAdmin = computed(() => authStore.user?.role === 'admin')

const testsTitle = computed(() => isAdmin.value ? 'Все тесты' : 'Тесты')

function navigateToCreate() {
  router.push('/tests/create')
}

function viewTest(id: string) {
  router.push(`/tests/${id}`)
}

async function loadTests() {
  try {
    await testsStore.fetchTests(1, 100)
  } catch (err: any) {
    logger.error('Failed to load tests', 'TestsView', err)
  }
}

async function handleDelete(id: string) {
  if (!confirm('Are you sure you want to delete this test?')) {
    return
  }

  try {
    await testsStore.deleteTest(id)
    logger.info('Test deleted successfully', 'TestsView', { testId: id })
  } catch (err: any) {
    logger.error('Failed to delete test', 'TestsView', err)
    alert('Failed to delete test: ' + err.message)
  }
}

function getStatusClass(status: string): string {
  switch (status) {
    case TestStatus.DRAFT:
      return 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
    case TestStatus.PUBLISHED:
      return 'bg-green-500/20 text-green-400 border border-green-500/30'
    case TestStatus.ARCHIVED:
      return 'bg-gray-500/20 text-gray-400 border border-gray-500/30'
    default:
      return 'bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/30'
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays < 7) return `${diffDays}d ago`

  return date.toLocaleDateString()
}

onMounted(() => {
  loadTests()
})
</script>
