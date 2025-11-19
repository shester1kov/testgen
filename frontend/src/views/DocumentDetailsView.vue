<template>
  <div>
    <div class="mb-8">
      <router-link
        to="/documents"
        class="inline-flex items-center text-text-secondary hover:text-neon-orange transition-colors mb-4"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Back to Documents
      </router-link>
      <h1 class="text-3xl font-bold text-text-primary mb-2">Document Details</h1>
      <p class="text-text-muted">Document ID: <span class="text-text-secondary font-mono">{{ id }}</span></p>
    </div>

    <div class="card-cyber">
      <p class="text-text-secondary">Document details will be displayed here</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const id = route.params.id
</script>
