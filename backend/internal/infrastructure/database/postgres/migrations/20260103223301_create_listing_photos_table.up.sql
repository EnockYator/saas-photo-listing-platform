CREATE TABLE IF NOT EXISTS listing_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL,
    
    url TEXT NOT NULL,
    position INT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_listing_photos_listing
    FOREIGN KEY (listing_id) REFERENCES listings(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_listing_photos_listing_id ON listing_photos(listing_id);
CREATE INDEX idx_listing_photos_is_primary ON listing_photos(is_primary);
CREATE INDEX idx_listing_photos_id ON listing_photos(id);
