<template>
  <div class="app">
    <header class="top">
      <div class="brand">
        <h1>多租户协同SaaS平台</h1>
        <p class="subtitle">阶段三：协作闭环与工程化</p>
      </div>
      <nav class="nav">
        <RouterLink to="/dashboard">业务看板</RouterLink>
        <RouterLink to="/workspace">协作</RouterLink>
        <RouterLink v-if="isAdmin" to="/admin">管理员</RouterLink>
        <RouterLink to="/bootstrap">初始化</RouterLink>
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
import { useSessionStore } from './stores/session'

const router = useRouter()
const session = useSessionStore()
const isAuthed = computed(() => session.isAuthenticated)
const isAdmin = computed(() => session.isAdmin)

function logout() {
  session.clearSession()
  router.push('/login')
}
</script>

<style scoped>
:global(body) {
  margin: 0;
  background: radial-gradient(1200px 600px at 10% -10%, #e7f1ff 0%, transparent 55%),
    radial-gradient(900px 500px at 110% 10%, #fff2df 0%, transparent 50%),
    #f6f8fb;
  color: #0f172a;
  font-family: 'Noto Sans SC', 'Source Han Sans SC', 'HarmonyOS Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.app {
  max-width: 1240px;
  margin: 0 auto;
  padding: 2.5rem 2.5rem 4rem;
  min-height: 100vh;
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  padding: 1.25rem 1.5rem;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 18px;
  backdrop-filter: blur(10px);
  box-shadow: 0 12px 40px rgba(15, 23, 42, 0.08);
}

.brand h1 {
  margin: 0 0 0.25rem 0;
  font-size: 2.1rem;
  letter-spacing: 0.5px;
}

.subtitle {
  margin: 0;
  color: #475569;
  font-size: 0.95rem;
}

.nav {
  display: flex;
  gap: 0.6rem;
  align-items: center;
  padding: 0.35rem;
  background: rgba(15, 23, 42, 0.04);
  border-radius: 999px;
}

.nav a {
  text-decoration: none;
  color: #1e293b;
  font-weight: 600;
  padding: 0.45rem 0.9rem;
  border-radius: 999px;
  transition: all 0.2s ease;
}

.nav a:hover {
  background: rgba(15, 23, 42, 0.08);
}

.ghost {
  border: 1px solid rgba(15, 23, 42, 0.2);
  background: white;
  padding: 0.45rem 0.9rem;
  border-radius: 999px;
  cursor: pointer;
  font-weight: 600;
}

.content {
  padding-top: 1.8rem;
}

:global(.nav a.router-link-active) {
  color: white;
  background: #0f172a;
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.2);
}
</style>


