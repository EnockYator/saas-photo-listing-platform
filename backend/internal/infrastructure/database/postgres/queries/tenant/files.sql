-- name: CreateFile :one
INSERT INTO files (tenant_id, listing_id, user_id, original_url, watermarked_url, watermark_type, thumbnail_url, file_size_bytes, mime_type)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListFilesByListing :many
SELECT *
FROM files
WHERE tenant_id = $1
  AND listing_id = $2
ORDER BY created_at DESC;

-- name: ListFilesByUser :many
SELECT *
FROM files
WHERE tenant_id = $1
  AND user_id = $2
ORDER BY created_at DESC;

