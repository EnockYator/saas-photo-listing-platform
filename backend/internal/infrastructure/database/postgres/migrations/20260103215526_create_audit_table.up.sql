CREATE TABLE IF NOT EXISTS audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    
    action TEXT NOT NULL,
    changed_data JSONB,
    
    performed_by UUID,
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_entity_type ON audit(entity_type);
CREATE INDEX idx_audit_entity_id ON audit(entity_id);
CREATE INDEX idx_audit_performed_by ON audit(performed_by);