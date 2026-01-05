CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,  -- e.g., 'free', 'pro', 'business'
    
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger to keep updated_at fresh
CREATE TRIGGER trg_plans_updated_at
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
