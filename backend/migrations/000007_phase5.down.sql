DROP TABLE IF EXISTS audit_logs;

ALTER TABLE tenants DROP COLUMN IF EXISTS max_requests_per_minute;
ALTER TABLE tenants DROP COLUMN IF EXISTS max_storage_bytes;
ALTER TABLE tenants DROP COLUMN IF EXISTS max_members;
ALTER TABLE tenants DROP COLUMN IF EXISTS max_projects;
