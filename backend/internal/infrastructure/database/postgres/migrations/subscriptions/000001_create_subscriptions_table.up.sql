CREATE TYPE subscription_status AS ENUM (
    'active',
    'inactive',
    'canceled',
    'past_due'
);

CREATE TABLE  IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,  
    
    status subscription_status NOT NULL DEFAULT 'inactive',
    
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    end_at TIMESTAMPTZ NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3️⃣ Trigger bound to the table
CREATE TRIGGER trg_subscriptions_updated_at
BEFORE UPDATE ON subscriptions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_subscriptions_tenant_id
    ON subscriptions(tenant_id);

CREATE INDEX idx_subscriptions_tenant_status
    ON subscriptions(tenant_id, status);

CREATE INDEX idx_subscriptions_created_at
    ON subscriptions(created_at);

