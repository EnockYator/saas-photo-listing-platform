CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),

    amount DECIMAL(10,2) NOT NULL
    CONSTRAINT invoice_amount_positive_check CHECK(amount >= 0),

    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    
    status TEXT NOT NULL DEFAULT 'pending'
        CONSTRAINT invoice_status_check CHECK (status IN ('pending', 'paid', 'failed')),

    issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    paid_at TIMESTAMPTZ,


    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger for updated_at
CREATE TRIGGER trg_invoice_updated_at
BEFORE UPDATE ON invoices
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE INDEX idx_invoice_tenant_id
    ON invoices(tenant_id);

CREATE INDEX idx_invoice_subscription_id
    ON invoices(subscription_id);

CREATE INDEX idx_invoice_status
    ON invoices(status);


