-- name: CreateTenantStorageUsage :one
INSERT INTO tenant_storage_usage (tenant_id, used_storage_bytes, created_at)
VALUES ($1, 0, NOW())
RETURNING *;

-- name: IncrementTenantStorageUsage :one
UPDATE tenant_storage_usage
SET used_storage_bytes = used_storage_bytes + $2,
    updated_at = NOW()
WHERE tenant_id = $1
RETURNING used_storage_bytes;

-- name: DecrementTenantStorageUsage :one
UPDATE tenant_storage_usage
SET used_storage_bytes = GREATEST(used_storage_bytes - $2, 0),
    updated_at = NOW()
WHERE tenant_id = $1
RETURNING used_storage_bytes;

-- name: GetTenantStorageUsage :one
SELECT used_storage_bytes, created_at, updated_at
FROM tenant_storage_usage
WHERE tenant_id = $1;

