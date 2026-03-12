<template>
  <section class="card">
    <h2>登录</h2>
    <p class="hint">
      使用已初始化的管理员账号或你自己的账号。
      <RouterLink class="link" to="/bootstrap">没有账号？去初始化</RouterLink>
    </p>
    <form @submit.prevent="submit">
      <label>
        租户 ID
        <input v-model="form.tenantId" placeholder="租户 UUID" required />
      </label>
      <label>
        邮箱
        <input v-model="form.email" type="email" required />
      </label>
      <label>
        密码
        <input v-model="form.password" type="password" required />
      </label>
      <button :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const route = useRoute()
const session = useSessionStore()
const loading = ref(false)
const error = ref('')
const form = reactive({
  tenantId: '',
  email: '',
  password: ''
})

if (typeof route.query.tenant === 'string') {
  form.tenantId = route.query.tenant
}
if (typeof route.query.email === 'string') {
  form.email = route.query.email
}

async function submit() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch('http://localhost:8080/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        tenant_id: form.tenantId,
        email: form.email,
        password: form.password
      })
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data.error || '登录失败')
    }
    const data = await res.json()
    session.setSession(data.access_token, data.user.role, data.user.tenant_id)
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.card {
  max-width: 420px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  background: white;
}
label {
  display: grid;
  gap: 0.4rem;
  margin: 0.8rem 0;
  color: #334155;
}
input {
  padding: 0.5rem 0.7rem;
  border: 1px solid #cbd5f5;
  border-radius: 8px;
}
button {
  margin-top: 0.8rem;
  padding: 0.6rem 1rem;
  border: none;
  background: #0f172a;
  color: white;
  border-radius: 8px;
  cursor: pointer;
}
.error {
  margin-top: 0.6rem;
  color: #ef4444;
}
.hint {
  color: #64748b;
}
.link {
  margin-left: 0.4rem;
  color: #0f172a;
  text-decoration: underline;
}
</style>
