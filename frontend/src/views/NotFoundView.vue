<template>
  <div class="text-center py-16">
    <div class="card-cyber max-w-lg mx-auto">
      <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-cyber-pink/20 flex items-center justify-center animate-pulse">
        <svg class="w-10 h-10 text-cyber-pink" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
      </div>
      <h1 class="text-6xl font-bold text-glow mb-4">404</h1>
      <p class="text-xl text-text-primary mb-2">Page not found</p>
      <p class="text-text-muted mb-8">The page you're looking for doesn't exist or has been moved.</p>
      <router-link
        to="/dashboard"
        class="btn-neon inline-block"
      >
        <svg class="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        Go to Dashboard
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
// 404 page
</script>
