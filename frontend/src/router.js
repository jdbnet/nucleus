import { createRouter, createWebHistory } from 'vue-router'
import TargetManager from './components/TargetManager.vue'
import ScansHistory from './components/ScansHistory.vue'
import FindingsInspector from './components/FindingsInspector.vue'

const routes = [
  { path: '/', component: TargetManager },
  { path: '/scans', component: ScansHistory },
  { path: '/scans/:id', component: FindingsInspector, props: true }
]

export const router = createRouter({
  history: createWebHistory(),
  routes
})
