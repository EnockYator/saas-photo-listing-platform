CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    type TEXT NOT NULL
        CONSTRAINT plan_type_check CHECK (type IN ('free', 'basic', 'business')),
    
    price DECIMAL(10,2) NOT NULL
        CONSTRAINT plan_price_positive_check CHECK (price >= 0),

    billing_cycle TEXT NOT NULL DEFAULT 'monthly'
 
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 975dd66 (chore(github): rebase main to dev (#45))
    ON plans(billing_cycle);
=======
    ON plans(billing_cycle);
>>>>>>> 9e4fead (chore(github): rebase main to dev (#45))
<<<<<<< HEAD
=======
    ON plans(billing_cycle);
>>>>>>> e215485 (refactor: place all miration files in migrations/ directory rather than sub-directories)
=======
>>>>>>> 975dd66 (chore(github): rebase main to dev (#45))
