
-- Create listings table
CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        
    title VARCHAR(200) NOT NULL,
    description TEXT,

    status TEXT NOT NULL DEFAULT 'draft'
        CONSTRAINT listings_listing_status_check CHECK (status IN ('draft', 'published')),

    visibility TEXT NOT NULL DEFAULT 'private'
        CONSTRAINT listings_listing_visibility_check CHECK (visibility IN ('private', 'public')),


    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ  DEFAULT NULL
);

-- Attach shared trigger
CREATE TRIGGER trg_listings_updated_at
BEFORE UPDATE ON listings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Indexes (query-driven)
CREATE INDEX idx_listings_tenant_user
    ON listings(tenant_id, user_id);

CREATE INDEX idx_listings_tenant_status
    ON listings(tenant_id, status);

CREATE INDEX idx_listings_tenant_visibility
    ON listings(tenant_id, visibility);
