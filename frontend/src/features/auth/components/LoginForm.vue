<template>
  <div>
    <h2 class="text-2xl font-bold text-text-primary mb-6">Вход в систему</h2>

    <form @submit.prevent="handleLogin" class="space-y-6">
      <div>
        <label for="email" class="block text-sm font-medium text-text-secondary mb-2">Email</label>
        <input
          id="email"
          v-model="formData.email"
          type="email"
          required
          placeholder="your.email@example.com"
          class="input-neon w-full"
        />
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-text-secondary mb-2">Пароль</label>
        <input
          id="password"
          v-model="formData.password"
          type="password"
          required
          placeholder="Введите пароль"
          class="input-neon w-full"
        />
      </div>

      <div v-if="error" class="text-sm text-cyber-pink bg-cyber-pink/10 border border-cyber-pink/20 rounded-lg px-4 py-3">
        {{ error }}
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="btn-neon w-full"
      >
        <span v-if="loading" class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Вход...
        </span>
        <span v-else>Войти</span>
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-text-secondary">
      Нет аккаунта?
      <router-link to="/register" class="font-medium text-neon-orange hover:text-neon-orange-light transition-colors">
        Зарегистрироваться
      </router-link>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(formData.value)
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>
