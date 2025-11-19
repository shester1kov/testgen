<template>
  <div>
    <div class="mb-8">
      <router-link
        to="/tests"
        class="inline-flex items-center text-text-secondary hover:text-neon-orange transition-colors mb-4"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Back to Tests
      </router-link>
      <h1 class="text-3xl font-bold text-text-primary mb-2">Create Test</h1>
      <p class="text-text-secondary">Generate test questions from your documents using AI</p>
    </div>

    <div class="card-cyber">
      <p class="text-text-secondary">Test creation form will be displayed here</p>
    </div>
  </div>
</template>

<script setup lang="ts">
// Create test logic will be implemented later
</script>
