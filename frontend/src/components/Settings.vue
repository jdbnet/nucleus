<template>
  <div class="space-y-6">
    <div v-if="message" class="p-3 rounded-lg text-sm text-center border" :class="messageType === 'error' ? 'bg-red-500/20 border-red-500/50 text-red-400' : 'bg-accent/20 border-accent/50 text-accent'">
      {{ message }}
    </div>

    <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle backdrop-blur-sm">
      <h2 class="text-xl font-semibold mb-4 text-heading">General</h2>
      <form @submit.prevent="saveSettings" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-body mb-1">Instance URL</label>
          <input v-model="form.instance_url" type="url" placeholder="https://nucleus.example.com" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          <p class="text-xs text-body mt-1">Used to generate "View scan" links in email and webhook notifications.</p>
        </div>

        <h3 class="text-lg font-semibold text-heading pt-2">Email (SMTP)</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-body mb-1">SMTP Host</label>
            <input v-model="form.smtp_host" type="text" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">SMTP Port</label>
            <input v-model="form.smtp_port" type="text" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">From Address</label>
            <input v-model="form.smtp_from" type="email" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">To Address</label>
            <input v-model="form.smtp_to" type="email" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">SMTP Username</label>
            <input v-model="form.smtp_user" type="text" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">SMTP Password</label>
            <input v-model="form.smtp_pass" type="password" :placeholder="smtpPassSet ? '********' : ''" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
        </div>

        <h3 class="text-lg font-semibold text-heading pt-2">Notifications</h3>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Minimum Severity to Notify</label>
          <select v-model="form.notify_min_severity" class="w-full md:w-1/2 bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
            <option value="always">Always (every scan, including zero findings)</option>
            <option value="info">Info and above</option>
            <option value="low">Low and above</option>
            <option value="medium">Medium and above</option>
            <option value="high">High and above</option>
            <option value="critical">Critical only</option>
          </select>
          <p class="text-xs text-body mt-1">Controls when email and webhook notifications are sent after a scan completes.</p>
        </div>

        <h3 class="text-lg font-semibold text-heading pt-2">Data Retention</h3>
        <div class="md:w-1/3">
          <label class="block text-sm font-medium text-body mb-1">Retention Days</label>
          <input v-model="form.retention_days" type="number" min="1" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
        </div>

        <h3 class="text-lg font-semibold text-heading pt-2">Web Access</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-body mb-1">Web Username</label>
            <input v-model="form.web_user" type="text" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-body mb-1">Web Password</label>
            <input v-model="form.web_pass" type="password" :placeholder="webPassSet ? '********' : ''" class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
          </div>
        </div>
        <p class="text-xs text-body">Leave username and password empty to disable dashboard authentication.</p>

        <button type="submit" class="bg-accent hover:bg-accent/80 text-bg font-bold py-2 px-6 rounded-lg shadow-md hover:shadow-lg transition-all hover:cursor-pointer">
          Save Settings
        </button>
      </form>
    </div>

    <div class="bg-surface rounded-xl p-6 shadow-lg border border-border-subtle backdrop-blur-sm">
      <h2 class="text-xl font-semibold mb-4 text-heading">{{ editingWebhook ? 'Edit Webhook' : 'Add Webhook' }}</h2>
      <form @submit.prevent="saveWebhook" class="grid grid-cols-1 md:grid-cols-5 gap-4 items-end">
        <div>
          <label class="block text-sm font-medium text-body mb-1">Name</label>
          <input v-model="webhookForm.name" type="text" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
        </div>
        <div>
          <label class="block text-sm font-medium text-body mb-1">Provider</label>
          <select v-model="webhookForm.provider" required class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
            <option value="discord">Discord</option>
            <option value="slack">Slack</option>
            <option value="gotify">Gotify</option>
            <option value="ntfy">Ntfy</option>
            <option value="teams">Microsoft Teams</option>
            <option value="generic">Generic</option>
          </select>
        </div>
        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-body mb-1">Webhook URL</label>
          <input v-model="webhookForm.url" type="url" required placeholder="https://..." class="w-full bg-bg border border-border-subtle rounded-lg px-4 py-2 text-heading focus:ring-2 focus:ring-accent focus:border-transparent outline-none transition-all">
        </div>
        <div class="flex items-center gap-3">
          <label class="flex items-center gap-2 text-sm text-body">
            <input v-model="webhookForm.enabled" type="checkbox" class="rounded border-border-subtle">
            Enabled
          </label>
          <button type="submit" class="bg-accent hover:bg-accent/80 text-bg font-bold py-2 px-4 rounded-lg shadow-md hover:cursor-pointer">
            {{ editingWebhook ? 'Update' : 'Add' }}
          </button>
          <button v-if="editingWebhook" type="button" @click="resetWebhookForm" class="text-body hover:text-heading text-sm hover:cursor-pointer">Cancel</button>
        </div>
      </form>
    </div>

    <div class="bg-surface rounded-xl overflow-hidden shadow-lg border border-border-subtle">
      <div class="px-6 py-4 border-b border-border-subtle">
        <h2 class="text-xl font-semibold text-heading">Webhooks</h2>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-bg/50 text-body">
            <tr>
              <th class="px-6 py-3 font-medium">Name</th>
              <th class="px-6 py-3 font-medium">Provider</th>
              <th class="px-6 py-3 font-medium">URL</th>
              <th class="px-6 py-3 font-medium">Enabled</th>
              <th class="px-6 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border-subtle">
            <tr v-for="wh in webhooks" :key="wh.id" class="hover:bg-border-subtle/30 transition-colors">
              <td class="px-6 py-4 text-heading font-medium">{{ wh.name }}</td>
              <td class="px-6 py-4 text-body capitalize">{{ wh.provider }}</td>
              <td class="px-6 py-4 text-body font-mono text-xs truncate max-w-xs">{{ wh.url }}</td>
              <td class="px-6 py-4 text-body">{{ wh.enabled ? 'Yes' : 'No' }}</td>
              <td class="px-6 py-4 text-right space-x-3">
                <button @click="testWebhook(wh.id)" class="text-accent hover:text-accent/80 font-medium transition-colors hover:cursor-pointer">Test</button>
                <button @click="editWebhook(wh)" class="text-blue-400 hover:text-blue-300 font-medium transition-colors hover:cursor-pointer">Edit</button>
                <button @click="deleteWebhook(wh.id)" class="text-red-400 hover:text-red-300 font-medium transition-colors hover:cursor-pointer">Delete</button>
              </td>
            </tr>
            <tr v-if="webhooks.length === 0">
              <td colspan="5" class="px-6 py-8 text-center text-body">No webhooks configured.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'

