CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    type TEXT NOT NULL
        CONSTRAINT plan_type_check CHECK (type IN ('free', 'basic', 'business')),
    
    price DECIMAL(10,2) NOT NULL
        CONSTRAINT plan_price_positive_check CHECK (price >= 0),

    billing_cycle TEXT NOT NULL DEFAULT 'monthly'
        CONSTRAINT plan_billing_cycle_check CHECK (billing_cycle IN ('monthly', 'yearly')),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_plans_updated_at
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Index for querying by type
CREATE INDEX idx_plan_type
    ON plans(type);

-- Index for querying by billing_cycle
CREATE INDEX idx_plans_billing_cycle
    ON plans(billing_cycle);
