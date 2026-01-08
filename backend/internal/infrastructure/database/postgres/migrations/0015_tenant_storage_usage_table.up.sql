CREATE TABLE IF NOT EXISTS tenant_storage_usage (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,

    used_storage_bytes BIGINT NOT NULL DEFAULT 0
        CONSTRAINT used_storage_bytes_positive_check CHECK (used_storage_bytes >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_tenant_storage_usage_timestamps
    CHECK (updated_at >= created_at)

);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_tenant_storage_usage_updated_at
BEFORE UPDATE ON tenant_storage_usage
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Index for tracking by update time
CREATE INDEX idx_tenant_storage_usage_updated_at
    ON tenant_storage_usage(updated_at);
