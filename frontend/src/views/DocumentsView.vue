<template>
  <div class="documents-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">Documents</h1>
        <p class="page-subtitle">Upload and manage your learning materials</p>
      </div>
    </div>

    <div class="content-grid">
      <!-- Upload Section -->
      <div class="upload-section">
        <h2 class="section-title">Upload New Document</h2>
        <DocumentUpload @upload-success="handleUploadSuccess" />
      </div>

      <!-- Documents List -->
      <div class="documents-section">
        <DocumentList @view-document="handleViewDocument" />
      </div>
    </div>

    <!-- Document Text Preview Modal -->
    <Teleport to="body">
      <div v-if="showPreview" class="modal-overlay" @click="closePreview">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3 class="modal-title">{{ selectedDocument?.title }}</h3>
            <button type="button" class="modal-close" @click="closePreview">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div v-if="selectedDocument?.parsed_text" class="parsed-text">
              {{ selectedDocument.parsed_text }}
            </div>
            <div v-else class="no-text">
              <p>This document has not been parsed yet.</p>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDocumentsStore } from '@/features/documents/stores/documentsStore'
import DocumentUpload from '@/features/documents/components/DocumentUpload.vue'
import DocumentList from '@/features/documents/components/DocumentList.vue'
import type { Document } from '@/features/documents/types/document.types'

const documentsStore = useDocumentsStore()
const showPreview = ref(false)
const selectedDocument = ref<Document | null>(null)

function handleUploadSuccess() {
  // Reload the documents list after successful upload
  documentsStore.fetchDocuments(1)
}

function handleViewDocument(document: Document) {
  selectedDocument.value = document
  showPreview.value = true
}

function closePreview() {
  showPreview.value = false
  selectedDocument.value = null
}
</script>

<style scoped>
.documents-page {
  @apply max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8;
}

.page-header {
  @apply mb-8;
}

.page-title {
  @apply text-3xl font-bold text-gray-900 mb-2;
}

.page-subtitle {
  @apply text-gray-600;
}

.content-grid {
  @apply grid grid-cols-1 lg:grid-cols-3 gap-8;
}

.upload-section {
  @apply lg:col-span-1;
}

.documents-section {
  @apply lg:col-span-2;
}

.section-title {
  @apply text-xl font-semibold text-gray-900 mb-4;
}

/* Modal styles */
.modal-overlay {
  @apply fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4;
}

.modal-content {
  @apply bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] flex flex-col;
}

.modal-header {
  @apply flex items-center justify-between p-6 border-b border-gray-200;
}

.modal-title {
  @apply text-xl font-semibold text-gray-900;
}

.modal-close {
  @apply text-gray-400 hover:text-gray-600 transition-colors;
}

.modal-body {
  @apply p-6 overflow-y-auto flex-1;
}

.parsed-text {
  @apply whitespace-pre-wrap text-gray-700 leading-relaxed;
}

.no-text {
  @apply text-center py-12 text-gray-500;
}
</style>
