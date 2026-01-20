<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      class="px-4 py-2 bg-neon-orange/20 hover:bg-neon-orange/30 border border-neon-orange/50 rounded-lg text-neon-orange font-medium transition-colors flex items-center gap-2"
      :disabled="loading"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <span v-if="loading">Загрузка...</span>
      <span v-else>Экспорт</span>
      <svg
        class="w-4 h-4 transition-transform"
        :class="{ 'rotate-180': isOpen }"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-72 bg-bg-secondary/95 backdrop-blur-md border border-cyber-blue/30 rounded-lg shadow-xl z-50 overflow-hidden"
    >
      <div class="p-2">
        <div class="px-3 py-2 text-xs font-semibold text-text-muted uppercase tracking-wider border-b border-cyber-blue/20 mb-2">
          Выберите формат экспорта
        </div>

        <button
          v-for="format in formats"
          :key="format.format"
          @click="handleExport(format.format)"
          class="w-full text-left px-3 py-2 rounded hover:bg-cyber-blue/10 transition-colors group"
        >
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 mt-1">
              <div class="w-8 h-8 rounded bg-cyber-blue/20 flex items-center justify-center text-cyber-blue text-xs font-bold">
                {{ getFormatIcon(format.format) }}
              </div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-text-primary group-hover:text-neon-orange transition-colors">
                {{ format.name }}
              </div>
              <div class="text-xs text-text-muted mt-0.5">
                {{ format.description }}
              </div>
            </div>
          </div>
        </button>

        <div v-if="formats.length === 0 && !loading" class="px-3 py-4 text-center text-text-muted text-sm">
          Нет доступных форматов
        </div>

        <div v-if="loading" class="px-3 py-4 text-center text-text-muted text-sm">
          Загрузка форматов...
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { FormatInfo, ExportFormat } from '@/features/tests/types/export.types'
import { testService } from '@/services/testService'
import logger from '@/utils/logger'

interface Props {
  testId: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  export: [format: ExportFormat]
  error: [message: string]
}>()

const isOpen = ref(false)
const loading = ref(false)
const formats = ref<FormatInfo[]>([])
const dropdownRef = ref<HTMLElement | null>(null)

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

function closeDropdown() {
  isOpen.value = false
}

async function loadFormats() {
  loading.value = true
  try {
    const fetchedFormats = await testService.getExportFormats()
    formats.value = fetchedFormats
    logger.info('Export formats loaded', 'ExportDropdown', { count: fetchedFormats.length })
  } catch (err: any) {
    logger.error('Failed to load export formats', 'ExportDropdown', err)
    emit('error', err.message || 'Не удалось загрузить форматы экспорта')
  } finally {
    loading.value = false
  }
}

async function handleExport(format: ExportFormat) {
  closeDropdown()
  emit('export', format)

  try {
    await testService.exportToFormat(props.testId, format)
    logger.info('Test exported successfully', 'ExportDropdown', { format })
  } catch (err: any) {
    logger.error('Failed to export test', 'ExportDropdown', err)
    emit('error', err.message || `Не удалось экспортировать тест в формат ${format}`)
  }
}

function getFormatIcon(format: ExportFormat): string {
  const icons: Record<string, string> = {
    json: 'JS',
    moodle_xml: 'XML',
    stepik_csv: 'CSV',
    gift: 'GFT',
    aiken: 'AKN',
    blackboard: 'BB',
    qti: 'QTI',
  }
  return icons[format] || '?'
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    closeDropdown()
  }
}

onMounted(() => {
  loadFormats()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
