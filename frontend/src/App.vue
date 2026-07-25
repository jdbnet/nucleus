<template>
  <div class="min-h-screen bg-bg text-body flex flex-col font-sans">
    <header class="bg-surface border-b border-border-subtle shadow-sm relative z-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <img src="/favicon.svg" alt="Nucleus Logo" class="w-8 h-8 drop-shadow-md" />
          <h1 class="text-xl font-bold text-heading">Nucleus</h1>
        </div>

        <template v-if="$route.path !== '/login'">
          <!-- Desktop nav -->
          <nav class="hidden md:flex space-x-4 items-center">
            <router-link
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="px-3 py-2 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading"
              :active-class="link.exact ? '' : 'bg-bg text-accent shadow-inner'"
              :exact-active-class="link.exact ? 'bg-bg text-accent shadow-inner' : ''"
            >{{ link.label }}</router-link>
            <button @click="logout" class="ml-4 px-3 py-1.5 border border-border-subtle rounded-md text-sm font-medium text-body hover:text-red-400 hover:border-red-400/50 transition-colors hover:cursor-pointer">Logout</button>
          </nav>

          <!-- Mobile hamburger -->
          <button
            type="button"
            class="md:hidden p-2 rounded-md text-body hover:text-heading hover:bg-border-subtle transition-colors hover:cursor-pointer"
            :aria-expanded="menuOpen"
            aria-label="Toggle navigation menu"
            @click="menuOpen = !menuOpen"
          >
            <svg v-if="!menuOpen" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </template>
      </div>

      <!-- Mobile menu -->
      <nav
        v-if="$route.path !== '/login' && menuOpen"
        class="md:hidden border-t border-border-subtle bg-surface"
      >
        <div class="max-w-7xl mx-auto px-4 py-3 flex flex-col gap-1">
          <router-link
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="px-3 py-2.5 rounded-md text-sm font-medium transition-colors hover:bg-border-subtle hover:text-heading"
            :active-class="link.exact ? '' : 'bg-bg text-accent shadow-inner'"
            :exact-active-class="link.exact ? 'bg-bg text-accent shadow-inner' : ''"
            @click="menuOpen = false"
          >{{ link.label }}</router-link>
          <button
            @click="logout"
            class="mt-1 px-3 py-2.5 border border-border-subtle rounded-md text-sm font-medium text-body hover:text-red-400 hover:border-red-400/50 transition-colors hover:cursor-pointer text-left"
          >Logout</button>
        </div>
      </nav>
    </header>
    
    <main class="flex-1 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative">
      <!-- Glow effect -->
      <div class="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-3xl h-64 bg-accent/10 blur-[120px] pointer-events-none"></div>
      <router-view class="relative z-10"></router-view>
    </main>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const menuOpen = ref(false)

const navLinks = [
  { to: '/', label: 'Dashboard', exact: true },
  { to: '/targets', label: 'Targets' },
  { to: '/scans', label: 'Scans History' },
  { to: '/settings', label: 'Settings' },
]

watch(() => route.path, () => {
  menuOpen.value = false
})

const logout = async () => {
  menuOpen.value = false
  await fetch('/api/auth/logout', { method: 'POST' })
  router.push('/login')
}
</script>
