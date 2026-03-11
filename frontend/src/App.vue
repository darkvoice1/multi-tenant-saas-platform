<template>
  <div class="app">
    <header class="top">
      <div class="brand">
        <h1>多租户协同SaaS平台</h1>
        <p class="subtitle">阶段二：认证与RBAC</p>
      </div>
      <nav class="nav">
        <RouterLink to="/dashboard">控制台</RouterLink>
        <RouterLink to="/admin">管理员</RouterLink>
        <RouterLink to="/login">登录</RouterLink>
        <button v-if="isAuthed" class="ghost" @click="logout">退出登录</button>
      </nav>
    </header>
    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { clearSession, getToken } from './auth'

const router = useRouter()
const isAuthed = computed(() => getToken() !== null)

function logout() {
  clearSession()
  router.push('/login')
}
</script>

<style scoped>
.app {
  font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, sans-serif;
  padding: 2rem 2.5rem 3rem;
  color: #0f172a;
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 1.5rem;
}

.brand h1 {
  margin: 0 0 0.25rem 0;
  font-size: 2rem;
}

.subtitle {
  margin: 0;
  color: #64748b;
}

.nav {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.nav a {
  text-decoration: none;
  color: #1e293b;
  font-weight: 600;
}

.ghost {
  border: 1px solid #cbd5f5;
  background: transparent;
  padding: 0.4rem 0.8rem;
  border-radius: 8px;
  cursor: pointer;
}

.content {
  padding-top: 1.5rem;
}
</style>
