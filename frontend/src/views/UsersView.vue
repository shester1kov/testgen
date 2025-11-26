<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">Управление пользователями</h1>
        <p class="text-text-secondary">
          {{ isAdmin
            ? 'Управляйте ролями и правами пользователей'
            : 'Просматривайте студентов и назначайте им тесты'
          }}
        </p>
      </div>
    </div>

    <UserList />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/features/auth/stores/authStore'
import UserList from '@/features/users/components/UserList.vue'

const authStore = useAuthStore()

const isAdmin = computed(() => {
  return authStore.user?.role === 'admin'
})
</script>
