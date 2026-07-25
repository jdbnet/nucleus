<template>
  <div class="flex items-center justify-center min-h-[70vh]">
    <div class="w-full max-w-md bg-surface p-8 rounded-xl shadow-lg border border-border-subtle backdrop-blur-sm">
      <div class="flex justify-center mb-6">
        <img src="/favicon.svg" alt="Nucleus Logo" class="w-16 h-16 drop-shadow-md" />
      </div>
      <h2 class="text-2xl font-bold text-center text-heading mb-6">Sign In to Nucleus</h2>
      
      <form @submit.prevent="login" class="space-y-4">
        <div v-if="error" class="p-3 bg-red-500/20 border border-red-500/50 text-red-400 rounded-lg text-sm text-center">
          {{ error }}
        </div>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Username</label>
          <input v-model="username" type="text" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
        </div>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Password</label>
          <input v-model="password" type="password" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
        </div>
        <button type="submit" class="w-full bg-accent hover:bg-accent/80 text-bg font-bold py-2 px-4 rounded-lg shadow-md hover:shadow-lg transition-all hover:cursor-pointer mt-2">
          Login
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const username = ref('')
const password = ref('')
const error = ref('')

const isSafeRedirect = (path) => {
  return typeof path === 'string' && path.startsWith('/') && !path.startsWith('//')
}

const login = async () => {
  error.value = ''
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, password: password.value })
    })
    
    if (res.ok) {
      const redirect = route.query.redirect
      router.push(isSafeRedirect(redirect) ? redirect : '/')
    } else {
      error.value = 'Invalid username or password'
    }
  } catch (e) {
    error.value = 'Network error'
  }
}
</script>
