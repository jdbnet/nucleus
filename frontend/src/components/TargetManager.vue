<template>
  <div class="space-y-6">
    <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle backdrop-blur-sm">
      <h2 class="text-xl font-semibold mb-4 text-heading">Add New Target</h2>
      <form @submit.prevent="addTarget" class="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
        <div>
          <label class="block text-sm font-medium text-body mb-1">Name</label>
          <input v-model="form.name" type="text" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all" placeholder="Prod App">
        </div>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Address / Subnet</label>
          <input v-model="form.address" type="text" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all" placeholder="10.0.0.1/24 or example.com">
        </div>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Schedule</label>
          <div class="flex space-x-2">
            <select v-model="form.schedule" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
              <option value="manual">Manual Only</option>
              <option value="@midnight">Daily at Midnight</option>
              <option value="@hourly">Hourly</option>
              <option value="0 0 * * 0">Weekly on Sunday</option>
              <option value="custom">Custom Cron...</option>
            </select>
            <input v-if="form.schedule === 'custom'" v-model="form.customSchedule" type="text" placeholder="0 2 * * 1" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
        </div>
        <div>
          <button type="submit" class="w-full bg-accent hover:bg-accent/80 text-bg font-bold py-2 px-4 rounded-lg shadow-md hover:shadow-lg transition-all transform hover:-translate-y-0.5">
            Add Target
          </button>
        </div>
      </form>
    </div>

    <div class="bg-surface rounded-xl overflow-hidden shadow-lg border border-border-subtle">
      <div class="px-6 py-4 border-b border-border-subtle flex justify-between items-center">
        <h2 class="text-xl font-semibold text-heading">Configured Targets</h2>
        <button @click="loadTargets" class="text-body hover:text-heading transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
        </button>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-bg/50 text-body">
            <tr>
              <th class="px-6 py-3 font-medium">Name</th>
              <th class="px-6 py-3 font-medium">Address</th>
              <th class="px-6 py-3 font-medium">Schedule</th>
              <th class="px-6 py-3 font-medium">Last Scan</th>
              <th class="px-6 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-subtle">
            <tr v-for="target in targets" :key="target.id" class="hover:bg-border-subtle/30 transition-colors">
              <td class="px-6 py-4 text-heading font-medium">{{ target.name }}</td>
              <td class="px-6 py-4 text-body font-mono text-xs">{{ target.address }}</td>
              <td class="px-6 py-4 text-body">
                <span class="px-2 py-1 bg-bg rounded text-xs border border-border-subtle">{{ target.schedule }}</span>
              </td>
              <td class="px-6 py-4 text-body">
                {{ target.last_scan_at ? new Date(target.last_scan_at).toLocaleString() : 'Never' }}
              </td>
              <td class="px-6 py-4 text-right space-x-3">
                <button @click="runScan(target.id)" class="text-accent hover:text-accent/80 font-medium transition-colors">Run Now</button>
                <button @click="deleteTarget(target.id)" class="text-red-400 hover:text-red-300 font-medium transition-colors">Delete</button>
              </td>
            </tr>
            <tr v-if="targets.length === 0">
              <td colspan="5" class="px-6 py-8 text-center text-body">No targets configured. Add one above.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const targets = ref([])
const form = ref({ name: '', address: '', schedule: 'manual', customSchedule: '' })

const loadTargets = async () => {
  try {
    const res = await fetch('/api/targets')
    targets.value = await res.json() || []
  } catch (e) {
    console.error(e)
  }
}

const addTarget = async () => {
  let finalSchedule = form.value.schedule
  if (finalSchedule === 'custom') {
    finalSchedule = form.value.customSchedule.trim()
    if (!finalSchedule) {
      alert('Please enter a valid custom cron expression')
      return
    }
  }

  await fetch('/api/targets', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: form.value.name,
      address: form.value.address,
      schedule: finalSchedule
    })
  })
  form.value = { name: '', address: '', schedule: 'manual', customSchedule: '' }
  loadTargets()
}

const deleteTarget = async (id) => {
  if (!confirm('Delete target?')) return
  await fetch(`/api/targets/${id}`, { method: 'DELETE' })
  loadTargets()
}

const runScan = async (id) => {
  await fetch(`/api/targets/${id}/scan`, { method: 'POST' })
  alert('Scan triggered successfully. Check Scans History.')
  loadTargets()
}

onMounted(loadTargets)
</script>
