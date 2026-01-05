CREATE TABLE plan_limits (
    plan_id UUID PRIMARY KEY REFERENCES plans(id) ON DELETE CASCADE,

    max_storage_bytes BIGINT NOT NULL,
    max_upload_bytes BIGINT NOT NULL,

    CHECK (max_storage_bytes >= 0),
    CHECK (max_upload_bytes >= 0),
    CHECK (max_upload_bytes <= max_storage_bytes)
);
