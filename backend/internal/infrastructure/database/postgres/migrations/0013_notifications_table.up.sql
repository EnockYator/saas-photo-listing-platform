
--  Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    message TEXT NOT NULL,
    
    type TEXT NOT NULL
        CONSTRAINT notifications_type_check CHECK (type IN ('info', 'warning', 'alert')),
    
    data JSONB,            -- dynamic payload

    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_notification_timestamps
        CHECK (updated_at >= created_at)
);

-- Trigger bound to the table
CREATE TRIGGER trg_notifications_updated_at
BEFORE UPDATE ON notifications
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Fetch unread notifications fast
CREATE INDEX idx_notifications_user_unread
    ON notifications(user_id, created_at DESC)
    WHERE is_read = FALSE;

-- Admin view, all notifications per tenant & user 
CREATE INDEX idx_notifications_tenant_user
    ON notifications(tenant_id, user_id, created_at DESC);

CREATE INDEX idx_notifications_type
    ON notifications(tenant_id, type);
