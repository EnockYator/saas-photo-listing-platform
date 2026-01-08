-- name: CreateUsageStats :one
INSERT INTO usage_stats (tenant_id, user_id, total_uploads, total_storage_used_bytes, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING id, tenant_id, user_id, total_uploads, total_storage_used_bytes, created_at, updated_at;

-- name: GetUsageStatsByUser :one
SELECT tenant_id, user_id, total_uploads, total_storage_used_bytes, created_at, updated_at
FROM usage_stats
WHERE tenant_id = $1
  AND user_id = $2;

-- name: IncrementUsageStats :one
UPDATE usage_stats
SET total_uploads = total_uploads + $3,
    total_storage_used_bytes = total_storage_used_bytes + $4
WHERE tenant_id = $1
  AND user_id = $2
RETURNING tenant_id, user_id, total_uploads, total_storage_used_bytes, updated_at;

-- name: ResetUsageStats :one
UPDATE usage_stats
SET total_uploads = 0,
    total_storage_used_bytes = 0
WHERE tenant_id = $1
  AND user_id = $2
RETURNING tenant_id, user_id, total_uploads, total_storage_used_bytes, updated_at;

-- name: ListTopUsersByStorage :many
SELECT user_id, total_uploads, total_storage_used_bytes
FROM usage_stats
WHERE tenant_id = $1
ORDER BY total_storage_used_bytes DESC
LIMIT $2;
