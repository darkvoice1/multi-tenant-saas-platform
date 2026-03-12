<template>
  <section class="page">
    <header class="head">
      <div>
        <h2>业务看板</h2>
        <p class="muted">聚合租户内项目/任务/通知与近期动态。</p>
      </div>
      <div class="head-actions">
        <button class="ghost" @click="load" :disabled="loading">刷新</button>
        <RouterLink class="btn" to="/workspace">进入协作</RouterLink>
      </div>
    </header>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="data" class="grid">
      <div class="card card--tenant card--span-4">
        <h3>租户</h3>
        <p class="kv">
          <span class="kv-label">名称</span>
          <strong class="kv-value">{{ data.tenant.name }}</strong>
        </p>
        <p class="kv">
          <span class="kv-label">Slug</span>
          <strong class="kv-value">{{ data.tenant.slug }}</strong>
        </p>
        <p class="kv">
          <span class="kv-label">ID</span>
          <strong class="kv-value mono">{{ data.tenant.id }}</strong>
        </p>
        <p class="kv">
          <span class="kv-label">存储</span>
          <strong class="kv-value">{{ data.lists.storage_backend }}</strong>
        </p>
      </div>

      <div class="card card--kpi card--span-4">
        <h3>核心指标</h3>
        <div class="kpi">
          <div class="kpi-item">
            <div class="kpi-num">{{ data.metrics.project_count }}</div>
            <div class="kpi-label">项目</div>
          </div>
          <div class="kpi-item">
            <div class="kpi-num">{{ data.metrics.task_count }}</div>
            <div class="kpi-label">任务</div>
          </div>
          <div class="kpi-item">
            <div class="kpi-num">{{ data.metrics.my_open_task_count }}</div>
            <div class="kpi-label">我的未完成任务</div>
          </div>
          <div class="kpi-item">
            <div class="kpi-num">{{ data.metrics.pending_review_task_count }}</div>
            <div class="kpi-label">待审批</div>
          </div>
          <div class="kpi-item">
            <div class="kpi-num">{{ data.metrics.unread_notification_count }}</div>
            <div class="kpi-label">未读通知</div>
          </div>
        </div>
      </div>

      <div class="card card--status card--span-4">
        <h3>任务状态分布</h3>
        <ul class="list">
          <li v-for="s in data.metrics.task_status_counts" :key="s.status" class="row">
            <span class="tag" :class="'tag--' + s.status">{{ statusLabel(s.status) }}</span>
            <span class="muted">{{ s.count }}</span>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-6">
        <h3>团队 7 天内到期</h3>
        <div v-if="data.lists.due_soon_tasks.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="t in data.lists.due_soon_tasks" :key="t.id" class="row">
            <RouterLink class="linklike" :to="{ path: '/workspace', query: { task: t.id } }">{{ t.title }}</RouterLink>
            <span class="muted">{{ t.due_at ? fmt(t.due_at) : '-' }}</span>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-6">
        <h3>我的 7 天内到期</h3>
        <div v-if="data.lists.my_due_soon_tasks.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="t in data.lists.my_due_soon_tasks" :key="t.id" class="row">
            <RouterLink class="linklike" :to="{ path: '/workspace', query: { task: t.id } }">{{ t.title }}</RouterLink>
            <span class="muted">{{ t.due_at ? fmt(t.due_at) : '-' }}</span>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-6">
        <h3>最近更新任务</h3>
        <div v-if="data.lists.recent_tasks.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="t in data.lists.recent_tasks" :key="t.id" class="row">
            <RouterLink class="linklike" :to="{ path: '/workspace', query: { task: t.id } }">{{ t.title }}</RouterLink>
            <span class="muted">{{ statusLabel(t.status) }}</span>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-6">
        <h3>待审批任务</h3>
        <div v-if="data.lists.pending_review_tasks.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="t in data.lists.pending_review_tasks" :key="t.id" class="row">
            <RouterLink class="linklike" :to="{ path: '/workspace', query: { task: t.id } }">{{ t.title }}</RouterLink>
            <span class="muted">{{ t.priority }}</span>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-8">
        <h3>最近评论</h3>
        <div v-if="data.lists.recent_comments.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="c in data.lists.recent_comments" :key="c.id" class="comment">
            <div class="comment-head">
              <strong>{{ c.user_name || c.user_email }}</strong>
              <span class="muted">{{ fmt(c.created_at) }}</span>
            </div>
            <RouterLink class="comment-body linklike" :to="{ path: '/workspace', query: { task: c.task_id } }">{{
              c.content
            }}</RouterLink>
          </li>
        </ul>
      </div>

      <div class="card card--list card--span-4">
        <h3>最近通知</h3>
        <div v-if="data.lists.recent_notifications.length === 0" class="muted">暂无</div>
        <ul v-else class="list">
          <li v-for="n in data.lists.recent_notifications" :key="n.id" class="row">
            <span class="title">{{ n.message }}</span>
            <span class="muted">{{ n.read_at ? '已读' : '未读' }}</span>
          </li>
        </ul>
      </div>
    </div>

    <div v-else class="muted">加载中...</div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiFetch } from '../api'
import type { DashboardSummary } from '../types'

