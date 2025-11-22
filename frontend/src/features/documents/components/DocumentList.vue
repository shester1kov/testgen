<template>
  <div class="document-list">
    <!-- Header -->
    <div class="list-header">
      <h2 class="list-title">My Documents</h2>
      <p class="list-subtitle">{{ total }} document{{ total !== 1 ? 's' : '' }} total</p>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading && documents.length === 0" class="loading-state">
      <svg class="loading-spinner" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
      <p class="loading-text">Loading documents...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="documents.length === 0" class="empty-state">
      <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="empty-title">No documents yet</h3>
      <p class="empty-subtitle">Upload your first document to get started</p>
    </div>

    <!-- Document grid -->
    <div v-else class="documents-grid">
      <DocumentCard v-for="document in documents" :key="document.id" :document="document"
        @view-text="handleViewText" />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button type="button" class="pagination-button" :disabled="currentPage === 1" @click="handlePageChange(currentPage - 1)">
        <svg class="pagination-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Previous
      </button>

      <div class="pagination-info">
        Page {{ currentPage }} of {{ totalPages }}
      </div>

      <button type="button" class="pagination-button" :disabled="currentPage === totalPages"
        @click="handlePageChange(currentPage + 1)">
        Next
        <svg class="pagination-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="error-message">
      <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="error-text">{{ error }}</p>
      <button type="button" class="error-retry" @click="handleRetry">
        Retry
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useDocumentsStore } from '../stores/documentsStore'
import DocumentCard from './DocumentCard.vue'
import type { Document } from '../types/document.types'

const emit = defineEmits<{
  (e: 'view-document', document: Document): void
}>()

const documentsStore = useDocumentsStore()
const { documents, total, currentPage, totalPages, loading: isLoading, error } = storeToRefs(documentsStore)

onMounted(() => {
  loadDocuments()
})

async function loadDocuments() {
  try {
    await documentsStore.fetchDocuments(currentPage.value)
  } catch (err) {
    // Error is handled by store
  }
}

function handlePageChange(page: number) {
  documentsStore.fetchDocuments(page)
}

function handleViewText(document: Document) {
  emit('view-document', document)
}

function handleRetry() {
  documentsStore.clearError()
  loadDocuments()
}
</script>

<style scoped>
.document-list {
  @apply space-y-6;
}

.list-header {
  @apply border-b border-gray-200 pb-4;
}

.list-title {
  @apply text-2xl font-bold text-gray-900;
}

.list-subtitle {
  @apply text-sm text-gray-600 mt-1;
}

.loading-state {
  @apply flex flex-col items-center justify-center py-12;
}

.loading-spinner {
  @apply w-12 h-12 text-blue-600 animate-spin mb-4;
}

.loading-text {
  @apply text-gray-600;
}

.empty-state {
  @apply flex flex-col items-center justify-center py-12 text-center;
}

.empty-icon {
  @apply w-16 h-16 text-gray-400 mb-4;
}

.empty-title {
  @apply text-xl font-semibold text-gray-900 mb-2;
}

.empty-subtitle {
  @apply text-gray-600;
}

.documents-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-4;
}

.pagination {
  @apply flex items-center justify-between border-t border-gray-200 pt-4;
}

.pagination-button {
  @apply px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700;
  @apply hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
  @apply flex items-center gap-2;
}

.pagination-icon {
  @apply w-5 h-5;
}

.pagination-info {
  @apply text-sm text-gray-700;
}

.error-message {
  @apply flex items-center justify-between gap-4 p-4 bg-red-50 border border-red-200 rounded-md;
}

.error-icon {
  @apply w-5 h-5 text-red-500 flex-shrink-0;
}

.error-text {
  @apply text-sm text-red-700 flex-1;
}

.error-retry {
  @apply px-3 py-1 bg-red-100 text-red-700 rounded-md text-sm font-medium;
  @apply hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500;
}
</style>
