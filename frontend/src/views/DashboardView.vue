<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-text-primary mb-2">Дашборд</h1>
      <p class="text-text-secondary">
        {{ isTeacherOrAdmin
          ? 'Добро пожаловать в TestGen - Генерация тестов на основе ИИ'
          : 'Добро пожаловать в TestGen - Просматривайте и проходите назначенные тесты'
        }}
      </p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="card-cyber text-center py-12 mb-8">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-cyber-blue/20 flex items-center justify-center animate-pulse">
        <svg class="w-8 h-8 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </div>
      <p class="text-text-secondary">Загрузка статистики...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="card-cyber border-red-500/20 bg-red-500/5 p-6 mb-8">
      <div class="flex items-center gap-3 mb-4">
        <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-red-400">Ошибка загрузки статистики</h3>
      </div>
      <p class="text-red-300 mb-4">{{ error }}</p>
      <button @click="loadStats" class="btn-neon">Попробовать снова</button>
    </div>

    <!-- Stats Cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <!-- Documents Card (Teacher/Admin only) -->
      <div v-if="isTeacherOrAdmin" class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">Всего документов</p>
            <p class="text-3xl font-bold text-text-primary mt-2">{{ stats.documents_count }}</p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-neon-orange/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Tests Card -->
      <div class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">
              {{ isTeacherOrAdmin ? 'Созданных тестов' : 'Назначенных тестов' }}
            </p>
            <p class="text-3xl font-bold text-text-primary mt-2">{{ stats.tests_count }}</p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-cyber-blue/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Questions/Score Card -->
      <div class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">
              {{ isTeacherOrAdmin ? 'Всего вопросов' : 'Средний балл' }}
            </p>
            <p class="text-3xl font-bold text-text-primary mt-2">
              {{ isTeacherOrAdmin ? stats.questions_count : '0%' }}
            </p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-cyber-purple/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-cyber-purple" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                :d="isTeacherOrAdmin
                  ? 'M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
                  : 'M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z'"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card-cyber">
      <h2 class="text-xl font-bold text-text-primary mb-4">Быстрые действия</h2>

      <!-- Teacher/Admin Actions -->
      <div v-if="isTeacherOrAdmin" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <router-link
          to="/documents"
          class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 hover:border-neon-orange transition-all duration-300 group"
        >
          <div class="w-10 h-10 rounded-lg bg-neon-orange/20 flex items-center justify-center mr-4 group-hover:bg-neon-orange/30">
            <svg class="w-5 h-5 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary group-hover:text-neon-orange transition-colors">Загрузить документ</p>
            <p class="text-sm text-text-muted">Добавить новый учебный материал</p>
          </div>
        </router-link>

        <router-link
          to="/tests/create"
          class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 hover:border-neon-orange transition-all duration-300 group"
        >
          <div class="w-10 h-10 rounded-lg bg-cyber-blue/20 flex items-center justify-center mr-4 group-hover:bg-cyber-blue/30">
            <svg class="w-5 h-5 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary group-hover:text-neon-orange transition-colors">Создать тест</p>
            <p class="text-sm text-text-muted">Сгенерировать новые вопросы</p>
          </div>
        </router-link>
      </div>

      <!-- Student Actions -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <router-link
          to="/tests"
          class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 hover:border-neon-orange transition-all duration-300 group"
        >
          <div class="w-10 h-10 rounded-lg bg-cyber-blue/20 flex items-center justify-center mr-4 group-hover:bg-cyber-blue/30">
            <svg class="w-5 h-5 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary group-hover:text-neon-orange transition-colors">Просмотр тестов</p>
            <p class="text-sm text-text-muted">Посмотреть назначенные тесты</p>
          </div>
        </router-link>

        <div class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 opacity-50 cursor-not-allowed">
          <div class="w-10 h-10 rounded-lg bg-cyber-purple/20 flex items-center justify-center mr-4">
            <svg class="w-5 h-5 text-cyber-purple" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary">Режим практики</p>
            <p class="text-sm text-text-muted">Скоро</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import { statsService, type DashboardStats } from '@/services/statsService'
import logger from '@/utils/logger'

const authStore = useAuthStore()

const isTeacherOrAdmin = computed(() => {
  const role = authStore.user?.role
  return role === 'teacher' || role === 'admin'
})

const stats = ref<DashboardStats>({
  documents_count: 0,
  tests_count: 0,
  questions_count: 0,
})

const loading = ref(false)
const error = ref<string | null>(null)

async function loadStats() {
  loading.value = true
  error.value = null

  try {
    stats.value = await statsService.getDashboardStats()
    logger.info('Dashboard stats loaded successfully', 'DashboardView', stats.value)
  } catch (err: any) {
    error.value = err.message || 'Не удалось загрузить статистику'
    logger.error('Failed to load dashboard stats', 'DashboardView', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>
