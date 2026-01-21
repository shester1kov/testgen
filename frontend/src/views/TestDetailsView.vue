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

      <!-- Loading State -->
      <div v-if="testsStore.loading" class="animate-pulse">
        <div class="h-8 bg-cyber-blue/20 rounded w-1/3 mb-2"></div>
        <div class="h-4 bg-cyber-blue/10 rounded w-1/4"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="card-cyber border-red-500/20 bg-red-500/5 p-6">
        <div class="flex items-center gap-3 mb-4">
          <svg class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 class="text-lg font-semibold text-red-400">Ошибка загрузки теста</h3>
        </div>
        <p class="text-red-300 mb-4">{{ error }}</p>
        <button @click="loadTest" class="btn-neon">
          Попробовать снова
        </button>
      </div>

      <!-- Test Header -->
      <div v-else-if="test">
        <div class="flex items-start justify-between mb-4">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h1 class="text-3xl font-bold text-text-primary">{{ test.title }}</h1>
              <button
                @click="showTestEditModal = true"
                class="p-2 text-cyber-blue hover:text-neon-orange transition-colors"
                title="Редактировать тест"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
            </div>
            <p v-if="test.description" class="text-text-muted">{{ test.description }}</p>
          </div>
          <div class="flex items-center gap-3">
            <!-- Export Dropdown -->
            <ExportDropdown
              :test-id="test.id"
              @error="handleExportError"
            />
            <span
              :class="getStatusClass(test.status)"
              class="px-4 py-2 rounded-full text-sm font-medium"
            >
              {{ test.status }}
            </span>
          </div>
        </div>

        <!-- Test Meta Info -->
        <div class="flex flex-wrap gap-4 text-sm text-text-secondary mb-6">
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Вопросов: {{ test.total_questions }}</span>
          </div>

          <div class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Создан {{ formatDate(test.created_at) }}</span>
          </div>

          <div v-if="test.moodle_synced" class="flex items-center gap-2 text-green-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>Синхронизирован с Moodle</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Questions List -->
    <div v-if="test && test.questions && test.questions.length > 0" class="space-y-6">
      <div
        v-for="(question, index) in sortedQuestions"
        :key="question.id"
        class="card-cyber"
      >
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-start gap-4 flex-1">
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-neon-orange/20 flex items-center justify-center text-neon-orange font-semibold">
              {{ index + 1 }}
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span
                  :class="getQuestionTypeClass(question.question_type)"
                  class="px-3 py-1 rounded-full text-xs font-medium"
                >
                  {{ formatQuestionType(question.question_type) }}
                </span>
                <span
                  :class="getDifficultyClass(question.difficulty)"
                  class="px-3 py-1 rounded-full text-xs font-medium"
                >
                  {{ question.difficulty }}
                </span>
                <span class="text-xs text-text-muted">{{ question.points }} pts</span>
              </div>
              <p class="text-text-primary text-lg mb-4">{{ question.question_text }}</p>

              <!-- Edit Button -->
              <button
                @click="openEditModal(question)"
                class="mb-3 px-3 py-1.5 text-sm bg-cyber-blue/20 hover:bg-cyber-blue/30 border border-cyber-blue/50 rounded text-cyber-blue font-medium transition-colors flex items-center gap-1"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Редактировать вопрос
              </button>

              <!-- Answers -->
              <div class="space-y-2">
                <div
                  v-for="answer in sortedAnswers(question.answers)"
                  :key="answer.id"
                  :class="[
                    'flex items-center gap-3 p-3 rounded-lg border transition-colors',
                    answer.is_correct
                      ? 'border-green-500/30 bg-green-500/10'
                      : 'border-cyber-blue/20 bg-bg-secondary/50'
                  ]"
                >
                  <div
                    :class="[
                      'flex-shrink-0 w-5 h-5 rounded-full border-2 flex items-center justify-center',
                      answer.is_correct
                        ? 'border-green-500 bg-green-500/20'
                        : 'border-cyber-blue/50'
                    ]"
                  >
                    <svg
                      v-if="answer.is_correct"
                      class="w-3 h-3 text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <span :class="answer.is_correct ? 'text-green-400' : 'text-text-secondary'">
                    {{ answer.answer_text }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- No Questions State -->
    <div v-else-if="test && (!test.questions || test.questions.length === 0)" class="card-cyber text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-yellow-500/20 flex items-center justify-center">
        <svg class="w-8 h-8 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-text-primary mb-2">Вопросов пока нет</h3>
      <p class="text-text-muted">В этом тесте пока нет вопросов.</p>
    </div>

    <!-- Question Edit Modal -->
    <QuestionEditModal
      :show="showEditModal"
      :question="editingQuestion"
      :test-id="id"
      @close="closeEditModal"
      @saved="handleQuestionSaved"
    />

    <!-- Test Edit Modal -->
    <TestEditModal
      :show="showTestEditModal"
      :test="test"
      @close="showTestEditModal = false"
      @saved="handleTestSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useTestsStore } from '@/features/tests/stores/testsStore'
import { TestStatus, QuestionType, Difficulty, type Question, type Answer } from '@/features/tests/types/test.types'
import QuestionEditModal from '@/features/tests/components/QuestionEditModal.vue'
import TestEditModal from '@/features/tests/components/TestEditModal.vue'
import ExportDropdown from '@/features/tests/components/ExportDropdown.vue'
import logger from '@/utils/logger'

const route = useRoute()
const testsStore = useTestsStore()
const id = route.params.id as string

const test = computed(() => testsStore.currentTest)
const error = ref<string | null>(null)

const showEditModal = ref(false)
const editingQuestion = ref<Question | null>(null)
const showTestEditModal = ref(false)

const sortedQuestions = computed(() => {
  if (!test.value?.questions) return []
  return [...test.value.questions].sort((a, b) => a.order_num - b.order_num)
})

function sortedAnswers(answers: Answer[]) {
  return [...answers].sort((a, b) => a.order_num - b.order_num)
}

async function loadTest() {
  error.value = null
  try {
    await testsStore.fetchTest(id)
  } catch (err: any) {
    error.value = err.message || 'Не удалось загрузить тест'
    logger.error('Failed to load test', 'TestDetailsView', err)
  }
}

function getStatusClass(status: string): string {
  switch (status) {
    case TestStatus.DRAFT:
      return 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
    case TestStatus.PUBLISHED:
      return 'bg-green-500/20 text-green-400 border border-green-500/30'
    case TestStatus.ARCHIVED:
      return 'bg-gray-500/20 text-gray-400 border border-gray-500/30'
    default:
      return 'bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/30'
  }
}

function getQuestionTypeClass(_type: string): string {
  return 'bg-purple-500/20 text-purple-400 border border-purple-500/30'
}

function getDifficultyClass(difficulty: string): string {
  switch (difficulty) {
    case Difficulty.EASY:
      return 'bg-green-500/20 text-green-400 border border-green-500/30'
    case Difficulty.MEDIUM:
      return 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
    case Difficulty.HARD:
      return 'bg-red-500/20 text-red-400 border border-red-500/30'
    default:
      return 'bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/30'
  }
}

function formatQuestionType(type: string): string {
  switch (type) {
    case QuestionType.SINGLE_CHOICE:
      return 'Один ответ'
    case QuestionType.MULTIPLE_CHOICE:
      return 'Множественный выбор'
    case QuestionType.TRUE_FALSE:
      return 'Верно/Неверно'
    case QuestionType.SHORT_ANSWER:
      return 'Короткий ответ'
    default:
      return type
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString()
}

function handleExportError(message: string) {
  error.value = message
  logger.error('Export error', 'TestDetailsView', new Error(message))
}

function openEditModal(question: Question) {
  editingQuestion.value = question
  showEditModal.value = true
}

function closeEditModal() {
  showEditModal.value = false
  editingQuestion.value = null
}

async function handleQuestionSaved(updatedQuestion: Question) {
  logger.info('Question saved, reloading test', 'TestDetailsView', { questionId: updatedQuestion.id })
  await loadTest()
}

async function handleTestSaved(updatedTest: any) {
  logger.info('Test saved, reloading test', 'TestDetailsView', { testId: updatedTest.id })
  await loadTest()
}

onMounted(() => {
  loadTest()
})
</script>
