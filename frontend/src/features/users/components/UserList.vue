<template>
  <div class="card-cyber">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-text-primary">Управление пользователями</h2>
        <p class="text-sm text-text-muted mt-1">Всего пользователей: {{ total }}</p>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading && users.length === 0" class="flex flex-col items-center justify-center py-12">
      <div class="w-12 h-12 border-4 border-neon-orange/30 border-t-neon-orange rounded-full animate-spin mb-4"></div>
      <p class="text-text-muted">Загрузка пользователей...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="users.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-dark-600 flex items-center justify-center">
        <svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-text-primary mb-2">Пользователи не найдены</h3>
      <p class="text-text-muted">В системе нет пользователей</p>
    </div>

    <!-- Users table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-dark-500">
            <th class="text-left py-3 px-4 text-sm font-semibold text-text-primary">Email</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-text-primary">Полное имя</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-text-primary">Роль</th>
            <th v-if="canChangeRoles" class="text-right py-3 px-4 text-sm font-semibold text-text-primary">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="border-b border-dark-600 hover:bg-dark-600/50 transition-colors">
            <td class="py-3 px-4 text-sm text-text-secondary">{{ user.email }}</td>
            <td class="py-3 px-4 text-sm text-text-secondary">{{ user.full_name }}</td>
            <td class="py-3 px-4">
              <span class="px-2.5 py-1 rounded-full text-xs font-medium" :class="getRoleBadgeClass(user.role)">
                {{ formatRole(user.role) }}
              </span>
            </td>
            <td v-if="canChangeRoles" class="py-3 px-4 text-right">
              <button
                type="button"
                class="px-3 py-1.5 bg-neon-orange/10 border border-neon-orange/30 text-neon-orange rounded-lg text-sm font-medium hover:bg-neon-orange/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="isProcessing"
                @click="openRoleModal(user)"
              >
                Изменить роль
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between mt-6 pt-4 border-t border-dark-500">
      <button
        type="button"
        class="px-4 py-2 border border-dark-500 rounded-lg text-sm font-medium text-text-secondary hover:bg-dark-600 hover:border-neon-orange transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Назад
      </button>

      <div class="text-sm text-text-muted">
        Страница <span class="text-neon-orange font-medium">{{ currentPage }}</span> из {{ totalPages }}
      </div>

      <button
        type="button"
        class="px-4 py-2 border border-dark-500 rounded-lg text-sm font-medium text-text-secondary hover:bg-dark-600 hover:border-neon-orange transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        :disabled="currentPage === totalPages"
        @click="handlePageChange(currentPage + 1)"
      >
        Вперёд
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="mt-4 flex items-center justify-between gap-4 p-4 bg-cyber-pink/10 border border-cyber-pink/30 rounded-lg">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-cyber-pink flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-sm text-cyber-pink flex-1">{{ error }}</p>
      </div>
      <button
        type="button"
        class="px-3 py-1.5 bg-cyber-pink/20 text-cyber-pink rounded-lg text-sm font-medium hover:bg-cyber-pink/30 transition-colors"
        @click="handleRetry"
      >
        Повторить
      </button>
    </div>

    <!-- Role Change Modal -->
    <Teleport to="body">
      <div v-if="showRoleModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click="closeRoleModal">
        <div class="card-cyber max-w-md w-full" @click.stop>
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-text-primary">Изменение роли пользователя</h3>
            <button type="button" class="text-text-muted hover:text-neon-orange transition-colors" @click="closeRoleModal">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div v-if="selectedUser" class="mb-6">
            <p class="text-sm text-text-muted mb-1">Пользователь</p>
            <p class="text-base text-text-primary font-medium">{{ selectedUser.email }}</p>
            <p class="text-sm text-text-secondary">{{ selectedUser.full_name }}</p>
          </div>

          <div class="mb-6">
            <label class="block text-sm font-medium text-text-primary mb-2">Новая роль</label>
            <div class="space-y-2">
              <label
                v-for="role in availableRoles"
                :key="role.value"
                class="flex items-center p-3 border rounded-lg cursor-pointer transition-all"
                :class="selectedRole === role.value
                  ? 'border-neon-orange bg-neon-orange/10'
                  : 'border-dark-500 hover:border-dark-600'"
              >
                <input
                  type="radio"
                  :value="role.value"
                  v-model="selectedRole"
                  class="w-4 h-4 text-neon-orange focus:ring-neon-orange focus:ring-offset-0 bg-dark-600 border-dark-500"
                />
                <div class="ml-3">
                  <p class="text-sm font-medium text-text-primary">{{ role.label }}</p>
                  <p class="text-xs text-text-muted">{{ role.description }}</p>
                </div>
              </label>
            </div>
          </div>

          <div class="flex gap-3">
            <button
              type="button"
              class="flex-1 px-4 py-2 border border-dark-500 text-text-secondary rounded-lg text-sm font-medium hover:bg-dark-600 transition-all"
              @click="closeRoleModal"
            >
              Отмена
            </button>
            <button
              type="button"
              class="flex-1 px-4 py-2 bg-neon-orange text-dark-800 rounded-lg text-sm font-medium hover:bg-neon-orange-light transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="!selectedRole || isProcessing"
              @click="handleRoleChange"
            >
              {{ isProcessing ? 'Обновление...' : 'Изменить роль' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useUsersStore } from '../stores/usersStore'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'

const usersStore = useUsersStore()
const authStore = useAuthStore()

const { users, total, currentPage, totalPages, loading: isLoading, error } = storeToRefs(usersStore)
const { user: currentUser } = storeToRefs(authStore)

const isProcessing = ref(false)
const showRoleModal = ref(false)
const selectedUser = ref<User | null>(null)
const selectedRole = ref<string>('')

const availableRoles = [
  { value: 'student', label: 'Студент', description: 'Может просматривать и проходить тесты' },
  { value: 'teacher', label: 'Преподаватель', description: 'Может создавать и управлять тестами' },
  { value: 'admin', label: 'Администратор', description: 'Полный доступ к системе' }
]

const canChangeRoles = computed(() => {
  // Only admin can change roles
  return currentUser.value?.role === 'admin'
})

const canViewUsers = computed(() => {
  // Both teacher and admin can view users
  const role = currentUser.value?.role
  return role === 'teacher' || role === 'admin'
})

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  try {
    await usersStore.fetchUsers(currentPage.value)
  } catch (err) {
    // Error is handled by store
  }
}

