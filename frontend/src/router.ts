import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { isAuthenticated, hasAnyRole } from './auth'
import LoginView from './views/LoginView.vue'
import DashboardView from './views/DashboardView.vue'
import AdminView from './views/AdminView.vue'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: LoginView, meta: { public: true } },
  { path: '/dashboard', component: DashboardView },
  { path: '/admin', component: AdminView, meta: { roles: ['admin'] } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  if (to.meta.public) return true
  if (!isAuthenticated()) {
    return { path: '/login' }
  }
  const roles = to.meta.roles as string[] | undefined
  if (roles && !hasAnyRole(roles)) {
    return { path: '/dashboard' }
  }
  return true
})

export default router
