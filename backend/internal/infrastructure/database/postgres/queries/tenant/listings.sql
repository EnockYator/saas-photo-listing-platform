-- name: CreateListing :one
INSERT INTO listings (tenant_id, user_id, title, description, status, visibility, created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING id, tenant_id, user_id, title, description, status, visibility, created_at, updated_at, deleted_at;

-- name: GetListingByID :one
SELECT *
FROM listings
WHERE tenant_id = $1
  AND id = $2
  AND deleted_at IS NULL;

-- name: ListListingsByTenantUser :many
SELECT *
FROM listings
WHERE tenant_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3;

-- name: UpdateListing :one
UPDATE listings
SET status = $3, visibility = $4, updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
RETURNING *;

-- name: SoftDeleteListing :one
UPDATE listings
SET deleted_at = NOW(), updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
RETURNING id;


