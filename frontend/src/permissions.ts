import type { Role } from './stores/session'

export type Permission =
  | 'tenant:read'
  | 'tenant:write'
  | 'user:read'
  | 'user:write'
  | 'project:read'
  | 'project:write'
  | 'task:read'
  | 'task:write'
  | 'audit:read'
  | 'audit:write'
  | 'admin:ping'

const matrix: Record<Permission, Role[]> = {
  'tenant:read': ['admin', 'manager'],
  'tenant:write': ['admin'],
  'user:read': ['admin', 'manager'],
  'user:write': ['admin'],
  'project:read': ['admin', 'manager', 'member', 'guest'],
  'project:write': ['admin', 'manager'],
  'task:read': ['admin', 'manager', 'member', 'guest'],
  'task:write': ['admin', 'manager', 'member'],
  'audit:read': ['admin'],
  'audit:write': ['admin'],
  'admin:ping': ['admin']
}

export function isAllowed(role: Role | null, perm: Permission): boolean {
  if (!role) return false
  return matrix[perm].includes(role)
}

