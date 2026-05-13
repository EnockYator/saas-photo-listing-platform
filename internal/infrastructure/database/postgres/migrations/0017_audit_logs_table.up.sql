
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,

    entity_id UUID NOT NULL,
    
    entity_type TEXT NOT NULL
        CONSTRAINT audit_logs_entity_type_check CHECK (entity_type IN (
            'user',
            'listing',
            'file',
            'subscription',
            'plan',
            'notification',
            'tenant'
        )),
    
    action TEXT NOT NULL
        CONSTRAINT audit_logs_action_check CHECK (action IN (
            'create',
            'update',
            'delete',
            'login',
            'logout',
            'publish',
            'unpublish',
            'share',
            'other'
        )),

    
    changed_data JSONB,
    
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_tenant_time
    ON audit_logs(tenant_id, performed_at DESC);

CREATE INDEX idx_audit_logs_performed_by
    ON audit_logs(performed_by);

CREATE INDEX idx_audit_logs_entity 
    ON audit_logs(entity_id, entity_type);