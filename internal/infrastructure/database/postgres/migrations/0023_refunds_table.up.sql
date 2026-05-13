CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY,
    payment_id UUID NOT NULL REFERENCES payments(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    amount DECIMAL(10,2) NOT NULL
        CONSTRAINT refunds_amount_positive_check CHECK(amount >= 0),
    
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_refunds_payment
    ON refunds(payment_id);

CREATE INDEX idx_refunds_tenant
    ON refunds(tenant_id);

