CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_tenant_users_username UNIQUE (tenant_id, username),
    CONSTRAINT uq_tenant_users_email UNIQUE (tenant_id, email),

    CONSTRAINT chk_users_username_length
        CHECK (length(username) BETWEEN 3 AND 50),

    CONSTRAINT chk_users_username_format
        CHECK (username ~ '^[a-zA-Z0-9_]+$'),

    CONSTRAINT chk_users_email_format
        CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),

    CONSTRAINT chk_users_password_hash_length
        CHECK (length(password_hash) >= 60),

    CONSTRAINT chk_users_timestamps
        CHECK (updated_at >= created_at)
);


-- updated_at trigger
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes
CREATE INDEX idx_users_tenant_id
    ON users (tenant_id);

CREATE INDEX idx_users_created_at_desc
    ON users (created_at DESC);
