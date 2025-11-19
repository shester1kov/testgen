<template>
  <div class="min-h-screen bg-dark-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <!-- Glow orbs -->
      <div class="absolute top-1/4 -left-20 w-96 h-96 bg-neon-orange opacity-10 rounded-full blur-3xl animate-pulse-slow"></div>
      <div class="absolute bottom-1/4 -right-20 w-96 h-96 bg-cyber-blue opacity-10 rounded-full blur-3xl animate-pulse-slow" style="animation-delay: 1s;"></div>

      <!-- Cyber grid -->
      <div class="absolute inset-0 opacity-5" style="background-image: linear-gradient(#ff6b35 1px, transparent 1px), linear-gradient(90deg, #ff6b35 1px, transparent 1px); background-size: 50px 50px;"></div>
    </div>

    <!-- Content -->
    <div class="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
      <!-- Logo with glow -->
      <div class="text-center mb-8">
        <h1 class="text-5xl font-bold text-glow mb-2 animate-glow">
          TestGen
        </h1>
        <div class="flex items-center justify-center space-x-2 text-text-secondary">
          <span class="text-sm uppercase tracking-wider">AI-Powered Test Generation</span>
          <div class="w-2 h-2 rounded-full bg-neon-orange animate-pulse"></div>
        </div>
      </div>

      <!-- Auth card -->
      <div class="card-cyber backdrop-blur-sm bg-dark-700/90">
        <router-view />
      </div>

      <!-- Footer -->
      <p class="mt-6 text-center text-xs text-text-muted">
        Powered by AI • Built with Vue 3 & Tailwind CSS
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
// Auth layout for login/register
</script>
