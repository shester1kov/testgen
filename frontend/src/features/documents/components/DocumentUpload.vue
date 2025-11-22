<template>
  <div class="document-upload">
    <div class="upload-area" :class="{ 'drag-over': isDragOver }" @drop.prevent="handleDrop"
      @dragover.prevent="isDragOver = true" @dragleave.prevent="isDragOver = false">
      <input ref="fileInput" type="file" :accept="acceptedFormats" class="hidden" @change="handleFileSelect" />

      <div class="upload-content">
        <svg class="upload-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>

        <p class="upload-title">
          {{ isDragOver ? 'Drop file here' : 'Upload a document' }}
        </p>

        <p class="upload-subtitle">
          Drag and drop or
          <button type="button" class="browse-button" @click="$refs.fileInput.click()">
            browse
          </button>
        </p>

        <p class="upload-formats">Supported: PDF, DOCX, PPTX, TXT, MD (max 50MB)</p>
      </div>
    </div>

    <!-- Selected file info -->
    <div v-if="selectedFile" class="selected-file">
      <div class="file-info">
        <svg class="file-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <div class="file-details">
          <p class="file-name">{{ selectedFile.name }}</p>
          <p class="file-size">{{ formatFileSize(selectedFile.size) }}</p>
        </div>
      </div>

      <button type="button" class="remove-button" @click="clearFile">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Title input -->
    <div v-if="selectedFile" class="title-input">
      <label for="doc-title" class="title-label">Document Title (optional)</label>
      <input id="doc-title" v-model="title" type="text" class="title-field"
        :placeholder="selectedFile.name" />
    </div>

    <!-- Upload button -->
    <div v-if="selectedFile" class="upload-actions">
      <button type="button" class="cancel-button" :disabled="isUploading" @click="clearFile">
        Cancel
      </button>
      <button type="button" class="upload-button" :disabled="isUploading" @click="handleUpload">
        <svg v-if="isUploading" class="spinner" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        {{ isUploading ? 'Uploading...' : 'Upload' }}
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="error-message">
      <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      {{ error }}
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
    error.value = err.message || 'Failed to upload document'
    emit('upload-error', error.value)
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

<style scoped>
.document-upload {
  @apply space-y-4;
}

.upload-area {
  @apply border-2 border-dashed border-gray-300 rounded-lg p-8 text-center cursor-pointer transition-colors;
  @apply hover:border-gray-400;
}

.upload-area.drag-over {
  @apply border-blue-500 bg-blue-50;
}

.upload-content {
  @apply space-y-3;
}

.upload-icon {
  @apply w-12 h-12 mx-auto text-gray-400;
}

.upload-title {
  @apply text-lg font-medium text-gray-700;
}

.upload-subtitle {
  @apply text-sm text-gray-600;
}

.browse-button {
  @apply text-blue-600 hover:text-blue-700 font-medium underline;
}

.upload-formats {
  @apply text-xs text-gray-500;
}

.hidden {
  @apply sr-only;
}

.selected-file {
  @apply flex items-center justify-between p-4 bg-gray-50 rounded-lg border border-gray-200;
}

.file-info {
  @apply flex items-center gap-3;
}

.file-icon {
  @apply w-8 h-8 text-blue-600;
}

.file-details {
  @apply text-left;
}

.file-name {
  @apply text-sm font-medium text-gray-900;
}

.file-size {
  @apply text-xs text-gray-500;
}

.remove-button {
  @apply p-1 text-gray-400 hover:text-gray-600 transition-colors;
}

.remove-button svg {
  @apply w-5 h-5;
}

.title-input {
  @apply space-y-2;
}

.title-label {
  @apply block text-sm font-medium text-gray-700;
}

.title-field {
  @apply w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm;
  @apply focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
}

.upload-actions {
  @apply flex gap-3 justify-end;
}

.cancel-button {
  @apply px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700;
  @apply hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.upload-button {
  @apply px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium flex items-center gap-2;
  @apply hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.spinner {
  @apply w-4 h-4 animate-spin;
}

.error-message {
  @apply flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-700;
}

.error-icon {
  @apply w-5 h-5 text-red-500;
}
</style>
