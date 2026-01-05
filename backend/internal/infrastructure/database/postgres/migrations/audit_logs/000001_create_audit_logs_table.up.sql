
CREATE TYPE audit_action AS ENUM (
    'CREATE',
    'UPDATE',
    'DELETE',
    'LOGIN',
    'LOGOUT',
    'PUBLISH',
    'UNPUBLISH',
    'SHARE'
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- User who performed the action
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Target entity of the action  / What was affected
    entity_id UUID NOT NULL,
    entity_type TEXT NOT NULL, --- e.g., "photo", "listing", "user", "subscription", "tenant" etc.
    
    -- What happened
    action audit_action NOT NULL,  --- e.g., "create", "update", "delete", "login", "logout" etc.
    
    -- What changed
    changed_data JSONB,
    
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_tenant_time
    ON audit_logs(tenant_id, performed_at DESC);

CREATE INDEX idx_audit_logs_performed_by
    ON audit_logs(performed_by);

CREATE INDEX idx_audit_logs_entity 
    ON audit_logs(entity, entity_type);