-- name: CreateSubscription :one
INSERT INTO subscriptions (tenant_id, plan_id, status, started_at, end_at, created_at)
VALUES ($1, $2, 'inactive', $3, $4, NOW())
RETURNING id, tenant_id, plan_id, status, started_at, end_at, created_at, updated_at;

-- name: GetSubscriptionByID :one
SELECT *
FROM subscriptions
WHERE tenant_id = $1
  AND id = $2;

-- name: ListSubscriptionsByTenant :many
SELECT *
FROM subscriptions
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: UpdateSubscriptionStatus :one
UPDATE subscriptions
SET status = $3, updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
RETURNING id, tenant_id, plan_id, status, started_at, end_at, created_at, updated_at;

-- name: UpdateSubscriptionPlan :one
UPDATE subscriptions
SET plan_id = $3, updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
RETURNING id, tenant_id, plan_id, status, started_at, end_at, created_at, updated_at;
