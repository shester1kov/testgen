<template>
  <div>
    <div class="mb-8">
      <router-link
        to="/tests"
        class="inline-flex items-center text-text-secondary hover:text-neon-orange transition-colors mb-4"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Назад к тестам
      </router-link>
      <h1 class="text-3xl font-bold text-text-primary mb-2">Создать тест</h1>
      <p class="text-text-secondary">Генерируйте тестовые вопросы из ваших документов с помощью ИИ</p>
    </div>

    <div class="card-cyber">
      <form @submit.prevent="handleSubmit">
        <!-- Select Document -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-text-primary mb-2">
            Выберите документ
            <span class="text-red-500">*</span>
          </label>

          <!-- Loading State -->
          <div v-if="isLoadingDocuments" class="flex items-center gap-2 text-text-secondary mb-2">
            <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span class="text-sm">Загрузка документов...</span>
          </div>

          <select
            v-model="form.documentId"
            class="input-cyber w-full"
            :disabled="isLoadingDocuments"
            required
          >
            <option value="">
              {{ isLoadingDocuments ? 'Загрузка...' : '-- Выберите документ --' }}
            </option>
            <option
              v-for="doc in parsedDocuments"
              :key="doc.id"
              :value="doc.id"
            >
              {{ doc.title }} ({{ doc.file_type.toUpperCase() }})
            </option>
          </select>

          <p v-if="!isLoadingDocuments && parsedDocuments.length === 0" class="text-text-muted text-sm mt-2">
            Нет обработанных документов. Пожалуйста, загрузите и обработайте документ на
            <router-link to="/documents" class="text-cyber-blue hover:underline">странице документов</router-link>.
          </p>
          <p v-else-if="!isLoadingDocuments && parsedDocuments.length > 0" class="text-text-muted text-sm mt-2">
            Доступно документов: {{ parsedDocuments.length }}
          </p>
        </div>

        <!-- Test Title -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-text-primary mb-2">
            Название теста
            <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            class="input-cyber w-full"
            placeholder="Например, Введение в программирование"
            required
            minlength="3"
          />
        </div>

        <!-- Number of Questions -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-text-primary mb-2">
            Количество вопросов
            <span class="text-red-500">*</span>
          </label>
          <input
            v-model.number="form.numQuestions"
            type="number"
            class="input-cyber w-full"
            min="1"
            max="50"
            required
          />
        </div>

        <!-- Difficulty -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-text-primary mb-2">
            Сложность
            <span class="text-red-500">*</span>
          </label>
          <div class="grid grid-cols-3 gap-4">
            <button
              type="button"
              @click.prevent="form.difficulty = 'easy'"
              :class="[
                'py-3 px-4 rounded-lg font-medium transition-all cursor-pointer',
                form.difficulty === 'easy'
                  ? 'bg-green-500 text-white border-2 border-green-500 shadow-lg'
                  : 'bg-bg-secondary text-text-secondary border-2 border-border-primary hover:border-green-500 hover:bg-green-500/10'
              ]"
            >
              Легкий
            </button>
            <button
              type="button"
              @click.prevent="form.difficulty = 'medium'"
              :class="[
                'py-3 px-4 rounded-lg font-medium transition-all cursor-pointer',
                form.difficulty === 'medium'
                  ? 'bg-yellow-500 text-white border-2 border-yellow-500 shadow-lg'
                  : 'bg-bg-secondary text-text-secondary border-2 border-border-primary hover:border-yellow-500 hover:bg-yellow-500/10'
              ]"
            >
              Средний
            </button>
            <button
              type="button"
              @click.prevent="form.difficulty = 'hard'"
              :class="[
                'py-3 px-4 rounded-lg font-medium transition-all cursor-pointer',
                form.difficulty === 'hard'
                  ? 'bg-red-500 text-white border-2 border-red-500 shadow-lg'
                  : 'bg-bg-secondary text-text-secondary border-2 border-border-primary hover:border-red-500 hover:bg-red-500/10'
              ]"
            >
              Сложный
            </button>
          </div>
        </div>

        <!-- LLM Provider -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-text-primary mb-2">
            ИИ Провайдер
          </label>
          <select
            v-model="form.llmProvider"
            class="input-cyber w-full"
          >
            <option value="yandexgpt">YandexGPT (Рекомендуется для России)</option>
            <option value="perplexity">Perplexity AI</option>
            <option value="openai">OpenAI GPT-4</option>
          </select>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="mb-6 p-4 bg-red-500/10 border border-red-500 rounded-lg">
          <p class="text-red-400 text-sm">{{ errorMessage }}</p>
        </div>

        <!-- Submit Button -->
        <div class="flex justify-end gap-4 items-center">
          <router-link
            to="/tests"
            class="inline-flex items-center justify-center px-6 py-3 border-2 border-border-primary text-text-secondary rounded-lg hover:border-neon-orange hover:text-neon-orange transition-all font-medium"
          >
            Отмена
          </router-link>
          <button
            type="submit"
            class="btn-neon"
            :disabled="isLoading || isLoadingDocuments || parsedDocuments.length === 0"
          >
            <svg
              v-if="isLoading"
              class="animate-spin -ml-1 mr-3 h-5 w-5 inline-block"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ isLoading ? 'Генерация...' : isLoadingDocuments ? 'Загрузка...' : 'Сгенерировать тест' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useDocumentsStore } from '@/features/documents/stores/documentsStore'
