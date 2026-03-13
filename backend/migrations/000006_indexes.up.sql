CREATE INDEX IF NOT EXISTS idx_tasks_tenant_status ON tasks(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_due_at ON tasks(tenant_id, due_at);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_assignee ON tasks(tenant_id, assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_tenant_updated_at ON tasks(tenant_id, updated_at);

CREATE INDEX IF NOT EXISTS idx_task_comments_tenant_created_at ON task_comments(tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_read_at ON notifications(tenant_id, read_at);
