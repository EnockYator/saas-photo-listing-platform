-- name: CreateShareLink :one
INSERT INTO share_links (
    listing_id, tenant_id, permission, token, expires_at, max_views
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, listing_id, tenant_id, permission, token, expires_at, max_views, view_count, created_at, updated_at;

-- name: GetShareLinkByToken :one
SELECT *
FROM share_links
WHERE token = $1
  AND tenant_id = $2
  AND (expires_at > NOW());

-- name: IncrementShareLinkView :one
UPDATE share_links
SET view_count = view_count + 1, updated_at = NOW()
WHERE id = $1
  AND (view_count < max_views OR max_views = 0)
RETURNING id, view_count;

-- name: ListShareLinksByListing :many
SELECT id, permission, token, expires_at, max_views, view_count, created_at, updated_at
FROM share_links
WHERE tenant_id = $1
  AND listing_id = $2
ORDER BY created_at DESC;

-- name: DeleteShareLink :exec
DELETE FROM share_links
WHERE tenant_id = $1
  AND id = $2;
