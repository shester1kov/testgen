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
        Back to Tests
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
          <h3 class="text-lg font-semibold text-red-400">Error loading test</h3>
        </div>
        <p class="text-red-300 mb-4">{{ error }}</p>
        <button @click="loadTest" class="btn-neon">
          Try Again
        </button>
      </div>

      <!-- Test Header -->
      <div v-else-if="test">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h1 class="text-3xl font-bold text-text-primary mb-2">{{ test.title }}</h1>
            <p v-if="test.description" class="text-text-muted">{{ test.description }}</p>
          </div>
          <span
            :class="getStatusClass(test.status)"
            class="px-4 py-2 rounded-full text-sm font-medium"
          >
            {{ test.status }}
          </span>
        </div>

        <!-- Test Meta Info -->
        <div class="flex flex-wrap gap-4 text-sm text-text-secondary mb-6">
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{{ test.total_questions }} questions</span>
          </div>

          <div class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Created {{ formatDate(test.created_at) }}</span>
          </div>

          <div v-if="test.moodle_synced" class="flex items-center gap-2 text-green-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span>Synced to Moodle</span>
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
      <h3 class="text-xl font-semibold text-text-primary mb-2">No questions yet</h3>
      <p class="text-text-muted">This test doesn't have any questions yet.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useTestsStore } from '@/features/tests/stores/testsStore'
import { TestStatus, QuestionType, Difficulty, type Question, type Answer } from '@/features/tests/types/test.types'
import logger from '@/utils/logger'

const route = useRoute()
const testsStore = useTestsStore()
const id = route.params.id as string

const test = computed(() => testsStore.currentTest)
const error = ref<string | null>(null)

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
    error.value = err.message || 'Failed to load test'
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

function getQuestionTypeClass(type: string): string {
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
      return 'Single Choice'
    case QuestionType.MULTIPLE_CHOICE:
      return 'Multiple Choice'
    case QuestionType.TRUE_FALSE:
      return 'True/False'
    case QuestionType.SHORT_ANSWER:
      return 'Short Answer'
    default:
      return type
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleString()
}

onMounted(() => {
  loadTest()
})
</script>
