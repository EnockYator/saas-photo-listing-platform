-- name: CreatePlan :one
INSERT INTO plans (type, price, billing_cycle, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id, type, price, billing_cycle, created_at, updated_at;

-- name: GetPlanByID :one
SELECT id, type, price, billing_cycle, created_at, updated_at
FROM plans
WHERE id = $1;

-- name: ListPlans :many
SELECT id, type, price, billing_cycle, created_at, updated_at
FROM plans
ORDER BY created_at DESC;

-- name: UpdatePlan :one
UPDATE plans
SET price = $2, billing_cycle = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, type, price, billing_cycle, created_at, updated_at;

