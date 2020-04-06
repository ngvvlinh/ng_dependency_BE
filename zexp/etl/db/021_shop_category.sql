CREATE TABLE if not exists shop_category (
    id bigint PRIMARY KEY,
    shop_id int8,
    supplier_id int8,
    parent_id int8,
    name text,
    status int2,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    rid int8
);

alter table shop_category
    drop column if exists supplier_id;