CREATE TABLE if not exists shop_collection (
    id bigint PRIMARY KEY,
    shop_id bigint,
    name text,
    description text,
    desc_html text,
    short_desc text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    partner_id INT8,
    external_id TEXT,
    rid bigint
);