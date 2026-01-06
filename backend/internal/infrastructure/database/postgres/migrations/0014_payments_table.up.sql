
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    
    amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0), --- e.g., 49.99
    currency VARCHAR(3) NOT NULL,
    
    status TEXT NOT NULL DEFAULT 'pending'
        CONSTRAINT payments_status_check CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),


    method TEXT NOT NULL
        CONSTRAINT payments_method_check CHECK (method IN (
            'credit_card',
            'paypal',
            'bank_transfer',
            'mpesa',
            'stripe'
        )),

    provider TEXT NOT NULL
        CONSTRAINT payments_provider_check CHECK (provider IN (
            'stripe',
            'paypal',
            'mpesa',
            'bank'
        )),

    provider_payment_id TEXT NOT NULL,  -- gateway reference ID
    idempotency_key TEXT NOT NULL,  -- to prevent duplicate payments
   
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (provider, provider_payment_id),
    UNIQUE (idempotency_key)
);

-- 3️⃣ Trigger bound to the table
CREATE TRIGGER trg_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_payments_tenant_subscription
    ON payments(tenant_id, subscription_id);

CREATE INDEX idx_payments_tenant_user
    ON payments(tenant_id, user_id);

CREATE INDEX idx_payments_tenant_status
    ON payments(tenant_id, status);