const loading = ref(false)
const error = ref('')
const data = ref<DashboardSummary | null>(null)

function fmt(iso: string) {
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

function statusLabel(s: string) {
  switch (s) {
    case 'todo':
      return '待办'
    case 'in_progress':
      return '进行中'
    case 'review':
      return '待审批'
    case 'done':
      return '已完成'
    case 'rejected':
      return '已拒绝'
    default:
      return s
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    data.value = await apiFetch<DashboardSummary>('/api/dashboard')
  } catch (e: any) {
    error.value = e?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 1.6rem;
}

.head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 1rem;
}

.head h2 {
  margin: 0;
  font-size: 1.8rem;
  letter-spacing: 0.4px;
}

.head p {
  margin: 0.4rem 0 0 0;
}

.head-actions {
  display: flex;
  gap: 0.6rem;
  align-items: center;
}

.grid {
  display: grid;
  gap: 1.1rem;
  grid-template-columns: repeat(12, minmax(0, 1fr));
}

.card--span-4 { grid-column: span 4; }
.card--span-6 { grid-column: span 6; }
.card--span-8 { grid-column: span 8; }
.card--span-12 { grid-column: span 12; }

.card {
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 16px;
  padding: 1.1rem 1.2rem;
  background: white;
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.06);
  position: relative;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  animation: fadeUp 0.45s ease both;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.1);
}

.card--tenant,
.card--kpi,
.card--status {
  min-height: 220px;
}

.card--list {
  min-height: 200px;
}

.card:nth-child(1) { animation-delay: 0.02s; }
.card:nth-child(2) { animation-delay: 0.04s; }
.card:nth-child(3) { animation-delay: 0.06s; }
.card:nth-child(4) { animation-delay: 0.08s; }
.card:nth-child(5) { animation-delay: 0.1s; }
.card:nth-child(6) { animation-delay: 0.12s; }
.card:nth-child(7) { animation-delay: 0.14s; }
.card:nth-child(8) { animation-delay: 0.16s; }

.kv {
  display: grid;
  grid-template-columns: 84px 1fr;
  align-items: center;
  gap: 0.6rem;
  margin: 0.4rem 0;
  padding: 0.45rem 0;
  border-bottom: 1px dashed rgba(15, 23, 42, 0.08);
}

.kv:last-of-type {
  border-bottom: none;
}

.kv-label {
  color: #64748b;
  font-size: 0.86rem;
  letter-spacing: 0.6px;
  text-transform: uppercase;
}

.kv-value {
  text-align: right;
  font-weight: 700;
  color: #0f172a;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.85rem;
}

.kpi {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.8rem;
  margin-top: 0.6rem;
}

.kpi-item {
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 14px;
  padding: 0.75rem 0.8rem;
  background: linear-gradient(135deg, #f8fafc, #eef2ff);
}

.kpi-num {
  font-size: 1.6rem;
  font-weight: 800;
}

.kpi-label {
  color: #475569;
  margin-top: 0.2rem;
}

.list {
  list-style: none;
  padding: 0;
  margin: 0.6rem 0 0 0;
  display: grid;
  gap: 0.55rem;
}

.row {
  display: flex;
  justify-content: space-between;
  gap: 0.8rem;
  align-items: center;
}

.title {
  font-weight: 600;
}

.linklike {
  color: inherit;
  text-decoration: none;
  font-weight: 600;
}

.linklike:hover {
  text-decoration: underline;
}

.tag {
  padding: 0.15rem 0.55rem;
  border-radius: 999px;
  font-size: 0.82rem;
  font-weight: 600;
  background: #e2e8f0;
  color: #1e293b;
}

.tag--todo { background: #e2e8f0; color: #1e293b; }
.tag--in_progress { background: #dbeafe; color: #1d4ed8; }
.tag--review { background: #fef3c7; color: #b45309; }
.tag--done { background: #dcfce7; color: #166534; }
.tag--rejected { background: #fee2e2; color: #b91c1c; }

.comment {
  border: 1px solid rgba(15, 23, 42, 0.08);
  border-radius: 12px;
  padding: 0.6rem;
  background: #f8fafc;
}

.comment-head {
  display: flex;
  justify-content: space-between;
  gap: 0.6rem;
}

.comment-body {
  margin-top: 0.35rem;
  color: #0f172a;
  white-space: pre-wrap;
  word-break: break-word;
}

.muted {
  color: #64748b;
}

.error {
  color: #ef4444;
}

.ghost {
  border: 1px solid rgba(15, 23, 42, 0.18);
  background: white;
  padding: 0.45rem 0.9rem;
  border-radius: 999px;
  cursor: pointer;
  font-weight: 600;
}

.btn {
  text-decoration: none;
  background: #0f172a;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 999px;
  box-shadow: 0 10px 20px rgba(15, 23, 42, 0.25);
}

@media (max-width: 1100px) {
  .grid {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
  .card--span-8 { grid-column: span 6; }
}

@media (max-width: 780px) {
  .grid {
    grid-template-columns: repeat(1, minmax(0, 1fr));
  }
  .card--span-4,
  .card--span-6,
  .card--span-8,
  .card--span-12 {
    grid-column: span 1;
  }
}

@keyframes fadeUp {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>








