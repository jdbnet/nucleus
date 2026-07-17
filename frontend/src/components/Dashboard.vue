<template>
  <div class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- High Level Metrics -->
      <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle flex flex-col justify-center">
        <h3 class="text-sm font-medium text-body uppercase tracking-wider mb-2">Total Targets</h3>
        <p class="text-4xl font-bold text-heading">{{ stats.total_targets || 0 }}</p>
      </div>
      <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle flex flex-col justify-center">
        <h3 class="text-sm font-medium text-body uppercase tracking-wider mb-2">Total Scans</h3>
        <p class="text-4xl font-bold text-heading">{{ stats.total_scans || 0 }}</p>
      </div>
      <div class="bg-surface rounded-xl p-6 shadow-lg border border-red-500/30 flex flex-col justify-center">
        <h3 class="text-sm font-medium text-red-400 uppercase tracking-wider mb-2">Critical Findings</h3>
        <p class="text-4xl font-bold text-red-500">{{ stats.severity?.critical || 0 }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Severity Distribution Pie Chart -->
      <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle">
        <h3 class="text-lg font-semibold text-heading mb-4">Severity Distribution</h3>
        <div class="h-64 flex justify-center">
          <Pie v-if="chartData.labels.length" :data="chartData" :options="chartOptions" />
          <div v-else class="flex items-center justify-center h-full text-body">No data available</div>
        </div>
      </div>

      <!-- Top Vulnerable Hosts -->
      <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle flex flex-col">
        <h3 class="text-lg font-semibold text-heading mb-4">Top Vulnerable Hosts</h3>
        <div class="flex-1 overflow-y-auto">
          <ul v-if="stats.top_hosts?.length > 0" class="divide-y divide-border-subtle">
            <li v-for="h in stats.top_hosts" :key="h.host" class="py-3 flex justify-between items-center">
              <span class="text-heading font-mono text-sm">{{ h.host }}</span>
              <span class="px-2.5 py-1 bg-red-500/20 text-red-400 border border-red-500/50 rounded-md text-xs font-bold">{{ h.count }} findings</span>
            </li>
          </ul>
          <div v-else class="flex items-center justify-center h-full text-body">
            No hosts with findings.
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Trends -->
    <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle">
      <h3 class="text-lg font-semibold text-heading mb-4">Recent Discovery Trends (Last 7 Days)</h3>
      <div class="h-64">
        <Bar v-if="trendData.labels.length" :data="trendData" :options="trendOptions" />
        <div v-else class="flex items-center justify-center h-full text-body">No trends available</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js'
import { Pie, Bar } from 'vue-chartjs'

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title)

const stats = ref({})

const loadStats = async () => {
  try {
    const res = await fetch('/api/dashboard/stats')
    if (res.status === 401) { window.location.href = '/login'; return }
    stats.value = await res.json()
  } catch (e) {
    console.error(e)
  }
}

const chartData = computed(() => {
  if (!stats.value.severity) return { labels: [], datasets: [] }
  
  const sev = stats.value.severity
  const values = [sev.critical, sev.high, sev.medium, sev.low, sev.info]
  
  if (values.every(v => v === 0)) return { labels: [], datasets: [] }

  return {
    labels: ['Critical', 'High', 'Medium', 'Low', 'Info'],
    datasets: [
      {
        backgroundColor: ['#ef4444', '#f97316', '#eab308', '#3b82f6', '#4b5563'],
        data: values,
        borderWidth: 1,
        borderColor: '#161b22'
      }
    ]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right',
      labels: { color: '#8b949e' }
    }
  }
}

const trendData = computed(() => {
  if (!stats.value.recent_trends || stats.value.recent_trends.length === 0) return { labels: [], datasets: [] }
  
  return {
    labels: stats.value.recent_trends.map(t => t.date),
    datasets: [
      {
        label: 'Findings Discovered',
        backgroundColor: '#1ebe8a',
        data: stats.value.recent_trends.map(t => t.count),
        borderRadius: 4
      }
    ]
  }
})

const trendOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
      grid: { color: '#30363d' },
      ticks: { color: '#8b949e' }
    },
    x: {
      grid: { display: false },
      ticks: { color: '#8b949e' }
    }
  },
  plugins: {
    legend: { display: false }
  }
}

let pollInterval = null

onMounted(() => {
  loadStats()
  pollInterval = setInterval(loadStats, 10000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>
