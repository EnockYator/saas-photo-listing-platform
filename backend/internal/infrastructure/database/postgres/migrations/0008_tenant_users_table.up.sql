CREATE TABLE IF NOT EXISTS tenant_users (

<<<<<<< HEAD
=======
CREATE TABLE IF NOT EXISTS tenant_users (
>>>>>>> e5eedb5 (chore(database): add indexes and constraints to database tables for fast performance and security)
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    role TEXT NOT NULL DEFAULT 'viewer'
        CONSTRAINT tenant_user_role_check CHECK (role IN ('admin', 'editor', 'viewer')),
<<<<<<< HEAD
 
=======
    
>>>>>>> e5eedb5 (chore(database): add indexes and constraints to database tables for fast performance and security)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    PRIMARY KEY (tenant_id, user_id)
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_tenant_users_updated_at
BEFORE UPDATE ON tenant_users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes
CREATE INDEX idx_tenant_user_user_id
    ON tenant_users(user_id);

CREATE INDEX idx_tenant_user_role
    ON tenant_users(role);
