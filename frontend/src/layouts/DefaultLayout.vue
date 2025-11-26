<template>
  <div class="min-h-screen bg-dark-900">
    <!-- Header with cyber gradient -->
    <header class="header-gradient sticky top-0 z-50 backdrop-blur-sm">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo with glow effect -->
          <div class="flex items-center">
            <router-link to="/dashboard" class="text-2xl font-bold text-glow hover:animate-glow">
              TestGen
            </router-link>
            <span class="ml-2 text-xs text-text-muted uppercase tracking-wider">AI Powered</span>
          </div>

          <!-- Navigation -->
          <nav class="hidden md:flex space-x-2">
            <router-link
              to="/dashboard"
              class="text-text-secondary hover:text-neon-orange px-4 py-2 rounded-lg text-sm font-medium
                     transition-all duration-300 hover:bg-dark-700
                     border border-transparent hover:border-dark-500"
              active-class="text-neon-orange bg-dark-700 border-dark-500 shadow-neon-sm"
            >
              Дашборд
            </router-link>
            <router-link
              v-if="isTeacherOrAdmin"
              to="/documents"
              class="text-text-secondary hover:text-neon-orange px-4 py-2 rounded-lg text-sm font-medium
                     transition-all duration-300 hover:bg-dark-700
                     border border-transparent hover:border-dark-500"
              active-class="text-neon-orange bg-dark-700 border-dark-500 shadow-neon-sm"
            >
              Документы
            </router-link>
            <router-link
              to="/tests"
              class="text-text-secondary hover:text-neon-orange px-4 py-2 rounded-lg text-sm font-medium
                     transition-all duration-300 hover:bg-dark-700
                     border border-transparent hover:border-dark-500"
              active-class="text-neon-orange bg-dark-700 border-dark-500 shadow-neon-sm"
            >
              Тесты
            </router-link>
            <router-link
              v-if="isTeacherOrAdmin"
              to="/users"
              class="text-text-secondary hover:text-neon-orange px-4 py-2 rounded-lg text-sm font-medium
                     transition-all duration-300 hover:bg-dark-700
                     border border-transparent hover:border-dark-500"
              active-class="text-neon-orange bg-dark-700 border-dark-500 shadow-neon-sm"
            >
              Пользователи
            </router-link>
          </nav>

          <!-- User menu -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2 px-3 py-1 rounded-lg bg-dark-700 border border-dark-500">
              <div class="w-2 h-2 rounded-full bg-cyber-blue animate-pulse"></div>
              <span class="text-sm text-text-primary font-medium">{{ user?.full_name }}</span>
            </div>
            <button
              @click="handleLogout"
              class="text-sm text-text-secondary hover:text-neon-orange font-medium
                     px-4 py-2 rounded-lg border border-dark-500 hover:border-neon-orange
                     transition-all duration-300 hover:shadow-neon-sm"
            >
              Выход
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main content with grid background pattern -->
    <main class="container mx-auto px-4 sm:px-6 lg:px-8 py-8 relative">
      <!-- Cyber grid background -->
      <div class="absolute inset-0 pointer-events-none opacity-5">
        <div class="absolute inset-0" style="background-image: linear-gradient(#ff6b35 1px, transparent 1px), linear-gradient(90deg, #ff6b35 1px, transparent 1px); background-size: 50px 50px;"></div>
      </div>

      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/authStore'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const isTeacherOrAdmin = computed(() => {
  const role = authStore.user?.role
  return role === 'teacher' || role === 'admin'
})

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
