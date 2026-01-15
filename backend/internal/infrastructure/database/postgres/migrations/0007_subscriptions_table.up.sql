CREATE TABLE  IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,  
    
    status TEXT NOT NULL DEFAULT 'inactive'
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 85236ed (fix: fix conflicts in branches (#46))
        CONSTRAINT subscription_status_check
            CHECK (status IN (
                'active',
                'inactive',
                'canceled',
                'past_due'
            )),
<<<<<<< HEAD
   
=======
        CONSTRAINT subscription_status_check CHECK (status IN (
<<<<<<< HEAD
=======
            'active',
>>>>>>> e5eedb5 (chore(database): add indexes and constraints to database tables for fast performance and security)
            'inactive',
            'canceled',
            'past_due'
        )),
=======
>>>>>>> 85236ed (fix: fix conflicts in branches (#46))
    
>>>>>>> 8790134 (chore(database): add indexes and constraints to database tables for fast performance and security)
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

CREATE INDEX idx_subscription_tenant_id
    ON subscriptions(tenant_id);

CREATE INDEX idx_subscription_tenant_status
    ON subscriptions(tenant_id, status);

CREATE INDEX idx_subscription_created_at
    ON subscriptions(created_at);