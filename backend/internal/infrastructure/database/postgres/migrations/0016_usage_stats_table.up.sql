CREATE TABLE IF NOT EXISTS usage_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    total_uploads BIGINT NOT NULL DEFAULT 0,
    total_storage_used_bytes BIGINT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_usage_stats_timestamps
        CHECK (updated_at >= created_at),

    CONSTRAINT chk_total_uploads_positive
        CHECK (total_uploads >= 0),

    CONSTRAINT chk_total_storage_used_bytes_positive
        CHECK (total_storage_used_bytes >= 0),

    UNIQUE (tenant_id, user_id)  -- one row per tenant+user
);

-- Trigger for updated_at
CREATE TRIGGER trg_usage_stats_updated_at
BEFORE UPDATE ON usage_stats
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Recommended indexes
CREATE INDEX idx_usage_stats_tenant_user
    ON usage_stats(tenant_id, user_id);

CREATE INDEX idx_usage_stats_updated_at
    ON usage_stats(updated_at);