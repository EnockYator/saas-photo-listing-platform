CREATE TABLE IF NOT EXISTS listing_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    
    position INT NOT NULL,  -- order of the photo in the listing
    is_cover BOOLEAN NOT NULL DEFAULT FALSE,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT chk_listing_photos_timestamps
        CHECK (updated_at >= created_at),

    -- Constraints
    -- One file appears once per listing
    UNIQUE (listing_id, file_id),

    -- Prevent duplicate positions within a listing
    UNIQUE (listing_id, position)

);

CREATE TRIGGER trg_listing_photos_updated_at
BEFORE UPDATE ON listing_photos
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes for performance
CREATE INDEX idx_listing_photos_tenant_listing
    ON listing_photos(tenant_id, listing_id);

CREATE INDEX idx_listing_photos_cover
    ON listing_photos(listing_id)
    WHERE is_cover = TRUE;

CREATE INDEX idx_listing_photos_published
    ON listing_photos(listing_id)
    WHERE is_published = TRUE;