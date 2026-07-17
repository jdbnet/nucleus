<template>
  <div class="bg-surface rounded-xl overflow-hidden shadow-lg border border-border-subtle">
    <div class="px-6 py-4 border-b border-border-subtle flex justify-between items-center">
      <h2 class="text-xl font-semibold text-heading">Scans History</h2>
      <button @click="loadScans" class="text-body hover:text-heading transition-colors hover:cursor-pointer">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
      </button>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-left text-sm">
        <thead class="bg-bg/50 text-body">
          <tr>
            <th class="px-6 py-3 font-medium">Target</th>
            <th class="px-6 py-3 font-medium">Status</th>
            <th class="px-6 py-3 font-medium">Started At</th>
            <th class="px-6 py-3 font-medium">Duration</th>
            <th class="px-6 py-3 font-medium">Findings</th>
            <th class="px-6 py-3 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border-subtle">
          <tr v-for="scan in scans" :key="scan.id" class="hover:bg-border-subtle/30 transition-colors">
            <td class="px-6 py-4 text-heading font-medium">{{ scan.target_name }}</td>
            <td class="px-6 py-4">
              <span v-if="scan.status === 'running'" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-900/50 text-blue-400 border border-blue-800">
                <svg class="animate-spin -ml-0.5 mr-1.5 h-3 w-3 text-blue-400" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Running
              </span>
              <span v-else-if="scan.status === 'completed'" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-900/50 text-green-400 border border-green-800">
                <svg class="mr-1.5 h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                Completed
              </span>
              <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-900/50 text-red-400 border border-red-800">
                Failed
              </span>
            </td>
            <td class="px-6 py-4 text-body">
              {{ new Date(scan.started_at).toLocaleString() }}
            </td>
            <td class="px-6 py-4 text-body">
              {{ getDuration(scan.started_at, scan.completed_at) }}
            </td>
            <td class="px-6 py-4">
              <div class="flex space-x-2">
                <span v-if="getSeverityCount(scan, 'critical') > 0" class="px-2 py-0.5 rounded text-xs font-bold bg-red-500/20 text-red-400 border border-red-500/50" title="Critical">{{ getSeverityCount(scan, 'critical') }}</span>
                <span v-if="getSeverityCount(scan, 'high') > 0" class="px-2 py-0.5 rounded text-xs font-bold bg-orange-500/20 text-orange-400 border border-orange-500/50" title="High">{{ getSeverityCount(scan, 'high') }}</span>
                <span v-if="getSeverityCount(scan, 'medium') > 0" class="px-2 py-0.5 rounded text-xs font-bold bg-yellow-500/20 text-yellow-400 border border-yellow-500/50" title="Medium">{{ getSeverityCount(scan, 'medium') }}</span>
                <span v-if="getSeverityCount(scan, 'low') > 0" class="px-2 py-0.5 rounded text-xs font-bold bg-blue-500/20 text-blue-400 border border-blue-500/50" title="Low">{{ getSeverityCount(scan, 'low') }}</span>
                <span v-if="getSeverityCount(scan, 'info') > 0" class="px-2 py-0.5 rounded text-xs font-bold bg-border-subtle/50 text-body border border-border-subtle" title="Info">{{ getSeverityCount(scan, 'info') }}</span>
                <span v-if="totalFindings(scan) === 0" class="text-body text-xs">None</span>
              </div>
            </td>
            <td class="px-6 py-4 text-right">
              <router-link :to="`/scans/${scan.id}`" class="text-accent hover:text-accent/80 font-medium transition-colors">View Details</router-link>
            </td>
          </tr>
          <tr v-if="scans.length === 0">
            <td colspan="6" class="px-6 py-8 text-center text-body">No scans recorded yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const scans = ref([])
let interval = null

const loadScans = async () => {
  try {
    const res = await fetch('/api/scans')
    if (res.status === 401) { window.location.href = '/login'; return }
    scans.value = await res.json() || []
  } catch (e) {
    console.error(e)
  }
}

const getDuration = (start, end) => {
  if (!end) return '...'
  const ms = new Date(end) - new Date(start)
  const sec = Math.floor(ms / 1000)
  if (sec < 60) return `${sec}s`
  return `${Math.floor(sec / 60)}m ${sec % 60}s`
}

const getSeverityCount = (scan, sev) => {
  return scan.finding_counts?.[sev] || 0
}

const totalFindings = (scan) => {
  if (!scan.finding_counts) return 0
  return Object.values(scan.finding_counts).reduce((a, b) => a + b, 0)
}

onMounted(() => {
  loadScans()
  interval = setInterval(loadScans, 5000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>
