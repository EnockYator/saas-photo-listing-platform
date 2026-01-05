
CREATE TYPE listing_status as ENUM ('draft', 'published');
CREATE TYPE listing_visibility as ENUM ('private', 'public');

CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        
    title VARCHAR(200) NOT NULL,
    description TEXT,

    status listing_status NOT NULL DEFAULT 'draft',
    visibility listing_visibility NOT NULL DEFAULT 'private',

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
