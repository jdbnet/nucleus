import { createRouter, createWebHistory } from 'vue-router'
import TargetManager from './components/TargetManager.vue'
import ScansHistory from './components/ScansHistory.vue'
import FindingsInspector from './components/FindingsInspector.vue'
import Login from './components/Login.vue'
import Dashboard from './components/Dashboard.vue'

const routes = [
  { path: '/login', component: Login },
  { path: '/', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/targets', component: TargetManager, meta: { requiresAuth: true } },
  { path: '/scans', component: ScansHistory, meta: { requiresAuth: true } },
  { path: '/scans/:id', component: FindingsInspector, props: true, meta: { requiresAuth: true } }
]

export const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  if (to.meta.requiresAuth) {
    try {
      const res = await fetch('/api/auth/status')
      const status = await res.json()
      if (status.auth_required && !status.authenticated) {
        next('/login')
      } else {
        next()
      }
    } catch (e) {
      next('/login')
    }
  } else {
    next()
  }
})
