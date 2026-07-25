<template>
  <div class="min-h-screen bg-bg text-body flex flex-col font-sans">
    <header class="bg-surface border-b border-border-subtle shadow-sm relative z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <img src="/favicon.svg" alt="Nucleus Logo" class="w-8 h-8 drop-shadow-md" />
          <h1 class="text-xl font-bold text-heading">Nucleus</h1>
        </div>
        <nav class="flex space-x-4 items-center" v-if="$route.path !== '/login'">
          <router-link to="/" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner" exact>Dashboard</router-link>
          <router-link to="/targets" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner">Targets</router-link>
          <router-link to="/scans" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner">Scans History</router-link>
          <router-link to="/settings" class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading" active-class="bg-bg text-accent shadow-inner">Settings</router-link>
          <button @click="logout" class="ml-4 px-3 py-1.5 border border-border-subtle rounded-md text-sm font-medium text-body hover:text-red-400 hover:border-red-400/50 transition-colors hover:cursor-pointer">Logout</button>
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

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const logout = async () => {
  await fetch('/api/auth/logout', { method: 'POST' })
  router.push('/login')
}
</script>
