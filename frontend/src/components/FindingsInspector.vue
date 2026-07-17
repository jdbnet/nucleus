<template>
  <div class="space-y-6">
    <div class="flex items-center space-x-4 mb-6">
      <router-link to="/scans" class="text-body hover:text-heading transition-colors flex items-center">
        <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path></svg>
        Back
      </router-link>
      <h2 class="text-2xl font-bold text-heading">Scan Findings Inspector</h2>
    </div>

    <div class="bg-surface rounded-xl overflow-hidden shadow-lg border border-border-subtle">
      <div class="px-6 py-4 border-b border-border-subtle flex flex-col md:flex-row md:justify-between md:items-center gap-4">
        <h3 class="text-lg font-medium text-heading whitespace-nowrap">
          Findings ({{ filteredFindings.length }} / {{ findings.length }})
        </h3>
        
        <div class="flex flex-col md:flex-row items-center gap-4 w-full md:w-auto">
          <!-- Severity Filters -->
          <div class="flex flex-wrap items-center gap-3 text-sm">
            <label v-for="sev in ['critical', 'high', 'medium', 'low', 'info']" :key="sev" class="flex items-center space-x-1 cursor-pointer">
              <input type="checkbox" v-model="selectedSeverities" :value="sev" class="rounded border-border-subtle bg-bg text-accent focus:ring-accent">
              <span class="capitalize text-body">{{ sev }}</span>
            </label>
          </div>
          
          <!-- Search Box -->
          <div class="relative w-full md:w-64">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-4 w-4 text-body" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
            </div>
            <input v-model="searchQuery" type="text" placeholder="Search findings..." class="w-full bg-bg border border-border-subtle rounded-lg pl-9 pr-4 py-2 text-sm text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          
          <button @click="loadFindings" class="text-body hover:text-heading transition-colors hover:cursor-pointer p-2 bg-bg rounded-lg border border-border-subtle" title="Refresh">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          </button>
        </div>
      </div>
      
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-bg/50 text-body">
            <tr>
              <th class="px-6 py-3 font-medium">Severity</th>
              <th class="px-6 py-3 font-medium">Name & Template</th>
              <th class="px-6 py-3 font-medium">Host</th>
              <th class="px-6 py-3 font-medium">Matched At</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-subtle">
            <tr v-for="finding in filteredFindings" :key="finding.id" class="hover:bg-border-subtle/30 transition-colors">
              <td class="px-6 py-4">
                <span :class="severityClass(finding.severity)" class="px-2.5 py-1 rounded-md text-xs font-bold uppercase tracking-wider border">
                  {{ finding.severity }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="text-heading font-medium">{{ finding.name }}</div>
                <div class="text-body font-mono text-xs mt-1">{{ finding.template_id }}</div>
                <div v-if="finding.description" class="text-body text-xs mt-2 max-w-md truncate" :title="finding.description">
                  {{ finding.description }}
                </div>
              </td>
              <td class="px-6 py-4 text-body font-mono text-xs">{{ finding.host }}</td>
              <td class="px-6 py-4 text-body text-xs break-all max-w-xs">{{ finding.matched_at }}</td>
            </tr>
            <tr v-if="filteredFindings.length === 0">
              <td colspan="4" class="px-6 py-12 text-center">
                <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-bg mb-4 border border-border-subtle">
                  <svg class="w-8 h-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
                </div>
                <h3 class="text-lg font-medium text-heading mb-1">No matches found</h3>
                <p class="text-body">Try adjusting your search query or severity filters.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  id: { type: String, required: true }
})

const findings = ref([])
const searchQuery = ref('')
const selectedSeverities = ref(['critical', 'high', 'medium', 'low', 'info'])

const filteredFindings = computed(() => {
  let result = findings.value
  
  // Filter by severity
  if (selectedSeverities.value.length > 0) {
    result = result.filter(f => selectedSeverities.value.includes(f.severity.toLowerCase()))
  } else {
    // If no severities are selected, show none
    return []
  }

  // Filter by search query
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(f => {
      return (f.name && f.name.toLowerCase().includes(query)) ||
             (f.template_id && f.template_id.toLowerCase().includes(query)) ||
             (f.host && f.host.toLowerCase().includes(query)) ||
             (f.description && f.description.toLowerCase().includes(query))
    })
  }
  
  return result
})

const loadFindings = async () => {
  try {
    const res = await fetch(`/api/scans/${props.id}/findings`)
    if (res.status === 401) { window.location.href = '/login'; return }
    findings.value = await res.json() || []
  } catch (e) {
    console.error(e)
  }
}

const severityClass = (sev) => {
  const map = {
    critical: 'bg-red-500/20 text-red-400 border-red-500/50',
    high: 'bg-orange-500/20 text-orange-400 border-orange-500/50',
    medium: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/50',
    low: 'bg-blue-500/20 text-blue-400 border-blue-500/50',
    info: 'bg-border-subtle/50 text-body border-border-subtle'
  }
  return map[sev] || map.info
}

onMounted(loadFindings)
</script>
