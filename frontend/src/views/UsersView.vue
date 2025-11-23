<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">Users Management</h1>
        <p class="text-text-secondary">
          {{ isAdmin
            ? 'Manage user roles and permissions'
            : 'View students and assign tests to them'
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
