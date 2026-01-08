-- name: AddTenantUser :one
INSERT INTO tenant_users (tenant_id, user_id, role, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING tenant_id, user_id, role, created_at, updated_at;

-- name: UpdateTenantUserRole :one
UPDATE tenant_users
SET role = $3, updated_at = NOW()
WHERE tenant_id = $1
  AND user_id = $2
RETURNING tenant_id, user_id, role, created_at, updated_at;

-- name: RemoveTenantUser :exec
DELETE FROM tenant_users
WHERE tenant_id = $1
  AND user_id = $2;

-- name: ListTenantUsers :many
SELECT user_id, role, created_at, updated_at
FROM tenant_users
WHERE tenant_id = $1
ORDER BY role ASC, created_at DESC;

-- name: ListUserTenants :many
SELECT tenant_id, role, created_at, updated_at
FROM tenant_users
WHERE user_id = $1;
