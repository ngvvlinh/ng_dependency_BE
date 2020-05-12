create table fb_external_user_shop_customer (
    created_at timestamptz,
    updated_at timestamptz,
    shop_id int8,
    fb_external_user_id text,
    customer_id int8,
    status INT2
);

CREATE UNIQUE INDEX ON fb_external_user_shop_customer (shop_id, customer_id, fb_external_user_id);
