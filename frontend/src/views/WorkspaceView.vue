<template>
  <section class="page">
    <header class="header">
      <div>
        <h2>协作工作台</h2>
        <p class="muted">项目、任务、评论、审批与附件都可以在这里完成。</p>
      </div>
      <button class="ghost" @click="loadAll" :disabled="loading">刷新数据</button>
    </header>

    <div class="grid">
      <div class="card">
        <h3>项目</h3>
        <form class="form" @submit.prevent="createProject">
          <input v-model="projectForm.name" placeholder="项目名称" required />
          <input v-model="projectForm.description" placeholder="项目描述" />
          <button :disabled="loading || !canProjectWrite">创建项目</button>
        </form>
        <ul class="list">
          <li v-for="p in projects" :key="p.id" :class="{ active: p.id === selectedProjectId }">
            <button class="link" @click="selectProject(p.id)">{{ p.name }}</button>
            <span class="muted">{{ p.description || '无描述' }}</span>
          </li>
        </ul>
      </div>

      <div class="card">
        <h3>任务</h3>
        <form class="form" @submit.prevent="createTask" v-if="selectedProjectId">
          <input v-model="taskForm.title" placeholder="任务标题" required />
          <select v-model="taskForm.priority">
            <option value="medium">中</option>
            <option value="low">低</option>
            <option value="high">高</option>
            <option value="urgent">紧急</option>
          </select>
          <input v-model="taskForm.dueAtLocal" type="datetime-local" />
          <button :disabled="loading || !canTaskWrite">创建任务</button>
        </form>
        <div v-else class="muted">先选择一个项目。</div>
        <ul class="list">
          <li v-for="t in tasks" :key="t.id" :class="{ active: t.id === selectedTaskId }">
            <button class="link" @click="selectTask(t.id)">{{ t.title }}</button>
            <span class="badge">{{ t.status }}</span>
            <span class="muted">{{ t.priority }}</span>
          </li>
        </ul>
        <div v-if="selectedTask" class="section">
          <h4>任务状态</h4>
          <div class="row">
            <select v-model="statusForm.status">
              <option value="todo">待办</option>
              <option value="in_progress">进行中</option>
              <option value="review">待审批</option>
              <option value="done">已完成</option>
              <option value="rejected">已拒绝</option>
            </select>
            <button @click="updateStatus" :disabled="loading || !canTaskWrite">更新状态</button>
          </div>
          <div class="row">
            <button class="ghost" @click="approveTask('approved')" :disabled="loading || !canApprove">审批通过</button>
            <button class="ghost" @click="approveTask('rejected')" :disabled="loading || !canApprove">审批拒绝</button>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>评论与附件</h3>
        <div v-if="selectedTask">
          <form class="form" @submit.prevent="createComment">
            <input v-model="commentForm.content" placeholder="评论内容，支持@邮箱" required />
            <button :disabled="loading || !canTaskWrite">发送评论</button>
          </form>
          <ul class="list">
            <li v-for="c in comments" :key="c.id">{{ c.content }}</li>
          </ul>

          <form class="form" @submit.prevent="uploadAttachment">
            <input ref="fileInput" type="file" />
            <button :disabled="loading || !canTaskWrite">上传附件</button>
          </form>
          <ul class="list">
            <li v-for="a in attachments" :key="a.id">
              <div class="att">
                <span>{{ a.file_name }} ({{ a.size_bytes }} bytes)</span>
                <div class="att-actions">
                  <button
                    class="ghost"
                    @click="previewAttachment(a)"
                    :disabled="loading || !a.content_type.toLowerCase().startsWith('image/')"
                  >
                    预览
                  </button>
                  <button class="ghost" @click="downloadAttachment(a)" :disabled="loading">下载</button>
                </div>
              </div>
            </li>
          </ul>
        </div>
        <div v-else class="muted">先选择一个任务。</div>
      </div>

      <div class="card">
        <h3>通知</h3>
        <button class="ghost" @click="loadNotifications" :disabled="loading">刷新通知</button>
        <ul class="list">
          <li v-for="n in notifications" :key="n.id">
            <span>{{ n.message }}</span>
            <button class="ghost" @click="markRead(n.id)" :disabled="loading">标记已读</button>
          </li>
        </ul>
      </div>
    </div>

    <div v-if="previewUrl" class="modal" @click.self="closePreview">
      <div class="modal-card">
        <header class="modal-head">
          <strong>预览：{{ previewName }}</strong>
          <button class="ghost" @click="closePreview">关闭</button>
        </header>
        <img class="preview" :src="previewUrl" alt="preview" />
      </div>
    </div>

    <p v-if="error" class="error">{{ error }}</p>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { apiFetch, apiFetchBlob } from '../api'
