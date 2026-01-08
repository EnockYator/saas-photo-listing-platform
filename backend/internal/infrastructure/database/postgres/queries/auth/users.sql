-- name: CreateUser :one
INSERT INTO users (email, password_hash, created_at)
VALUES ($1, $2, NOW())
RETURNING id, email, created_at;

-- name: GetUserByID :one
SELECT id, email, password_hash, created_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, email, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserPasswordById :one
UPDATE users
SET password_hash = $2
WHERE id = $1
RETURNING id, email, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE email = $1;

-- name: CountUsers :one
SELECT COUNT(*) AS count
FROM users;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING id, email, created_at;

-- name: DeleteUserByEmail :exec
DELETE FROM users
WHERE email = $1; 

-- name: GetUserPasswordHash :one
SELECT password_hash
FROM users
WHERE id = $1;

-- name: UpdateUsername :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING id, email, username, created_at;

-- name: ListUsersByIDs :many
SELECT id, email, username, created_at
FROM users
WHERE id = ANY($1)
ORDER BY created_at DESC;    

-- name: UpdateUserPasswordByEmail :one
UPDATE users
SET password_hash = $2
WHERE email = $1
RETURNING id, email, created_at;  

-- name: GetUserCreationDate :one
SELECT created_at
FROM users
WHERE id = $1;  

-- name: ListRecentUsers :many
SELECT id, email, username, created_at
FROM users
WHERE created_at >= NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;

-- name: ListUsersByCreationDate :many
SELECT id, email, username, created_at
FROM users
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC; 
 
-- name: CountUsersByCreationDate :one
SELECT COUNT(*) AS count
FROM users
WHERE created_at BETWEEN $1 AND $2; 
