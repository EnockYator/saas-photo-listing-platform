-- name: CreateAuthSession :one
INSERT INTO auth_sessions (user_id, tenant_id, refresh_token, expires_at, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING id, user_id, tenant_id, refresh_token, expires_at, created_at;

-- name: GetAuthSessionByToken :one
SELECT id, user_id, tenant_id, refresh_token, expires_at, created_at
FROM auth_sessions
WHERE refresh_token = $1;

-- name: DeleteAuthSession :exec
DELETE FROM auth_sessions
WHERE user_id = $1
  AND tenant_id = $2
  AND refresh_token = $3;

-- name: DeleteAllUserSessions :exec
DELETE FROM auth_sessions
WHERE user_id = $1
  AND tenant_id = $2;

-- name: ListUserSessions :many
SELECT id, refresh_token, expires_at, created_at
FROM auth_sessions
WHERE user_id = $1
  AND tenant_id = $2
ORDER BY created_at DESC;

-- name: DeleteExpiredSessions :exec
DELETE FROM auth_sessions
WHERE expires_at < NOW();
