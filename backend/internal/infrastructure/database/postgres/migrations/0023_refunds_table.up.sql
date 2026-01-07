CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY,
    payment_id UUID NOT NULL REFERENCES payments(id),
    
    amount DECIMAL(10,2) NOT NULL
        CONSTRAINT refunds_amount_positive_check CHECK(amount >= 0),
    
    reason TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
