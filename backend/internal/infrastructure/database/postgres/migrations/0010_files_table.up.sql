
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    original_url TEXT NOT NULL,
    watermarked_url TEXT,
    
    watermark_type TEXT
        CONSTRAINT files_watermark_type_check CHECK (watermark_type IN ('text', 'image')),
    
    thumbnail_url TEXT,
    
    file_size_bytes BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_file_timestamps
        CHECK (updated_at >= created_at),

    CHECK (
            (watermark_type IS NULL AND watermarked_url IS NULL)
        OR (watermark_type IS NOT NULL AND watermarked_url IS NOT NULL)
        )
);

CREATE TRIGGER trg_files_updated_at
BEFORE UPDATE ON files
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes (query-driven)
CREATE INDEX idx_files_tenant_listing
    ON files(tenant_id, listing_id);

CREATE INDEX idx_files_user
    ON files(user_id);

CREATE INDEX idx_files_mime_type
    ON files(mime_type);