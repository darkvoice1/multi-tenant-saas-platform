ALTER TABLE tenants ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
ALTER TABLE orgs ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deleted_at timestamptz;

CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants(deleted_at);
CREATE INDEX IF NOT EXISTS idx_orgs_deleted_at ON orgs(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