import { useSessionStore } from '../stores/session'
import { isAllowed } from '../permissions'
import type { Notification, Project, Task, TaskAttachment, TaskComment, TaskPriority, TaskStatus, UUID } from '../types'

const loading = ref(false)
const error = ref('')

const session = useSessionStore()
const route = useRoute()

const projects = ref<Project[]>([])
const tasks = ref<Task[]>([])
const comments = ref<TaskComment[]>([])
const attachments = ref<TaskAttachment[]>([])
const notifications = ref<Notification[]>([])

const selectedProjectId = ref<UUID | null>(null)
const selectedTaskId = ref<UUID | null>(null)

const projectForm = reactive({ name: '', description: '' })
const taskForm = reactive<{ title: string; priority: TaskPriority; dueAtLocal: string }>({
  title: '',
  priority: 'medium',
  dueAtLocal: ''
})
const statusForm = reactive<{ status: TaskStatus }>({ status: 'todo' })
const commentForm = reactive({ content: '' })

const fileInput = ref<HTMLInputElement | null>(null)
const previewUrl = ref<string | null>(null)
const previewName = ref<string>('')

const selectedTask = computed(() => tasks.value.find((t) => t.id === selectedTaskId.value))
const canProjectWrite = computed(() => isAllowed(session.role, 'project:write'))
const canTaskWrite = computed(() => isAllowed(session.role, 'task:write'))
const canApprove = computed(() => session.role === 'admin' || session.role === 'manager')

async function loadProjects() {
  projects.value = (await apiFetch<Project[]>('/api/projects')) || []
}

async function loadTasks(projectId: string) {
  tasks.value = (await apiFetch<Task[]>(`/api/projects/${projectId}/tasks`)) || []
}

async function loadComments(taskId: string) {
  comments.value = (await apiFetch<TaskComment[]>(`/api/tasks/${taskId}/comments`)) || []
}

async function loadAttachments(taskId: string) {
  attachments.value = (await apiFetch<TaskAttachment[]>(`/api/tasks/${taskId}/attachments`)) || []
}

async function loadNotifications() {
  notifications.value = (await apiFetch<Notification[]>('/api/notifications')) || []
}

async function loadAll() {
  await withLoading(async () => {
    await loadProjects()
    if (selectedProjectId.value) await loadTasks(selectedProjectId.value)
    if (selectedTaskId.value) {
      await loadComments(selectedTaskId.value)
      await loadAttachments(selectedTaskId.value)
    }
    await loadNotifications()
  })
}

async function createProject() {
  if (!canProjectWrite.value) {
    error.value = '当前角色无项目写权限'
    return
  }
  await withLoading(async () => {
    const body = JSON.stringify(projectForm)
    const proj = await apiFetch<Project>('/api/projects', { method: 'POST', body })
    projects.value.unshift(proj)
    projectForm.name = ''
    projectForm.description = ''
  })
}

async function selectProject(id: string) {
  selectedProjectId.value = id
  await withLoading(async () => {
    await loadTasks(id)
    selectedTaskId.value = null
    comments.value = []
    attachments.value = []
  })
}

async function createTask() {
  if (!selectedProjectId.value) return
  if (!canTaskWrite.value) {
    error.value = '当前角色无任务写权限'
    return
  }
  await withLoading(async () => {
    const payload: any = {
      title: taskForm.title,
      priority: taskForm.priority
    }
    if (taskForm.dueAtLocal) payload.due_at = new Date(taskForm.dueAtLocal).toISOString()
    const task = await apiFetch<Task>(`/api/projects/${selectedProjectId.value}/tasks`, {
      method: 'POST',
      body: JSON.stringify(payload)
    })
    tasks.value.unshift(task)
    taskForm.title = ''
    taskForm.dueAtLocal = ''
  })
}

async function selectTask(id: string) {
  selectedTaskId.value = id
  await withLoading(async () => {
    await loadComments(id)
    await loadAttachments(id)
  })
}

async function updateStatus() {
  if (!selectedTaskId.value) return
  if (!canTaskWrite.value) {
    error.value = '当前角色无任务写权限'
    return
  }
  await withLoading(async () => {
    const body = JSON.stringify({ status: statusForm.status })
    const updated = await apiFetch<Task>(`/api/tasks/${selectedTaskId.value}/status`, { method: 'POST', body })
    tasks.value = tasks.value.map((t) => (t.id === updated.id ? updated : t))
  })
}

async function approveTask(status: string) {
  if (!selectedTaskId.value) return
  if (!canApprove.value) {
    error.value = '当前角色无审批权限'
    return
  }
  await withLoading(async () => {
    await apiFetch(`/api/tasks/${selectedTaskId.value}/approve`, {
      method: 'POST',
      body: JSON.stringify({ status, comment: '' })
    })
    await loadTasks(selectedProjectId.value as string)
  })
}

