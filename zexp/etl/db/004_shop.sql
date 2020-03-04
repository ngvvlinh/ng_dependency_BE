CREATE TABLE if not exists shop (
    id bigint,
    rid bigint,
    name text,
    owner_id bigint,
    status int2,
    product_source_id int8,
    order_source_id int8,
    created_at timestamp with time zone DEFAULT date_trunc('second', now()),
    updated_at timestamp with time zone,
    rules jsonb,
    is_test smallint DEFAULT '0'::smallint,
    image_url text,
    phone text,
    website_url text,
    email text,
    deleted_at timestamp with time zone,
    address_id bigint,
    bank_account jsonb,
    contact_persons jsonb,
    inventory_overstock boolean,
    code text,
    auto_create_ffm boolean DEFAULT FALSE,
    recognized_hosts text[]
);