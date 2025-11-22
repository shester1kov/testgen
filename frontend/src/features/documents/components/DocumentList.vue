<template>
  <div class="flex flex-col gap-6">
    <!-- Header -->
    <div class="border-b border-gray-200 pb-4">
      <h2 class="text-2xl font-bold text-gray-900">My Documents</h2>
      <p class="text-sm text-gray-600 mt-1">{{ total }} document{{ total !== 1 ? 's' : '' }} total</p>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading && documents.length === 0" class="flex flex-col items-center justify-center py-12">
      <svg class="w-12 h-12 text-blue-600 animate-spin mb-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
      <p class="text-gray-600">Loading documents...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="documents.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
      <svg class="w-16 h-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="text-xl font-semibold text-gray-900 mb-2">No documents yet</h3>
      <p class="text-gray-600">Upload your first document to get started</p>
    </div>

    <!-- Document grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <DocumentCard v-for="document in documents" :key="document.id" :document="document"
        @view-text="handleViewText" />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-200 pt-4">
      <button type="button"
              class="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="currentPage === 1"
              @click="handlePageChange(currentPage - 1)">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Previous
      </button>

      <div class="text-sm text-gray-700">
        Page {{ currentPage }} of {{ totalPages }}
      </div>

      <button type="button"
              class="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="currentPage === totalPages"
              @click="handlePageChange(currentPage + 1)">
        Next
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="flex items-center justify-between gap-4 p-4 bg-red-50 border border-red-200 rounded-md">
      <svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-red-700 flex-1">{{ error }}</p>
      <button type="button"
              class="px-3 py-1 bg-red-100 text-red-700 rounded-md text-sm font-medium hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              @click="handleRetry">
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
