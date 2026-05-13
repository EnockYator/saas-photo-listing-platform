CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Ensure user cannot have the same role multiple times in the same tenant
    CONSTRAINT uq_user_role_tenant UNIQUE (user_id, role_id, tenant_id)
);

-- Quickly list roles per tenant
CREATE INDEX idx_user_roles_tenant
    ON user_roles(tenant_id);

-- Quickly look up roles for a user within a tenant
CREATE INDEX idx_user_roles_tenant_user
    ON user_roles(user_id, tenant_id);
