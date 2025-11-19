<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">Documents</h1>
        <p class="text-text-secondary">Upload and manage your learning materials</p>
      </div>
      <button class="btn-neon">
        <svg class="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Upload Document
      </button>
    </div>

    <!-- Empty State -->
    <div class="card-cyber text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-neon-orange/20 flex items-center justify-center">
        <svg class="w-8 h-8 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-text-primary mb-2">No documents yet</h3>
      <p class="text-text-muted mb-6">Upload your first document to start generating tests</p>
      <button class="btn-neon">
        Get Started
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Documents logic will be implemented later
</script>
