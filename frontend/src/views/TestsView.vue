<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-text-primary mb-2">Tests</h1>
        <p class="text-text-secondary">Generate and manage your test questions</p>
      </div>
      <button class="btn-neon">
        <svg class="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Create Test
      </button>
    </div>

    <!-- Empty State -->
    <div class="card-cyber text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-cyber-blue/20 flex items-center justify-center">
        <svg class="w-8 h-8 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-text-primary mb-2">No tests yet</h3>
      <p class="text-text-muted mb-6">Create your first test from uploaded documents</p>
      <button class="btn-neon">
        Generate Test
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Tests logic will be implemented later
</script>
