-- name: CreateTenant :one
INSERT INTO tenants (name, created_at)
VALUES ($1, NOW())
RETURNING id, name, created_at, updated_at;

-- name: GetTenantByID :one
SELECT id, name, created_at, updated_at
FROM tenants
WHERE id = $1;

-- name: ListTenants :many
SELECT id, name, created_at, updated_at
FROM tenants
ORDER BY created_at DESC
LIMIT $1;

-- name: UpdateTenantName :one
UPDATE tenants
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, name, created_at, updated_at;

