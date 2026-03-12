import { defineStore } from 'pinia'

export type Role = 'admin' | 'manager' | 'member' | 'guest'

type SessionState = {
  accessToken: string | null
  role: Role | null
  tenantId: string | null
}

const TOKEN_KEY = 'access_token'
const ROLE_KEY = 'role'
const TENANT_KEY = 'tenant_id'

export const useSessionStore = defineStore('session', {
  state: (): SessionState => ({
    accessToken: localStorage.getItem(TOKEN_KEY),
    role: (localStorage.getItem(ROLE_KEY) as Role | null) ?? null,
    tenantId: localStorage.getItem(TENANT_KEY)
  }),
  actions: {
    setSession(accessToken: string, role: Role, tenantId: string) {
      this.accessToken = accessToken
      this.role = role
      this.tenantId = tenantId
      localStorage.setItem(TOKEN_KEY, accessToken)
      localStorage.setItem(ROLE_KEY, role)
      localStorage.setItem(TENANT_KEY, tenantId)
    },
    clearSession() {
      this.accessToken = null
      this.role = null
      this.tenantId = null
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem(ROLE_KEY)
      localStorage.removeItem(TENANT_KEY)
    }
  },
  getters: {
    isAuthenticated: (s) => s.accessToken !== null,
    isAdmin: (s) => s.role === 'admin'
  }
})

