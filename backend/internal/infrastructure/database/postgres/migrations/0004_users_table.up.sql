
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    
    role TEXT NOT NULL DEFAULT 'viewer'
        CONSTRAINT users_role_check CHECK (role IN ('admin', 'editor', 'viewer')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Recommended indexes
CREATE INDEX idx_users_tenant_id
    ON users(tenant_id);

CREATE INDEX idx_users_role
    ON users(role);
