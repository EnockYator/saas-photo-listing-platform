-- name: ListRoles :many
SELECT id, name, created_at
FROM roles
ORDER BY name;

-- name: GetRoleByName :one
SELECT id, name, created_at
FROM roles
WHERE name = $1;

-- name: AssignRoleToUser :one
INSERT INTO user_roles (user_id, tenant_id, role_id)
VALUES ($1, $2, $3)
RETURNING id, user_id, tenant_id, role_id, assigned_at;

-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1
  AND tenant_id = $2
  AND role_id = $3;

-- name: ListUserRoles :many
SELECT ur.role_id, r.name
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1
  AND ur.tenant_id = $2;
