<template>
  <DesignModeBanner />
  <component :is="layout">
    <router-view :key="route.fullPath" />
  </component>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/authStore'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import DesignModeBanner from '@/components/DesignModeBanner.vue'

const route = useRoute()
const authStore = useAuthStore()

const layout = computed(() => {
  const layoutName = route.meta.layout as string
  return layoutName === 'auth' ? AuthLayout : DefaultLayout
})

onMounted(() => {
  authStore.initializeAuth()
})
</script>
