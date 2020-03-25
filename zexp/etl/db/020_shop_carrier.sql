CREATE TABLE if not exists shop_carrier (
    id INT8 PRIMARY KEY
    , shop_id INT8
    , full_name TEXT
    , note TEXT
    , status INT2
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , rid INT8
);