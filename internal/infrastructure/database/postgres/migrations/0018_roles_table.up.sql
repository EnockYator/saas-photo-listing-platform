CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    name TEXT NOT NULL DEFAULT 'viewer'
        CONSTRAINT roles_name_check CHECK (name IN ('super_admin', 'tenant_admin', 'tenant_editor', 'viewer')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

