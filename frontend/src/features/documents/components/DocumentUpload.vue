<template>
  <div class="card-cyber">
    <h2 class="text-lg font-semibold text-text-primary mb-4">Upload Document</h2>

    <div
      class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-all duration-300"
      :class="isDragOver ? 'border-neon-orange bg-neon-orange/10 shadow-neon-sm' : 'border-dark-500 hover:border-dark-600'"
      @click="fileInput?.click()"
      @drop.prevent="handleDrop"
      @dragover.prevent="isDragOver = true"
      @dragleave.prevent="isDragOver = false"
    >
      <input ref="fileInput" type="file" :accept="acceptedFormats" class="sr-only" @change="handleFileSelect" />

      <div class="flex flex-col items-center gap-3">
        <div class="w-12 h-12 rounded-full bg-neon-orange/20 flex items-center justify-center">
          <svg class="w-6 h-6 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
        </div>

        <p class="text-base font-medium text-text-primary">
          {{ isDragOver ? 'Drop file here' : 'Upload a document' }}
        </p>

        <p class="text-sm text-text-muted">
          Drag and drop or
          <button type="button" class="text-neon-orange hover:text-neon-orange-light font-medium underline" @click.stop>
            browse
          </button>
        </p>

        <p class="text-xs text-text-muted">Supported: PDF, DOCX, PPTX, TXT, MD (max 50MB)</p>
      </div>
    </div>

    <!-- Selected file info -->
    <div v-if="selectedFile" class="mt-4 p-4 bg-dark-600 rounded-lg border border-dark-500">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded bg-cyber-blue/20 flex items-center justify-center flex-shrink-0">
            <svg class="w-5 h-5 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div class="min-w-0">
            <p class="text-sm font-medium text-text-primary truncate">{{ selectedFile.name }}</p>
            <p class="text-xs text-text-muted">{{ formatFileSize(selectedFile.size) }}</p>
          </div>
        </div>

        <button type="button" class="p-1 text-text-muted hover:text-cyber-pink transition-colors" @click="clearFile">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Title input -->
    <div v-if="selectedFile" class="mt-4">
      <label for="doc-title" class="block text-sm font-medium text-text-secondary mb-2">Document Title (optional)</label>
      <input
        id="doc-title"
        v-model="title"
        type="text"
        class="input-neon w-full"
        :placeholder="selectedFile.name"
      />
    </div>

    <!-- Upload button -->
    <div v-if="selectedFile" class="mt-4 flex gap-3 justify-end">
      <button
        type="button"
        class="px-4 py-2 border border-dark-500 rounded-lg text-sm font-medium text-text-secondary hover:bg-dark-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="isUploading"
        @click="clearFile"
      >
        Cancel
      </button>
      <button
        type="button"
        class="btn-neon flex items-center gap-2 disabled:opacity-50"
        :disabled="isUploading"
        @click="handleUpload"
      >
        <svg v-if="isUploading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        {{ isUploading ? 'Uploading...' : 'Upload' }}
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="mt-4 flex items-center gap-2 p-3 bg-cyber-pink/10 border border-cyber-pink/30 rounded-lg">
      <svg class="w-5 h-5 text-cyber-pink flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-cyber-pink">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDocumentsStore } from '../stores/documentsStore'
import { FileType } from '../types/document.types'

const emit = defineEmits<{
  (e: 'upload-success'): void
  (e: 'upload-error', error: string): void
}>()

const documentsStore = useDocumentsStore()

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFile = ref<File | null>(null)
const title = ref('')
const isDragOver = ref(false)
const isUploading = ref(false)
const error = ref<string | null>(null)

const acceptedFormats = '.pdf,.docx,.pptx,.txt,.md'
const maxFileSize = 50 * 1024 * 1024 // 50MB

function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    validateAndSetFile(target.files[0])
  }
}

function handleDrop(event: DragEvent) {
  isDragOver.value = false
  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    validateAndSetFile(event.dataTransfer.files[0])
  }
}

function validateAndSetFile(file: File) {
  error.value = null

  // Check file size
  if (file.size > maxFileSize) {
    error.value = 'File size exceeds 50MB limit'
    return
  }

  // Check file type
  const extension = file.name.split('.').pop()?.toLowerCase()
  const validExtensions = Object.values(FileType)
  if (!extension || !validExtensions.includes(extension as any)) {
    error.value = `Invalid file type. Supported formats: ${validExtensions.join(', ').toUpperCase()}`
    return
  }

  selectedFile.value = file
  title.value = file.name.replace(/\.[^/.]+$/, '') // Remove extension
}

function clearFile() {
  selectedFile.value = null
  title.value = ''
  error.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

async function handleUpload() {
  if (!selectedFile.value) return

  isUploading.value = true
  error.value = null

  try {
    await documentsStore.uploadDocument({
      file: selectedFile.value,
      title: title.value || selectedFile.value.name,
    })

    emit('upload-success')
    clearFile()
  } catch (err: any) {
    const errorMessage = err.message || 'Failed to upload document'
    error.value = errorMessage
    emit('upload-error', errorMessage)
  } finally {
    isUploading.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>
