-- name: CreatePlanLimits :one
INSERT INTO plan_limits (plan_id, max_storage_bytes, max_upload_bytes, max_listings, max_listing_photos)
VALUES ($1, $2, $3, $4, $5)
RETURNING plan_id, max_storage_bytes, max_upload_bytes, max_listings, max_listing_photos;

-- name: GetPlanLimitsByPlan :one
SELECT plan_id, max_storage_bytes, max_upload_bytes, max_listings, max_listing_photos
FROM plan_limits
WHERE plan_id = $1;

-- name: UpdatePlanLimits :one
UPDATE plan_limits
SET max_storage_bytes = $2,
    max_upload_bytes = $3,
    max_listings = $4,
    max_listing_photos = $5
WHERE plan_id = $1
RETURNING plan_id, max_storage_bytes, max_upload_bytes, max_listings, max_listing_photos;

