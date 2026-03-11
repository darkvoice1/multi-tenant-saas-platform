export type Role = 'admin' | 'manager' | 'member' | 'guest'

const TOKEN_KEY = 'access_token'
const ROLE_KEY = 'role'
const TENANT_KEY = 'tenant_id'

export function setSession(token: string, role: Role, tenantId: string) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(ROLE_KEY, role)
  localStorage.setItem(TENANT_KEY, tenantId)
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(ROLE_KEY)
  localStorage.removeItem(TENANT_KEY)
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getRole(): Role | null {
  return localStorage.getItem(ROLE_KEY) as Role | null
}

export function getTenantId(): string | null {
  return localStorage.getItem(TENANT_KEY)
}

export function isAuthenticated(): boolean {
  return getToken() !== null
}

export function hasAnyRole(roles: string[]): boolean {
  const role = getRole()
  if (!role) return false
  return roles.includes(role)
}
