-- name: CreateInvoice :one
INSERT INTO invoices (tenant_id, subscription_id, amount, currency, status, issued_at)
VALUES ($1, $2, $3, $4, 'pending', NOW())
RETURNING id, tenant_id, subscription_id, amount, currency, status, issued_at, paid_at, updated_at;

-- name: GetInvoiceByID :one
SELECT id, tenant_id, subscription_id, amount, currency, status, issued_at, paid_at, updated_at
FROM invoices
WHERE tenant_id = $1
  AND id = $2;

-- name: ListInvoicesByTenant :many
SELECT id, subscription_id, amount, currency, status, issued_at, paid_at
FROM invoices
WHERE tenant_id = $1
ORDER BY issued_at DESC
LIMIT $2;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET status = $3, paid_at = $4
WHERE tenant_id = $1
  AND id = $2
RETURNING id, tenant_id, subscription_id, amount, currency, status, issued_at, paid_at, updated_at;

