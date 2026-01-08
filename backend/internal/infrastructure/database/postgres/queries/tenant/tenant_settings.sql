-- name: CreateTenantSettings :one
INSERT INTO tenant_settings (tenant_id, theme, watermark_enabled, watermark_text, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: GetTenantSettings :one
SELECT *
FROM tenant_settings
WHERE tenant_id = $1;

-- name: UpdateTenantSettings :one
UPDATE tenant_settings
SET theme = $2,
    watermark_enabled = $3,
    watermark_text = $4,
    updated_at = NOW()
WHERE tenant_id = $1
RETURNING *;
