import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import type { Role } from './stores/session'
import type { Permission } from './permissions'
import { isAllowed } from './permissions'

function getToken(): string | null {
  return localStorage.getItem('access_token')
}

function getRole(): Role | null {
  return (localStorage.getItem('role') as Role | null) ?? null
}

function isAuthenticated(): boolean {
  return getToken() !== null
}

function hasAnyRole(roles: Role[]): boolean {
  const role = getRole()
  if (!role) return false
  return roles.includes(role)
}

function hasAllPerms(perms: Permission[]): boolean {
  const role = getRole()
  return perms.every((p) => isAllowed(role, p))
}

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    roles?: Role[]
    perms?: Permission[]
  }
}

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', component: () => import('./views/LoginView.vue'), meta: { public: true } },
  { path: '/bootstrap', component: () => import('./views/BootstrapView.vue'), meta: { public: true } },
  { path: '/dashboard', component: () => import('./views/DashboardView.vue') },
  { path: '/workspace', component: () => import('./views/WorkspaceView.vue'), meta: { perms: ['project:read'] } },
  { path: '/admin', component: () => import('./views/AdminView.vue'), meta: { roles: ['admin'] } }
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
  if (to.meta.roles && !hasAnyRole(to.meta.roles)) {
    return { path: '/dashboard' }
  }
  if (to.meta.perms && !hasAllPerms(to.meta.perms)) {
    return { path: '/dashboard' }
  }
  return true
})

export default router

