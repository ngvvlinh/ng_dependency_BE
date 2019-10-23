CREATE TYPE shop_ledger_type as enum('cash', 'bank');

CREATE TABLE shop_ledger (
    id INT8 PRIMARY KEY
    , shop_id INT8 REFERENCES shop(id)
    , name TEXT NOT NULL
    , bank_account JSONB
    , status INT2
    , note TEXT
    , type shop_ledger_type NOT NULL
    , created_by INT8 REFERENCES "user"(id)
    , created_at TIMESTAMP WITH TIME ZONE
    , updated_at TIMESTAMP WITH TIME ZONE
    , deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('shop_ledger', '{id, shop_id}');

ALTER TABLE receipt ADD COLUMN shop_ledger_id INT8 REFERENCES shop_ledger(id);
ALTER TABLE "history".receipt ADD COLUMN shop_ledger_id INT8;
