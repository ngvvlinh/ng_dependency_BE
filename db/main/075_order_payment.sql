CREATE TABLE payment (
    id BIGINT PRIMARY KEY,
    amount INT,
    status INT2,
    state TEXT,
    payment_provider TEXT,
    external_trans_id TEXT,
    external_data JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON payment (external_trans_id, payment_provider);

ALTER TABLE "order"
    ADD COLUMN payment_status INT2,
    ADD COLUMN payment_id BIGINT REFERENCES payment(id);

ALTER TABLE history."order"
    ADD COLUMN payment_status INT2,
    ADD COLUMN payment_id BIGINT;

select init_history('payment', '{id}');
