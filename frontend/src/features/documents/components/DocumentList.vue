<template>
  <div class="card-cyber">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-lg font-semibold text-text-primary">My Documents</h2>
        <p class="text-sm text-text-muted mt-1">{{ total }} document{{ total !== 1 ? 's' : '' }} total</p>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading && documents.length === 0" class="flex flex-col items-center justify-center py-12">
      <div class="w-12 h-12 border-4 border-neon-orange/30 border-t-neon-orange rounded-full animate-spin mb-4"></div>
      <p class="text-text-muted">Loading documents...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="documents.length === 0" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-dark-600 flex items-center justify-center">
        <svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-text-primary mb-2">No documents yet</h3>
      <p class="text-text-muted">Upload your first document to get started</p>
    </div>

    <!-- Document grid -->
    <div v-else class="space-y-3">
      <DocumentCard
        v-for="document in documents"
        :key="document.id"
        :document="document"
        @view-text="handleViewText"
      />
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between mt-6 pt-4 border-t border-dark-500">
      <button
        type="button"
        class="px-4 py-2 border border-dark-500 rounded-lg text-sm font-medium text-text-secondary hover:bg-dark-600 hover:border-neon-orange transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Previous
      </button>

      <div class="text-sm text-text-muted">
        Page <span class="text-neon-orange font-medium">{{ currentPage }}</span> of {{ totalPages }}
      </div>

      <button
        type="button"
        class="px-4 py-2 border border-dark-500 rounded-lg text-sm font-medium text-text-secondary hover:bg-dark-600 hover:border-neon-orange transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        :disabled="currentPage === totalPages"
        @click="handlePageChange(currentPage + 1)"
      >
        Next
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- Error message -->
    <div v-if="error" class="mt-4 flex items-center justify-between gap-4 p-4 bg-cyber-pink/10 border border-cyber-pink/30 rounded-lg">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-cyber-pink flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-sm text-cyber-pink flex-1">{{ error }}</p>
      </div>
      <button
        type="button"
        class="px-3 py-1.5 bg-cyber-pink/20 text-cyber-pink rounded-lg text-sm font-medium hover:bg-cyber-pink/30 transition-colors"
        @click="handleRetry"
      >
        Retry
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
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
