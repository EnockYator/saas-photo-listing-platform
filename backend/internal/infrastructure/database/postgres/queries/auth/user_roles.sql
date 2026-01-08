-- name: AssignRoleToUser :one
INSERT INTO user_roles (user_id, role_id, tenant_id)
VALUES ($1, $2, $3)
RETURNING id, user_id, role_id, tenant_id, assigned_at;

-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1
  AND role_id = $2
  AND tenant_id = $3;

-- name: ListUserRoles :many
SELECT ur.role_id, r.name, r.description
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1
  AND ur.tenant_id = $2;

-- name: ListUsersByRole :many
SELECT u.id, u.username, u.email
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
WHERE ur.role_id = $1
  AND ur.tenant_id = $2;
