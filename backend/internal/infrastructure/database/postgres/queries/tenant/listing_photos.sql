-- name: AddListingPhoto :one
INSERT INTO listing_photos (tenant_id, listing_id, file_id, position, is_cover, is_published)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListListingPhotos :many
SELECT *
FROM listing_photos
WHERE tenant_id = $1
  AND listing_id = $2
  AND deleted_at IS NULL
ORDER BY position ASC;

-- name: SetCoverPhoto :exec
UPDATE listing_photos
SET is_cover = CASE WHEN id = $3 THEN TRUE ELSE FALSE END,
    updated_at = NOW()
WHERE tenant_id = $1
  AND listing_id = $2;

-- name: SoftDeleteListingPhoto :one
UPDATE listing_photos
SET deleted_at = NOW(), updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
RETURNING id;
