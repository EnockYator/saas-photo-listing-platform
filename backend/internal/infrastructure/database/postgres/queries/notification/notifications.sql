-- name: CreateNotification :one
INSERT INTO notifications (tenant_id, user_id, message, type, data, is_read, created_at)
VALUES ($1, $2, $3, $4, $5, FALSE, NOW())
RETURNING id, tenant_id, user_id, message, type, data, is_read, created_at;

-- name: GetUnreadNotifications :many
SELECT id, message, type, data, created_at
FROM notifications
WHERE tenant_id = $1
  AND user_id = $2
  AND is_read = FALSE
ORDER BY created_at DESC
LIMIT $3;

-- name: MarkNotificationRead :one
UPDATE notifications
SET is_read = TRUE, updated_at = NOW()
WHERE tenant_id = $1
  AND user_id = $2
  AND id = $3
RETURNING id, is_read, updated_at;

-- name: ListNotifications :many
SELECT id, message, type, data, is_read, created_at, updated_at
FROM notifications
WHERE tenant_id = $1
  AND user_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUnreadNotifications :one
SELECT COUNT(*) AS count
FROM notifications
WHERE tenant_id = $1
  AND user_id = $2
  AND is_read = FALSE;

-- name: ListNotificationsByType :many
SELECT id, message, is_read, created_at
FROM notifications
WHERE tenant_id = $1
  AND user_id = $2
  AND type = $3
ORDER BY created_at DESC;