import testService from '@/services/testService'
import logger from '@/utils/logger'

const router = useRouter()
const documentsStore = useDocumentsStore()

const form = ref({
  documentId: '',
  title: '',
  numQuestions: 10,
  difficulty: 'medium',
  llmProvider: 'yandexgpt'
})

const isLoading = ref(false)
const isLoadingDocuments = ref(false)
const errorMessage = ref('')

const parsedDocuments = computed(() => {
  const filtered = documentsStore.documents.filter(doc => doc.status === 'parsed')
  logger.debug('Parsed documents computed', 'CreateTestView', {
    total: documentsStore.documents.length,
    parsed: filtered.length,
    documents: filtered.map(d => ({ id: d.id, title: d.title, status: d.status }))
  })
  return filtered
})

// Watch for changes in documents store
watch(() => documentsStore.documents, (newDocs) => {
  logger.debug('Documents store changed', 'CreateTestView', {
    count: newDocs.length,
    parsed: newDocs.filter(d => d.status === 'parsed').length
  })
}, { deep: true })

onMounted(async () => {
  // Always load documents to ensure fresh data
  isLoadingDocuments.value = true
  try {
    logger.info('Loading documents for test creation', 'CreateTestView')
    await documentsStore.fetchDocuments()
    logger.info('Documents loaded', 'CreateTestView', {
      total: documentsStore.documents.length,
      parsed: parsedDocuments.value.length
    })
  } catch (err: any) {
    logger.error('Failed to load documents', 'CreateTestView', err)
    errorMessage.value = 'Не удалось загрузить документы. Пожалуйста, попробуйте снова.'
  } finally {
    isLoadingDocuments.value = false
  }
})

async function handleSubmit() {
  errorMessage.value = ''
  isLoading.value = true

  try {
    logger.info('Generating test', 'CreateTestView', form.value)

    const response = await testService.generateTest({
      document_id: form.value.documentId,
      title: form.value.title,
      num_questions: form.value.numQuestions,
      difficulty: form.value.difficulty,
      llm_provider: form.value.llmProvider,
      question_types: ['single_choice'] // Default for now
    })

    logger.info('Test generated successfully', 'CreateTestView', { testId: response.id })

    // Redirect to test details or edit page
    router.push(`/tests/${response.id}`)
  } catch (err: any) {
    logger.error('Failed to generate test', 'CreateTestView', err)
    errorMessage.value = err.response?.data?.error || 'Не удалось сгенерировать тест. Пожалуйста, попробуйте снова.'
  } finally {
    isLoading.value = false
  }
}
</script>
