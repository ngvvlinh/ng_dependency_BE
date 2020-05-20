CREATE TABLE transaction (
    id BIGINT PRIMARY KEY,
    amount INT,
    account_id BIGINT,
    status INT2,
    type TEXT,
    note TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

select init_history('transaction', '{id}');
