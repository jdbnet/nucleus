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
      <div class="px-6 py-4 border-b border-border-subtle flex justify-between items-center">
        <h3 class="text-lg font-medium text-heading">All Findings ({{ findings.length }})</h3>
        <button @click="loadFindings" class="text-body hover:text-heading transition-colors hover:cursor-pointer">Refresh</button>
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
            <tr v-for="finding in findings" :key="finding.id" class="hover:bg-border-subtle/30 transition-colors">
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
            <tr v-if="findings.length === 0">
              <td colspan="4" class="px-6 py-12 text-center">
                <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-bg mb-4 border border-border-subtle">
                  <svg class="w-8 h-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                </div>
                <h3 class="text-lg font-medium text-heading mb-1">No Findings</h3>
                <p class="text-body">This scan completed cleanly without any detected vulnerabilities.</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  id: { type: String, required: true }
})

const findings = ref([])

const loadFindings = async () => {
  try {
    const res = await fetch(`/api/scans/${props.id}/findings`)
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
