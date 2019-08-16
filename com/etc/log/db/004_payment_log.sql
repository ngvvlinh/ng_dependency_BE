CREATE TABLE payment (
    id BIGINT PRIMARY KEY,
    payment_provider TEXT,
    data JSONB,
    order_id TEXT,
    action TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
