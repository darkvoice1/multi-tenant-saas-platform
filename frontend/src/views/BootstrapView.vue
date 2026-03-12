<template>
  <section class="card">
    <h2>初始化租户</h2>
    <p class="hint">首次使用时创建租户与管理员账号（仅开发环境）。</p>

    <form class="form" @submit.prevent="submit" v-if="!result">
      <label>
        租户名称
        <input v-model="form.tenantName" required />
      </label>
      <label>
        租户标识（slug）
        <input v-model="form.tenantSlug" required />
      </label>
      <label>
        管理员姓名
        <input v-model="form.adminName" required />
      </label>
      <label>
        管理员邮箱
        <input v-model="form.adminEmail" type="email" required />
      </label>
      <label>
        管理员密码
        <input v-model="form.adminPassword" type="password" required />
      </label>
      <button :disabled="loading">{{ loading ? '初始化中...' : '一键初始化' }}</button>
    </form>

    <div v-else class="result">
      <h3>初始化成功</h3>
      <p>租户 ID：<strong>{{ result.tenant_id }}</strong></p>
      <p>管理员邮箱：<strong>{{ result.email }}</strong></p>
      <p>管理员密码：<strong>{{ result.password }}</strong></p>
      <button class="ghost" @click="goLogin">去登录</button>
    </div>

    <p v-if="error" class="error">{{ error }}</p>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { API_BASE } from '../api'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const result = ref<any | null>(null)

const form = reactive({
  tenantName: 'Demo Tenant',
  tenantSlug: 'demo',
  adminName: 'Admin',
  adminEmail: 'admin@example.com',
  adminPassword: 'Admin123!'
})

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch(`${API_BASE}/auth/bootstrap`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        tenant_name: form.tenantName,
        tenant_slug: form.tenantSlug,
        admin_name: form.adminName,
        admin_email: form.adminEmail,
        admin_password: form.adminPassword
      })
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data.error || '初始化失败')
    }
    result.value = await res.json()
  } catch (err: any) {
    error.value = err.message || '初始化失败'
  } finally {
    loading.value = false
  }
}

function goLogin() {
  if (!result.value) return
  router.push({
    path: '/login',
    query: { tenant: result.value.tenant_id, email: result.value.email }
  })
}
</script>

<style scoped>
.card {
  max-width: 520px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  background: white;
}

.form {
  display: grid;
  gap: 0.8rem;
  margin-top: 1rem;
}

label {
  display: grid;
  gap: 0.4rem;
  color: #334155;
}

input {
  padding: 0.5rem 0.7rem;
  border: 1px solid #cbd5f5;
  border-radius: 8px;
}

button {
  padding: 0.6rem 1rem;
  border: none;
  background: #0f172a;
  color: white;
  border-radius: 8px;
  cursor: pointer;
}

button.ghost {
  background: transparent;
  color: #0f172a;
  border: 1px solid #cbd5f5;
}

.hint {
  color: #64748b;
}

.error {
  margin-top: 0.6rem;
  color: #ef4444;
}

.result {
  margin-top: 1rem;
  display: grid;
  gap: 0.4rem;
}
</style>