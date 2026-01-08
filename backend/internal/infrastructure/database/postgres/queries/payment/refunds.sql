-- name: CreateRefund :one
INSERT INTO refunds (id, payment_id, tenant_id, amount, reason, created_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())
RETURNING id, payment_id, tenant_id, amount, reason, created_at;

-- name: ListRefundsByPayment :many
SELECT id, amount, reason, created_at
FROM refunds
WHERE tenant_id = $1
  AND payment_id = $2
ORDER BY created_at DESC;

-- name: TotalRefundedByInvoice :one
SELECT SUM(r.amount) AS total_refunded
FROM refunds r
JOIN payments p ON r.payment_id = p.id
WHERE r.tenant_id = $1
  AND p.invoice_id = $2;
