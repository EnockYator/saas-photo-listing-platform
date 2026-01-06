
CREATE TABLE IF NOT EXISTS share_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
<<<<<<< HEAD
=======
    
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
>>>>>>> c431e38 (refactor: place all miration files in migrations/ directory rather than sub-directories)

    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    permission TEXT NOT NULL
        CONSTRAINT share_links_permission_check CHECK (permission IN ('read', 'write', 'admin')),

    token TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,
    max_views INT NOT NULL DEFAULT 0,  -- 0 means unlimited views
    view_count INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_share_links_timestamps
        CHECK (updated_at >= created_at),

    -- Safety Constarints
    CONSTRAINT max_views_positive_check
        CHECK (max_views >= 0),
    
    CONSTRAINT view_counts_positive_check
        CHECK (view_count >= 0),

    CONSTRAINT view_counts_not_exceeding_max_views_check
        CHECK (view_count <= max_views OR max_views = 0)
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_share_links_updated_at
BEFORE UPDATE ON share_links
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes
CREATE INDEX idx_share_links_listing
    ON share_links(listing_id);

CREATE INDEX idx_share_links_token
    ON share_links(token);
