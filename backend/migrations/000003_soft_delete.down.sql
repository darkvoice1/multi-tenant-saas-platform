DROP INDEX IF EXISTS idx_tasks_deleted_at;
DROP INDEX IF EXISTS idx_projects_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_orgs_deleted_at;
DROP INDEX IF EXISTS idx_tenants_deleted_at;

ALTER TABLE tasks DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE projects DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE orgs DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE tenants DROP COLUMN IF EXISTS deleted_at;
