<template>
  <div class="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-shadow">
    <div class="flex items-start gap-4 mb-4">
      <div class="flex-shrink-0 w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
        <component :is="getFileIcon()" class="w-6 h-6 text-blue-600" />
      </div>

      <div class="flex-1 min-w-0">
        <h3 class="text-base font-semibold text-gray-900 truncate">{{ document.title }}</h3>
        <p class="text-sm text-gray-600 truncate">{{ document.file_name }}</p>
        <div class="flex items-center gap-2 mt-1 text-xs text-gray-500">
          <span class="inline-block">{{ formatFileSize(document.file_size) }}</span>
          <span class="text-gray-400">•</span>
          <span class="inline-block">{{ formatDate(document.created_at) }}</span>
        </div>
      </div>

      <div class="px-3 py-1 rounded-full text-xs font-medium flex-shrink-0" :class="getStatusClass()">
        {{ getStatusText() }}
      </div>
    </div>

    <div v-if="document.status === DocumentStatus.ERROR && document.error_msg"
         class="flex items-start gap-2 p-3 mb-4 bg-red-50 border border-red-200 rounded-md">
      <svg class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-red-700">{{ document.error_msg }}</p>
    </div>

    <div class="flex gap-2 flex-wrap">
      <button v-if="canParse" type="button"
              class="px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1.5 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500"
              :disabled="isProcessing"
              @click="handleParse">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        Parse Document
      </button>

      <button v-if="document.status === DocumentStatus.PARSED" type="button"
              class="px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1.5 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 bg-gray-100 text-gray-700 hover:bg-gray-200 focus:ring-gray-500"
              @click="handleViewText">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        View Text
      </button>

      <button type="button"
              class="px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1.5 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed bg-red-100 text-red-700 hover:bg-red-200 focus:ring-red-500"
              :disabled="isProcessing"
              @click="handleDelete">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
        Delete
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Document } from '../types/document.types'
import { DocumentStatus, FileType } from '../types/document.types'
import { useDocumentsStore } from '../stores/documentsStore'

interface Props {
  document: Document
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'view-text', document: Document): void
}>()

const documentsStore = useDocumentsStore()
const isProcessing = ref(false)

const canParse = computed(() => {
  return props.document.status === DocumentStatus.UPLOADED ||
    props.document.status === DocumentStatus.ERROR
})

function getFileIcon() {
  const icons: Record<string, string> = {
    [FileType.PDF]: 'pdf-icon',
    [FileType.DOCX]: 'doc-icon',
    [FileType.PPTX]: 'ppt-icon',
    [FileType.TXT]: 'txt-icon',
    [FileType.MD]: 'md-icon',
  }
  // Return generic document icon SVG component
  return 'svg'
}

function getStatusClass() {
  const classes: Record<string, string> = {
    [DocumentStatus.UPLOADED]: 'bg-gray-100 text-gray-800',
    [DocumentStatus.PARSING]: 'bg-blue-100 text-blue-800',
    [DocumentStatus.PARSED]: 'bg-green-100 text-green-800',
    [DocumentStatus.ERROR]: 'bg-red-100 text-red-800',
  }
  return classes[props.document.status] || ''
}

function getStatusText() {
  const texts: Record<string, string> = {
    [DocumentStatus.UPLOADED]: 'Uploaded',
    [DocumentStatus.PARSING]: 'Parsing...',
    [DocumentStatus.PARSED]: 'Parsed',
    [DocumentStatus.ERROR]: 'Error',
  }
  return texts[props.document.status] || 'Unknown'
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffInMs = now.getTime() - date.getTime()
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))

  if (diffInDays === 0) return 'Today'
  if (diffInDays === 1) return 'Yesterday'
  if (diffInDays < 7) return `${diffInDays} days ago`

  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

async function handleParse() {
  isProcessing.value = true
  try {
    await documentsStore.parseDocument(props.document.id)
  } catch (error) {
    // Error is handled by store
  } finally {
    isProcessing.value = false
  }
}

function handleViewText() {
  emit('view-text', props.document)
}

async function handleDelete() {
  if (!confirm('Are you sure you want to delete this document?')) {
    return
  }

  isProcessing.value = true
  try {
    await documentsStore.deleteDocument(props.document.id)
  } catch (error) {
    // Error is handled by store
  } finally {
    isProcessing.value = false
  }
}
</script>
