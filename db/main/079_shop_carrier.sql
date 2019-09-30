CREATE TABLE shop_carrier (
    id INT8 PRIMARY KEY REFERENCES shop_trader(id)
    , shop_id INT8 REFERENCES shop(id)
    , full_name TEXT
    , note TEXT
    , status INT2 NOT NULL
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);

select init_history('shop_carrier', '{id,shop_id}');