-- name: InsertAuditLog :one
INSERT INTO audit_logs (
    tenant_id,
    performed_by,
    entity_id,
    entity_type,
    action,
    changed_data,
    performed_at
)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING id, tenant_id, performed_by, entity_id, entity_type, action, changed_data, performed_at;

-- name: ListAuditLogsByTenant :many
SELECT id, performed_by, entity_id, entity_type, action, changed_data, performed_at
FROM audit_logs
WHERE tenant_id = $1
ORDER BY performed_at DESC
LIMIT $2;

-- name: ListAuditLogsByEntity :many
SELECT id, performed_by, action, changed_data, performed_at
FROM audit_logs
WHERE tenant_id = $1
  AND entity_id = $2
  AND entity_type = $3
ORDER BY performed_at DESC;

-- name: CountAuditLogsByTenant :one
SELECT COUNT(*) AS count
FROM audit_logs
WHERE tenant_id = $1;
