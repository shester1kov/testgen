<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/60 backdrop-blur-sm" @click="handleClose"></div>

        <!-- Modal -->
        <div class="flex min-h-full items-center justify-center p-4">
          <div class="relative w-full max-w-2xl card-cyber">
            <!-- Header -->
            <div class="flex items-center justify-between mb-6 pb-4 border-b border-cyber-blue/20">
              <h2 class="text-2xl font-bold text-text-primary">Редактировать тест</h2>
              <button
                @click="handleClose"
                class="text-text-muted hover:text-neon-orange transition-colors"
              >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Error Display -->
            <div v-if="error" class="mb-4 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p class="text-red-400">{{ error }}</p>
            </div>

            <!-- Form -->
            <form @submit.prevent="handleSubmit" class="space-y-6">
              <!-- Title -->
              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">
                  Название
                </label>
                <input
                  v-model="formData.title"
                  type="text"
                  class="input-cyber"
                  placeholder="Введите название теста..."
                  required
                  minlength="3"
                />
              </div>

              <!-- Description -->
              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">
                  Описание
                </label>
                <textarea
                  v-model="formData.description"
                  rows="4"
                  class="input-cyber"
                  placeholder="Введите описание теста (необязательно)..."
                ></textarea>
              </div>

              <!-- Actions -->
              <div class="flex items-center justify-end gap-3 pt-4 border-t border-cyber-blue/20">
                <button
                  type="button"
                  @click="handleClose"
                  class="px-6 py-2 bg-bg-secondary hover:bg-bg-secondary/70 border border-cyber-blue/20 rounded-lg text-text-secondary font-medium transition-colors"
                  :disabled="loading"
                >
                  Отмена
                </button>
                <button
                  type="submit"
                  class="btn-neon"
                  :disabled="loading || !isFormValid"
                >
                  {{ loading ? 'Сохранение...' : 'Сохранить изменения' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Test } from '@/features/tests/types/test.types'
import logger from '@/utils/logger'

interface Props {
  show: boolean
  test: Test | null
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', test: Test): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

interface TestFormData {
  title: string
  description: string
}

const loading = ref(false)
const error = ref<string | null>(null)

const formData = ref<TestFormData>({
  title: '',
  description: '',
})

const isFormValid = computed(() => {
  return formData.value.title.trim().length >= 3
})

watch(() => props.test, (newTest) => {
  if (newTest) {
    formData.value = {
      title: newTest.title,
      description: newTest.description || '',
    }
  }
}, { immediate: true })

async function handleSubmit() {
  if (!props.test || !isFormValid.value) return

  loading.value = true
  error.value = null

  try {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
    const url = `${API_BASE_URL}/tests/${props.test.id}`

    const response = await fetch(url, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData.value),
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: response.statusText }))
      throw new Error(errorData.error || 'Не удалось обновить тест')
    }

    const updatedTest = await response.json()
    logger.info('Test updated successfully', 'TestEditModal', { testId: props.test.id })

    emit('saved', updatedTest)
    emit('close')
  } catch (err: any) {
    error.value = err.message || 'Не удалось обновить тест'
    logger.error('Failed to update test', 'TestEditModal', err)
  } finally {
    loading.value = false
  }
}

function handleClose() {
  if (!loading.value) {
    emit('close')
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .card-cyber,
.modal-leave-active .card-cyber {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from .card-cyber,
.modal-leave-to .card-cyber {
  transform: scale(0.95);
  opacity: 0;
}
</style>
