<template>
  <div class="min-h-screen bg-dark-900 flex flex-col relative overflow-hidden">
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <!-- Glow orbs -->
      <div class="absolute top-1/4 -left-20 w-96 h-96 bg-neon-orange opacity-10 rounded-full blur-3xl animate-pulse-slow"></div>
      <div class="absolute bottom-1/4 -right-20 w-96 h-96 bg-cyber-blue opacity-10 rounded-full blur-3xl animate-pulse-slow" style="animation-delay: 1s;"></div>

      <!-- Cyber grid -->
      <div class="absolute inset-0 opacity-5" style="background-image: linear-gradient(#ff6b35 1px, transparent 1px), linear-gradient(90deg, #ff6b35 1px, transparent 1px); background-size: 50px 50px;"></div>
    </div>

    <!-- Header -->
    <header class="relative z-10 w-full px-6 py-6 flex justify-between items-center max-w-7xl mx-auto">
      <div class="flex items-center">
        <span class="text-2xl font-bold text-glow">TestGen</span>
        <span class="ml-2 text-xs text-text-muted uppercase tracking-wider hidden sm:inline-block">AI Powered</span>
      </div>
      
      <div class="flex space-x-4">
        <router-link
          to="/login"
          class="text-text-secondary hover:text-neon-orange px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300 hover:bg-dark-700 border border-transparent hover:border-dark-500"
        >
          Вход
        </router-link>
        <router-link
          to="/register"
          class="btn-primary text-sm px-6 py-2"
        >
          Регистрация
        </router-link>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-grow flex flex-col justify-center items-center relative z-10 px-4 sm:px-6 lg:px-8 text-center max-w-5xl mx-auto">
      <div class="mb-8">
        <h1 class="text-5xl sm:text-7xl font-bold mb-6 leading-tight text-neon-orange drop-shadow-[0_0_8px_rgba(255,107,53,0.5)]">
          Автоматическая генерация <br />
          учебных тестов
        </h1>
        <p class="mt-4 max-w-2xl mx-auto text-xl text-text-secondary">
          Используйте возможности искусственного интеллекта для мгновенного создания качественных тестовых заданий на основе ваших учебных материалов.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-12">
        <!-- Feature 1 -->
        <div class="card-cyber backdrop-blur-sm bg-dark-700/80 p-6 flex flex-col items-center">
          <div class="w-12 h-12 rounded-lg bg-dark-800 border border-dark-500 flex items-center justify-center mb-4 shadow-neon-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neon-orange" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-text-primary mb-2">Мгновенный результат</h3>
          <p class="text-text-secondary text-sm">Генерация тестов занимает секунды вместо часов ручной работы.</p>
        </div>

        <!-- Feature 2 -->
        <div class="card-cyber backdrop-blur-sm bg-dark-700/80 p-6 flex flex-col items-center">
          <div class="w-12 h-12 rounded-lg bg-dark-800 border border-dark-500 flex items-center justify-center mb-4 shadow-cyber-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-cyber-blue" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-text-primary mb-2">На основе документов</h3>
          <p class="text-text-secondary text-sm">Загружайте PDF или DOCX файлы, и ИИ создаст вопросы прямо по тексту.</p>
        </div>

        <!-- Feature 3 -->
        <div class="card-cyber backdrop-blur-sm bg-dark-700/80 p-6 flex flex-col items-center">
          <div class="w-12 h-12 rounded-lg bg-dark-800 border border-dark-500 flex items-center justify-center mb-4 shadow-neon-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-neon-orange" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-text-primary mb-2">Разные форматы</h3>
          <p class="text-text-secondary text-sm">Поддержка множественного выбора, верно/неверно и открытых вопросов.</p>
        </div>
      </div>

      <div class="mt-16 mb-8">
        <router-link
          to="/register"
          class="btn-primary text-lg px-8 py-4 shadow-neon animate-pulse-slow inline-flex items-center"
        >
          Начать использование
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </router-link>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 w-full py-6 text-center text-sm text-text-muted border-t border-dark-700/50">
      <p>TestGen &copy; {{ new Date().getFullYear() }} • Powered by AI</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
// Welcome page component (Landing)
</script>
