<template>
  <div class="min-h-screen bg-bg text-body flex flex-col font-sans">
    <header class="bg-surface border-b border-border-subtle shadow-sm relative z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 rounded-full bg-accent flex items-center justify-center font-bold text-bg shadow-lg border border-accent/80">N</div>
          <h1 class="text-xl font-bold text-heading">Nucleus</h1>
        </div>
        <nav class="flex space-x-4">
          <router-link to="/" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner">Targets</router-link>
          <router-link to="/scans" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner">Scans History</router-link>
        </nav>
      </div>
    </header>
    
    <main class="flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative">
      <!-- Glow effect -->
      <div class="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-3xl h-64 bg-accent/10 blur-[120px] pointer-events-none"></div>
      <router-view class="relative z-10"></router-view>
    </main>
  </div>
</template>
