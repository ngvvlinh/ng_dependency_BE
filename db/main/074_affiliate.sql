ALTER TABLE "user"
    ADD COLUMN ref_user_id BIGINT REFERENCES "user"(id),
    ADD COLUMN ref_sale_id BIGINT REFERENCES "user"(id);

ALTER TABLE history."user"
    ADD COLUMN ref_user_id BIGINT,
    ADD COLUMN ref_sale_id BIGINT;

-- Create shop **eTop Trading** belongs to user **eTop System**
INSERT INTO "account" (
    "id", "name", "type", "owner_id"
) VALUES (
    '1000015765615455091', 'eTop Trading', 'shop', '1000101010101010101'
);

INSERT INTO shop (
    "id", "name", "owner_id", "status", "created_at", "updated_at", "is_test", "try_on"
) VALUES (
    '1000015765615455091', 'eTop Trading', '1000101010101010101', 1, NOW(), NOW(), 1, 'open'
);

-- Add field **trading_shop_id** to table shop
ALTER TABLE "order" ADD COLUMN trading_shop_id BIGINT REFERENCES shop(id);

ALTER TABLE history."order" ADD COLUMN trading_shop_id BIGINT REFERENCES shop(id);

-- Add account type: **affiliate**
ALTER TYPE account_type ADD VALUE 'affiliate';

CREATE TABLE affiliate (
    id BIGINT PRIMARY KEY REFERENCES account(id),
    rid BIGINT,
    name TEXT not null,
    owner_id BIGINT NOT NULL REFERENCES "user"(id),
    status INT2,
    is_test SMALLINT DEFAULT '0'::SMALLINT,
    phone TEXT,
    email TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
