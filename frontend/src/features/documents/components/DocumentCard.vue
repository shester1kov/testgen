<template>
  <div class="group bg-dark-600 border border-dark-500 rounded-lg p-4 hover:border-neon-orange/50 transition-all duration-300">
    <div class="flex items-start gap-4">
      <!-- File Icon -->
      <div class="w-12 h-12 rounded-lg bg-cyber-blue/20 flex items-center justify-center flex-shrink-0 group-hover:shadow-cyber-blue transition-shadow">
        <svg class="w-6 h-6 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>

      <!-- Document Info -->
      <div class="flex-1 min-w-0">
        <h3 class="text-sm font-semibold text-text-primary truncate group-hover:text-neon-orange transition-colors">
          {{ document.title }}
        </h3>
        <p class="text-xs text-text-muted truncate mt-0.5">{{ document.file_name }}</p>
        <div v-if="document.user_name || document.user_email" class="flex items-center gap-2 mt-1">
          <span class="px-2 py-0.5 bg-cyber-blue/10 border border-cyber-blue/30 rounded text-xs text-cyber-blue">
            <span v-if="document.user_name">{{ document.user_name }}</span>
            <span v-if="document.user_email" class="text-text-muted ml-1">({{ document.user_email }})</span>
          </span>
        </div>
        <div class="flex items-center gap-3 mt-2 text-xs text-text-muted">
          <span>{{ formatFileSize(document.file_size) }}</span>
          <span class="text-dark-500">•</span>
          <span>{{ formatDate(document.created_at) }}</span>
        </div>
      </div>

      <!-- Status Badge -->
      <div class="flex-shrink-0">
        <span class="px-2.5 py-1 rounded-full text-xs font-medium" :class="getStatusClass()">
          {{ getStatusText() }}
        </span>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="document.status === DocumentStatus.ERROR && document.error_msg"
         class="mt-3 flex items-start gap-2 p-3 bg-cyber-pink/10 border border-cyber-pink/30 rounded-lg">
      <svg class="w-4 h-4 text-cyber-pink flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-xs text-cyber-pink">{{ document.error_msg }}</p>
    </div>

    <!-- Actions -->
    <div class="flex gap-2 mt-4">
      <button
        v-if="canParse"
        type="button"
        class="flex-1 px-3 py-2 bg-neon-orange/10 border border-neon-orange/30 text-neon-orange rounded-lg text-sm font-medium hover:bg-neon-orange/20 hover:shadow-neon-sm transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        :disabled="isProcessing"
        @click="handleParse"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        Parse
      </button>

      <button
        v-if="document.status === DocumentStatus.PARSED"
        type="button"
        class="flex-1 px-3 py-2 bg-cyber-blue/10 border border-cyber-blue/30 text-cyber-blue rounded-lg text-sm font-medium hover:bg-cyber-blue/20 transition-all flex items-center justify-center gap-2"
        @click="handleViewText"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        View
      </button>

      <button
        type="button"
        class="px-3 py-2 bg-cyber-pink/10 border border-cyber-pink/30 text-cyber-pink rounded-lg text-sm font-medium hover:bg-cyber-pink/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="isProcessing"
        @click="handleDelete"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
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

function getStatusClass() {
  const classes: Record<string, string> = {
    [DocumentStatus.UPLOADED]: 'bg-text-muted/20 text-text-muted border border-text-muted/30',
    [DocumentStatus.PARSING]: 'bg-cyber-blue/20 text-cyber-blue border border-cyber-blue/30',
    [DocumentStatus.PARSED]: 'bg-neon-orange/20 text-neon-orange border border-neon-orange/30',
    [DocumentStatus.ERROR]: 'bg-cyber-pink/20 text-cyber-pink border border-cyber-pink/30',
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

  // Reset time to midnight for accurate day comparison
  const dateOnly = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const nowOnly = new Date(now.getFullYear(), now.getMonth(), now.getDate())

  const diffInMs = nowOnly.getTime() - dateOnly.getTime()
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))

  if (diffInDays === 0) return 'Today'
  if (diffInDays === 1) return 'Yesterday'
  if (diffInDays > 1 && diffInDays < 7) return `${diffInDays} days ago`
  if (diffInDays < 0) return 'Just now' // Future dates (timezone issues)

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