function handlePageChange(page: number) {
  usersStore.fetchUsers(page)
}

function handleRetry() {
  usersStore.clearError()
  loadUsers()
}

function getRoleBadgeClass(role: string): string {
  const classes: Record<string, string> = {
    admin: 'bg-cyber-pink/20 text-cyber-pink border border-cyber-pink/30',
    teacher: 'bg-neon-orange/20 text-neon-orange border border-neon-orange/30',
    student: 'bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/30'
  }
  return classes[role] || 'bg-text-muted/20 text-text-muted border border-text-muted/30'
}

function formatRole(role: string): string {
  return role.charAt(0).toUpperCase() + role.slice(1)
}

function openRoleModal(user: User) {
  selectedUser.value = user
  selectedRole.value = user.role
  showRoleModal.value = true
}

function closeRoleModal() {
  showRoleModal.value = false
  selectedUser.value = null
  selectedRole.value = ''
}

async function handleRoleChange() {
  if (!selectedUser.value || !selectedRole.value) return

  isProcessing.value = true
  try {
    await usersStore.updateUserRole(
      selectedUser.value.id,
      selectedRole.value as 'admin' | 'teacher' | 'student'
    )
    closeRoleModal()
  } catch (err) {
    // Error is handled by store
  } finally {
    isProcessing.value = false
  }
}
</script>
