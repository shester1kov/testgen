<template>
  <div class="document-card">
    <div class="card-header">
      <div class="file-icon-container">
        <component :is="getFileIcon()" class="file-icon" />
      </div>

      <div class="document-info">
        <h3 class="document-title">{{ document.title }}</h3>
        <p class="document-filename">{{ document.file_name }}</p>
        <div class="document-meta">
          <span class="meta-item">{{ formatFileSize(document.file_size) }}</span>
          <span class="meta-separator">•</span>
          <span class="meta-item">{{ formatDate(document.created_at) }}</span>
        </div>
      </div>

      <div class="status-badge" :class="getStatusClass()">
        {{ getStatusText() }}
      </div>
    </div>

    <div v-if="document.status === DocumentStatus.ERROR && document.error_msg" class="error-section">
      <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="error-text">{{ document.error_msg }}</p>
    </div>

    <div class="card-actions">
      <button v-if="canParse" type="button" class="action-button primary" :disabled="isProcessing"
        @click="handleParse">
        <svg class="button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        Parse Document
      </button>

      <button v-if="document.status === DocumentStatus.PARSED" type="button" class="action-button secondary"
        @click="handleViewText">
        <svg class="button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        View Text
      </button>

      <button type="button" class="action-button danger" :disabled="isProcessing" @click="handleDelete">
        <svg class="button-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
    [DocumentStatus.UPLOADED]: 'status-uploaded',
    [DocumentStatus.PARSING]: 'status-parsing',
    [DocumentStatus.PARSED]: 'status-parsed',
    [DocumentStatus.ERROR]: 'status-error',
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

<style scoped>
.document-card {
  @apply bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-shadow;
}

.card-header {
  @apply flex items-start gap-4 mb-4;
}

.file-icon-container {
  @apply flex-shrink-0 w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center;
}

.file-icon {
  @apply w-6 h-6 text-blue-600;
}

.document-info {
  @apply flex-1 min-w-0;
}

.document-title {
  @apply text-base font-semibold text-gray-900 truncate;
}

.document-filename {
  @apply text-sm text-gray-600 truncate;
}

.document-meta {
  @apply flex items-center gap-2 mt-1 text-xs text-gray-500;
}

.meta-item {
  @apply inline-block;
}

.meta-separator {
  @apply text-gray-400;
}

.status-badge {
  @apply px-3 py-1 rounded-full text-xs font-medium flex-shrink-0;
}

.status-uploaded {
  @apply bg-gray-100 text-gray-800;
}

.status-parsing {
  @apply bg-blue-100 text-blue-800;
}

.status-parsed {
  @apply bg-green-100 text-green-800;
}

.status-error {
  @apply bg-red-100 text-red-800;
}

.error-section {
  @apply flex items-start gap-2 p-3 mb-4 bg-red-50 border border-red-200 rounded-md;
}

.error-icon {
  @apply w-5 h-5 text-red-500 flex-shrink-0 mt-0.5;
}

.error-text {
  @apply text-sm text-red-700;
}

.card-actions {
  @apply flex gap-2 flex-wrap;
}

.action-button {
  @apply px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1.5;
  @apply transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.action-button.primary {
  @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
}

.action-button.secondary {
  @apply bg-gray-100 text-gray-700 hover:bg-gray-200 focus:ring-gray-500;
}

.action-button.danger {
  @apply bg-red-100 text-red-700 hover:bg-red-200 focus:ring-red-500;
}

.button-icon {
  @apply w-4 h-4;
}
</style>