const form = ref({
  instance_url: '',
  smtp_host: '',
  smtp_port: '',
  smtp_from: '',
  smtp_to: '',
  smtp_user: '',
  smtp_pass: '',
  retention_days: '30',
  notify_min_severity: 'always',
  web_user: '',
  web_pass: '',
})

const webhookForm = ref({
  name: '',
  provider: 'discord',
  url: '',
  enabled: true,
})

const webhooks = ref([])
const editingWebhook = ref(false)
const editingWebhookId = ref(null)
const smtpPassSet = ref(false)
const webPassSet = ref(false)
const message = ref('')
const messageType = ref('success')

const showMessage = (text, type = 'success') => {
  message.value = text
  messageType.value = type
  setTimeout(() => { message.value = '' }, 4000)
}

const loadSettings = async () => {
  const res = await fetch('/api/settings')
  if (res.status === 401) { window.location.href = '/login'; return }
  const data = await res.json()
  form.value = {
    instance_url: data.instance_url || '',
    smtp_host: data.smtp_host || '',
    smtp_port: data.smtp_port || '',
    smtp_from: data.smtp_from || '',
    smtp_to: data.smtp_to || '',
    smtp_user: data.smtp_user || '',
    smtp_pass: '',
    retention_days: data.retention_days || '30',
    notify_min_severity: data.notify_min_severity || 'always',
    web_user: data.web_user || '',
    web_pass: '',
  }
  smtpPassSet.value = data.smtp_pass === '********'
  webPassSet.value = data.web_pass === '********'
}

const loadWebhooks = async () => {
  const res = await fetch('/api/webhooks')
  if (res.status === 401) { window.location.href = '/login'; return }
  webhooks.value = await res.json() || []
}

const saveSettings = async () => {
  const payload = { ...form.value }
  if (!payload.smtp_pass) delete payload.smtp_pass
  if (!payload.web_pass) delete payload.web_pass

  const res = await fetch('/api/settings', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (res.ok) {
    showMessage('Settings saved successfully')
    await loadSettings()
  } else {
    const text = await res.text()
    showMessage(text || 'Failed to save settings', 'error')
  }
}

const resetWebhookForm = () => {
  webhookForm.value = { name: '', provider: 'discord', url: '', enabled: true }
  editingWebhook.value = false
  editingWebhookId.value = null
}

const saveWebhook = async () => {
  const payload = { ...webhookForm.value }
  const url = editingWebhook.value ? `/api/webhooks/${editingWebhookId.value}` : '/api/webhooks'
  const method = editingWebhook.value ? 'PUT' : 'POST'

  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (res.ok) {
    showMessage(editingWebhook.value ? 'Webhook updated' : 'Webhook added')
    resetWebhookForm()
    await loadWebhooks()
  } else {
    const text = await res.text()
    showMessage(text || 'Failed to save webhook', 'error')
  }
}

const editWebhook = (wh) => {
  webhookForm.value = {
    name: wh.name,
    provider: wh.provider,
    url: wh.url,
    enabled: wh.enabled,
  }
  editingWebhook.value = true
  editingWebhookId.value = wh.id
}

const deleteWebhook = async (id) => {
  if (!confirm('Delete this webhook?')) return
  const res = await fetch(`/api/webhooks/${id}`, { method: 'DELETE' })
  if (res.ok) {
    showMessage('Webhook deleted')
    await loadWebhooks()
  } else {
    showMessage('Failed to delete webhook', 'error')
  }
}

const testWebhook = async (id) => {
  const res = await fetch(`/api/webhooks/${id}/test`, { method: 'POST' })
  if (res.ok) {
    showMessage('Test notification sent')
  } else {
    const text = await res.text()
    showMessage(text || 'Test failed', 'error')
  }
}

onMounted(async () => {
  await loadSettings()
  await loadWebhooks()
})
</script>
