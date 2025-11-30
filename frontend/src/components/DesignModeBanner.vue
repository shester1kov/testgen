<template>
  <div
    v-if="showBanner"
    class="fixed top-0 left-0 right-0 bg-gradient-to-r from-orange-500 to-pink-500 text-white px-4 py-2 text-center text-sm font-medium z-50 shadow-lg"
  >
    <div class="flex items-center justify-center gap-2">
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
          clip-rule="evenodd"
        />
      </svg>
      <span>🎨 DESIGN MODE: Аутентификация отключена | Mock пользователь: {{ mockUser.email }}</span>
      <button
        @click="hideBanner"
        class="ml-4 text-white/80 hover:text-white transition-colors"
        title="Скрыть баннер"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { isDesignMode, getMockUser } from '@/utils/designMode'

const showBanner = ref(false)
const mockUser = getMockUser()

onMounted(() => {
  // Show banner only in design mode and if user hasn't dismissed it
  if (isDesignMode()) {
    const dismissed = sessionStorage.getItem('design-mode-banner-dismissed')
    showBanner.value = !dismissed
  }
})

function hideBanner() {
  showBanner.value = false
  sessionStorage.setItem('design-mode-banner-dismissed', 'true')
}
</script>
