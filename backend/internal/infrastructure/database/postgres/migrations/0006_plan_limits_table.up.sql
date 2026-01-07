CREATE TABLE plan_limits (
    plan_id UUID PRIMARY KEY REFERENCES plans(id) ON DELETE CASCADE,

    max_storage_bytes BIGINT NOT NULL
        CONSTRAINT plan_limits_max_storage_bytes_check CHECK (max_storage_bytes >= 0),

    max_upload_bytes BIGINT NOT NULL
        CONSTRAINT plan_limits_max_upload_bytes_check CHECK (max_upload_bytes >= 0),

    max_listings INT NOT NULL
        CONSTRAINT plan_limits_max_listings_check CHECK (max_listings >= 0),

    max_listing_photos INT NOT NULL
        CONSTRAINT plan_limits_max_listing_photos_check CHECK (max_listing_photos >= 0),

    CONSTRAINT max_upload_bytes_not_exceeding_max_storage_bytes CHECK (max_upload_bytes <= max_storage_bytes)
);