async function createComment() {
  if (!selectedTaskId.value) return
  if (!canTaskWrite.value) {
    error.value = '当前角色无任务写权限'
    return
  }
  await withLoading(async () => {
    const body = JSON.stringify({ content: commentForm.content })
    const comment = await apiFetch<TaskComment>(`/api/tasks/${selectedTaskId.value}/comments`, { method: 'POST', body })
    comments.value.unshift(comment)
    commentForm.content = ''
    await loadNotifications()
  })
}

async function uploadAttachment() {
  if (!selectedTaskId.value) return
  if (!canTaskWrite.value) {
    error.value = '当前角色无任务写权限'
    return
  }
  const input = fileInput.value
  if (!input || !input.files || input.files.length === 0) {
    error.value = '请选择文件'
    return
  }
  await withLoading(async () => {
    const form = new FormData()
    form.append('file', input.files![0])
    const item = await apiFetch<TaskAttachment>(`/api/tasks/${selectedTaskId.value}/attachments`, {
      method: 'POST',
      body: form
    })
    attachments.value.unshift(item)
    input.value = ''
  })
}

async function previewAttachment(att: TaskAttachment) {
  await withLoading(async () => {
    const blob = await apiFetchBlob(`/api/attachments/${att.id}/preview`)
    if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = URL.createObjectURL(blob)
    previewName.value = att.file_name
  })
}

function closePreview() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
  previewUrl.value = null
  previewName.value = ''
}

async function downloadAttachment(att: TaskAttachment) {
  await withLoading(async () => {
    const blob = await apiFetchBlob(`/api/attachments/${att.id}/download`)
    const url = URL.createObjectURL(blob)
    try {
      const a = document.createElement('a')
      a.href = url
      a.download = att.file_name || 'attachment'
      document.body.appendChild(a)
      a.click()
      a.remove()
    } finally {
      URL.revokeObjectURL(url)
    }
  })
}

async function markRead(id: string) {
  await withLoading(async () => {
    await apiFetch(`/api/notifications/${id}/read`, { method: 'POST' })
    await loadNotifications()
  })
}

async function withLoading(fn: () => Promise<void>) {
  loading.value = true
  error.value = ''
  try {
    await fn()
  } catch (err: any) {
    error.value = err.message || '操作失败'
  } finally {
    loading.value = false
  }
}

function queryString(val: unknown): string {
  if (typeof val === 'string') return val
  if (Array.isArray(val) && typeof val[0] === 'string') return val[0]
  return ''
}

onMounted(() => {
  // Allow deep-linking from the dashboard (e.g. /workspace?task=...).
  const qTask = queryString(route.query.task)
  const qProject = queryString(route.query.project)

  withLoading(async () => {
    await loadProjects()

    if (qTask) {
      const task = await apiFetch<Task>(`/api/tasks/${qTask}`)
      selectedProjectId.value = task.project_id
      await loadTasks(task.project_id)
      selectedTaskId.value = task.id
      await loadComments(task.id)
      await loadAttachments(task.id)
      await loadNotifications()
      return
    }

    if (qProject) {
      selectedProjectId.value = qProject
      await loadTasks(qProject)
    }

    await loadNotifications()
  })
})
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.grid {
  display: grid;
  gap: 1.2rem;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1rem;
  background: #fff;
}

.form {
  display: grid;
  gap: 0.6rem;
  margin-bottom: 0.8rem;
}

.list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 0.6rem;
}

.list li {
  display: grid;
  gap: 0.3rem;
}

.list li.active {
  background: #f1f5f9;
  padding: 0.4rem;
  border-radius: 8px;
}

.link {
  border: none;
  background: none;
  padding: 0;
  text-align: left;
  font-weight: 600;
  cursor: pointer;
}

.row {
  display: flex;
  gap: 0.6rem;
  margin-bottom: 0.6rem;
}

.att {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
}

.att-actions {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
}

.badge {
  font-size: 0.8rem;
  background: #e2e8f0;
  padding: 0.1rem 0.4rem;
  border-radius: 6px;
  display: inline-block;
}

button {
  padding: 0.5rem 0.8rem;
  border-radius: 8px;
  border: none;
  background: #0f172a;
  color: white;
  cursor: pointer;
}

button.ghost {
  background: transparent;
  color: #0f172a;
  border: 1px solid #cbd5f5;
}

.muted {
  color: #64748b;
  font-size: 0.9rem;
}

.error {
  color: #ef4444;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.5);
  display: grid;
  place-items: center;
  padding: 2rem;
}

.modal-card {
  width: min(920px, 100%);
  background: white;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.8rem 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.preview {
  display: block;
  width: 100%;
  max-height: 70vh;
  object-fit: contain;
  background: #0b1220;
}
</style>
