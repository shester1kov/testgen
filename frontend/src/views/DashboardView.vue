<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-text-primary mb-2">Dashboard</h1>
      <p class="text-text-secondary">Welcome to TestGen - AI-Powered Test Generation</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">Total Documents</p>
            <p class="text-3xl font-bold text-text-primary mt-2">0</p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-neon-orange/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">Generated Tests</p>
            <p class="text-3xl font-bold text-text-primary mt-2">0</p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-cyber-blue/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card-cyber">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-text-muted text-sm uppercase tracking-wider">Total Questions</p>
            <p class="text-3xl font-bold text-text-primary mt-2">0</p>
          </div>
          <div class="w-12 h-12 rounded-lg bg-cyber-purple/20 flex items-center justify-center">
            <svg class="w-6 h-6 text-cyber-purple" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card-cyber">
      <h2 class="text-xl font-bold text-text-primary mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <router-link
          to="/documents"
          class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 hover:border-neon-orange transition-all duration-300 group"
        >
          <div class="w-10 h-10 rounded-lg bg-neon-orange/20 flex items-center justify-center mr-4 group-hover:bg-neon-orange/30">
            <svg class="w-5 h-5 text-neon-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary group-hover:text-neon-orange transition-colors">Upload Document</p>
            <p class="text-sm text-text-muted">Add new learning material</p>
          </div>
        </router-link>

        <router-link
          to="/tests"
          class="flex items-center p-4 bg-dark-600 rounded-lg border border-dark-500 hover:border-neon-orange transition-all duration-300 group"
        >
          <div class="w-10 h-10 rounded-lg bg-cyber-blue/20 flex items-center justify-center mr-4 group-hover:bg-cyber-blue/30">
            <svg class="w-5 h-5 text-cyber-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <p class="font-semibold text-text-primary group-hover:text-neon-orange transition-colors">Create Test</p>
            <p class="text-sm text-text-muted">Generate new test questions</p>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Dashboard logic will be implemented later
</script>
