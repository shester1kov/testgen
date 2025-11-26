<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/60 backdrop-blur-sm" @click="handleClose"></div>

        <!-- Modal -->
        <div class="flex min-h-full items-center justify-center p-4">
          <div class="relative w-full max-w-3xl card-cyber">
            <!-- Header -->
            <div class="flex items-center justify-between mb-6 pb-4 border-b border-cyber-blue/20">
              <h2 class="text-2xl font-bold text-text-primary">Редактировать вопрос</h2>
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
              <!-- Question Text -->
              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">
                  Текст вопроса
                </label>
                <textarea
                  v-model="formData.question_text"
                  rows="4"
                  class="input-cyber"
                  placeholder="Введите текст вопроса..."
                  required
                ></textarea>
              </div>

              <!-- Question Type and Difficulty -->
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-text-secondary mb-2">
                    Тип вопроса
                  </label>
                  <select v-model="formData.question_type" class="input-cyber" required>
                    <option value="single_choice">Один ответ</option>
                    <option value="multiple_choice">Множественный выбор</option>
                    <option value="true_false">Верно/Неверно</option>
                    <option value="short_answer">Короткий ответ</option>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium text-text-secondary mb-2">
                    Сложность
                  </label>
                  <select v-model="formData.difficulty" class="input-cyber" required>
                    <option value="easy">Легкий</option>
                    <option value="medium">Средний</option>
                    <option value="hard">Сложный</option>
                  </select>
                </div>
              </div>

              <!-- Points -->
              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">
                  Баллы
                </label>
                <input
                  v-model.number="formData.points"
                  type="number"
                  min="0.5"
                  step="0.5"
                  class="input-cyber"
                  placeholder="Баллы за этот вопрос"
                  required
                />
              </div>

              <!-- Answers -->
              <div>
                <div class="flex items-center justify-between mb-3">
                  <label class="block text-sm font-medium text-text-secondary">
                    Ответы
                  </label>
                  <button
                    type="button"
                    @click="addAnswer"
                    class="px-3 py-1 text-sm bg-cyber-blue/20 hover:bg-cyber-blue/30 border border-cyber-blue/50 rounded text-cyber-blue font-medium transition-colors flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                    </svg>
                    Добавить ответ
                  </button>
                </div>

                <div class="space-y-3">
                  <div
                    v-for="(answer, index) in formData.answers"
                    :key="index"
                    class="flex items-start gap-3 p-3 bg-bg-secondary/50 border border-cyber-blue/20 rounded-lg"
                  >
                    <div class="flex items-center gap-3 flex-1">
                      <!-- Correct Checkbox -->
                      <input
                        v-model="answer.is_correct"
                        type="checkbox"
                        @change="handleCorrectChange(index)"
                        class="w-5 h-5 rounded border-cyber-blue/50 bg-bg-secondary text-cyber-blue focus:ring-2 focus:ring-cyber-blue/50"
                      />

                      <!-- Answer Text -->
                      <input
                        v-model="answer.answer_text"
                        type="text"
                        class="flex-1 input-cyber"
                        placeholder="Текст ответа..."
                        required
                      />
                    </div>

                    <!-- Remove Button -->
                    <button
                      type="button"
                      @click="removeAnswer(index)"
                      class="text-red-400 hover:text-red-300 transition-colors flex-shrink-0"
                      :disabled="formData.answers.length <= 2"
                    >
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                </div>

                <p class="text-xs text-text-muted mt-2">
                  {{ formData.question_type === 'single_choice'
                    ? 'Отметьте галочкой один правильный ответ'
                    : 'Отметьте галочкой правильные ответы' }}
                </p>
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
import type { Question, Answer } from '@/features/tests/types/test.types'
import logger from '@/utils/logger'

interface Props {
  show: boolean
  question: Question | null
  testId: string
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', question: Question): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

interface AnswerFormData {
  id?: string
  answer_text: string
  is_correct: boolean
  order_num: number
}

interface QuestionFormData {
  question_text: string
  question_type: string
  difficulty: string
  points: number
  answers: AnswerFormData[]
}

const loading = ref(false)
const error = ref<string | null>(null)

const formData = ref<QuestionFormData>({
  question_text: '',
  question_type: 'single_choice',
  difficulty: 'medium',
  points: 1.0,
  answers: [
    { answer_text: '', is_correct: true, order_num: 0 },
    { answer_text: '', is_correct: false, order_num: 1 },
  ],
})

const isFormValid = computed(() => {
  if (!formData.value.question_text.trim()) return false
  if (formData.value.points <= 0) return false
  if (formData.value.answers.length < 2) return false

  const hasEmptyAnswer = formData.value.answers.some(a => !a.answer_text.trim())
  if (hasEmptyAnswer) return false

  const hasCorrectAnswer = formData.value.answers.some(a => a.is_correct)
  if (!hasCorrectAnswer) return false

  return true
})

watch(() => props.question, (newQuestion) => {
  if (newQuestion) {
    formData.value = {
      question_text: newQuestion.question_text,
      question_type: newQuestion.question_type,
      difficulty: newQuestion.difficulty,
      points: newQuestion.points,
      answers: newQuestion.answers.map((a, index) => ({
        id: a.id,
        answer_text: a.answer_text,
        is_correct: a.is_correct,
        order_num: index,
      })),
    }
  }
}, { immediate: true })

function addAnswer() {
  formData.value.answers.push({
    answer_text: '',
    is_correct: false,
    order_num: formData.value.answers.length,
  })
}

function removeAnswer(index: number) {
  if (formData.value.answers.length > 2) {
    formData.value.answers.splice(index, 1)
    // Re-index order_num
    formData.value.answers.forEach((answer, idx) => {
      answer.order_num = idx
    })
  }
}

function handleCorrectChange(index: number) {
  // For single_choice, only one answer can be correct
  if (formData.value.question_type === 'single_choice' && formData.value.answers[index].is_correct) {
    formData.value.answers.forEach((answer, idx) => {
      if (idx !== index) {
        answer.is_correct = false
      }
    })
  }
}

async function handleSubmit() {
  if (!props.question || !isFormValid.value) return

  loading.value = true
  error.value = null

  try {
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL
    const url = `${API_BASE_URL}/tests/${props.testId}/questions/${props.question.id}`

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
      throw new Error(errorData.error || 'Не удалось обновить вопрос')
    }

    const updatedQuestion = await response.json()
    logger.info('Question updated successfully', 'QuestionEditModal', { questionId: props.question.id })

    emit('saved', updatedQuestion)
    emit('close')
  } catch (err: any) {
    error.value = err.message || 'Не удалось обновить вопрос'
    logger.error('Failed to update question', 'QuestionEditModal', err)
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
