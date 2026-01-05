CRAETE TYPE share_permission AS ENUM ('read', 'write', 'admin');

CREATE TABLE IF NOT EXISTS share_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    permission share_permission NOT NULL,  -- e.g., 'read', 'write', 'admin'

    token TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,
    max_views INT NOT NULL DEFAULT 0,  -- 0 means unlimited views
    view_count INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Safety Constarints
    CHECK (max_views >= 0),
        CHECK (view_count >= 0),
        CHECK (view_count <= max_views OR max_views = 0)

);

-- 3️⃣ Trigger bound to the table
CREATE TRIGGER trg_share_links_updated_at
BEFORE UPDATE ON share_links
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_share_links_listing
    ON share_links(listing_id);
