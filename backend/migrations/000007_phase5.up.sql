ALTER TABLE tenants ADD COLUMN IF NOT EXISTS max_projects int NOT NULL DEFAULT 100;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS max_members int NOT NULL DEFAULT 100;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS max_storage_bytes bigint NOT NULL DEFAULT 1073741824;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS max_requests_per_minute int NOT NULL DEFAULT 600;

CREATE TABLE IF NOT EXISTS audit_logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action text NOT NULL,
    resource text NOT NULL,
    resource_id uuid,
    method text NOT NULL,
    path text NOT NULL,
    status_code int NOT NULL,
    ip text,
    user_agent text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_created_at ON audit_logs(tenant_id, created_at DESC);
