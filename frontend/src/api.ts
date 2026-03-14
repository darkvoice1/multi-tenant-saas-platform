const runtimeBase =
  typeof window !== 'undefined'
    ? `${window.location.protocol}//${window.location.hostname}:8080`
    : ''

export const API_BASE = (import.meta.env.VITE_API_BASE as string) || runtimeBase || 'http://localhost:8080'

type ApiError = { error?: string }

function clearSessionStorage() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('role')
  localStorage.removeItem('tenant_id')
}

function buildHeaders(init?: HeadersInit): Headers {
  const headers = new Headers(init || {})
  const token = localStorage.getItem('access_token')
  const tenantId = localStorage.getItem('tenant_id')

  if (token) headers.set('Authorization', `Bearer ${token}`)
  if (tenantId) headers.set('X-Tenant-ID', tenantId)

  return headers
}

async function parseError(res: Response): Promise<string> {
  if (res.status === 401) {
    // Token missing/expired/invalid (often after backend restart or JWT_SECRET change).
    clearSessionStorage()
    // Best-effort redirect back to login; keep throwing so callers can show error too.
    try {
      if (typeof window !== 'undefined') window.location.assign('/login')
    } catch {
      // ignore
    }
    return '未登录或登录已过期，请重新登录'
  }

  const data = (await res.json().catch(() => ({}))) as ApiError
  return data.error || `请求失败 (${res.status})`
}

export async function apiFetch<T = unknown>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = buildHeaders(options.headers)

  if (!(options.body instanceof FormData)) {
    if (options.body && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })
  if (!res.ok) {
    throw new Error(await parseError(res))
  }

  const text = await res.text()
  if (!text) return null as T

  try {
    return JSON.parse(text) as T
  } catch {
    return text as T
  }
}

export async function apiFetchBlob(path: string, options: RequestInit = {}): Promise<Blob> {
  const headers = buildHeaders(options.headers)
  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })
  if (!res.ok) {
    throw new Error(await parseError(res))
  }
  return await res.blob()
}