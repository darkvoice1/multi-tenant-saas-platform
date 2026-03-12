export type UUID = string

export type Project = {
  id: UUID
  tenant_id: UUID
  org_id?: UUID | null
  name: string
  description: string
  created_by?: UUID | null
  created_at: string
  updated_at: string
}

export type TaskStatus = 'todo' | 'in_progress' | 'review' | 'done' | 'rejected'
export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent'

export type Task = {
  id: UUID
  tenant_id: UUID
  project_id: UUID
  title: string
  status: TaskStatus
  assignee_id?: UUID | null
  priority: TaskPriority
  due_at?: string | null
  created_at: string
  updated_at: string
}

export type TaskComment = {
  id: UUID
  tenant_id: UUID
  task_id: UUID
  user_id: UUID
  content: string
  created_at: string
  updated_at: string
}

export type TaskAttachment = {
  id: UUID
  tenant_id: UUID
  task_id: UUID
  uploader_id: UUID
  file_name: string
  content_type: string
  size_bytes: number
  created_at: string
  updated_at: string
}

export type Notification = {
  id: UUID
  tenant_id: UUID
  user_id: UUID
  type: string
  message: string
  read_at?: string | null
  created_at: string
  updated_at: string
}

export type DashboardStatusCount = {
  status: string
  count: number
}

export type DashboardSummary = {
  tenant: { id: UUID; name: string; slug: string; status: string }
  metrics: {
    project_count: number
    task_count: number
    task_status_counts: DashboardStatusCount[]
    pending_review_task_count: number
    unread_notification_count: number
    my_open_task_count: number
  }
  lists: {
    due_soon_tasks: Array<{
      id: UUID
      project_id: UUID
      title: string
      status: string
      priority: string
      due_at?: string | null
      updated_at: string
    }>
    my_due_soon_tasks: Array<{
      id: UUID
      project_id: UUID
      title: string
      status: string
      priority: string
      due_at?: string | null
      updated_at: string
    }>
    recent_tasks: Array<{
      id: UUID
      project_id: UUID
      title: string
      status: string
      priority: string
      due_at?: string | null
      updated_at: string
    }>
    pending_review_tasks: Array<{
      id: UUID
      project_id: UUID
      title: string
      status: string
      priority: string
      due_at?: string | null
      updated_at: string
    }>
    recent_comments: Array<{
      id: UUID
      task_id: UUID
      content: string
      created_at: string
      user_name: string
      user_email: string
    }>
    recent_notifications: Array<{
      id: UUID
      type: string
      message: string
      read_at?: string | null
      created_at: string
    }>
    storage_backend: string
    storage_s3_bucket: string
  }
}
