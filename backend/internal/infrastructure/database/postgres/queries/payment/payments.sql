-- name: CreatePayment :one
INSERT INTO payments (user_id, tenant_id, invoice_id, subscription_id, amount, currency, status, method, provider, provider_payment_id, idempotency_key)
VALUES ($1, $2, $3, $4, $5, $6, 'pending', $7, $8, $9, $10)
RETURNING id, user_id, tenant_id, invoice_id, subscription_id, amount, currency, status, method, provider, provider_payment_id, idempotency_key, paid_at;

-- name: GetPaymentByID :one
SELECT *
FROM payments
WHERE tenant_id = $1
  AND id = $2;

-- name: ListPaymentsByTenant :many
SELECT *
FROM payments
WHERE tenant_id = $1
ORDER BY id DESC
LIMIT $2;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $3, paid_at = $4
WHERE tenant_id = $1
  AND id = $2
RETURNING *;
