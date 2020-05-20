create type trader_type as enum('customer', 'vendor', 'carrier');

ALTER TABLE shop_trader
    ADD COLUMN type trader_type not null;

ALTER TABLE history.shop_trader
    ADD COLUMN type trader_type;

create type receipt_type as enum('receipt', 'payment');

CREATE TABLE receipt (
    id INT8 PRIMARY KEY,
    shop_id INT8 NOT NULL REFERENCES shop(id),
    trader_id INT8 NOT NULL REFERENCES shop_trader(id),
    user_id INT8 NOT NULL REFERENCES "user"(id),
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    amount INT4,
    status INT2,
    type receipt_type NOT NULL,
    lines JSONB,
    order_ids INT8[],
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX ON receipt USING GIN(order_ids);
CREATE UNIQUE INDEX ON receipt (shop_id, code);
