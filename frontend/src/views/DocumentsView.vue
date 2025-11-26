<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">Документы</h1>
        <p class="text-text-secondary">Загружайте и управляйте учебными материалами</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Upload Section -->
      <div class="lg:col-span-1">
        <DocumentUpload @upload-success="handleUploadSuccess" />
      </div>

      <!-- Documents List -->
      <div class="lg:col-span-2">
        <DocumentList @view-document="handleViewDocument" />
      </div>
    </div>

    <!-- Document Text Preview Modal -->
    <Teleport to="body">
      <div v-if="showPreview" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click="closePreview">
        <div class="card-cyber max-w-4xl w-full max-h-[90vh] flex flex-col" @click.stop>
          <div class="flex items-center justify-between p-6 border-b border-dark-500">
            <h3 class="text-xl font-semibold text-text-primary">{{ selectedDocument?.title }}</h3>
            <button type="button" class="text-text-muted hover:text-neon-orange transition-colors" @click="closePreview">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="p-6 overflow-y-auto flex-1">
            <div v-if="selectedDocument?.parsed_text" class="whitespace-pre-wrap text-text-secondary leading-relaxed font-mono text-sm">
              {{ selectedDocument.parsed_text }}
            </div>
            <div v-else class="text-center py-12 text-text-muted">
              <p>Этот документ ещё не был обработан.</p>
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
