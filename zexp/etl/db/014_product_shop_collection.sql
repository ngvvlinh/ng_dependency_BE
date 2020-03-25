CREATE TABLE if not exists product_shop_collection (
    product_id bigint,
    shop_id bigint,
    collection_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    status smallint,
    rid bigint,
    primary key (product_id, collection_id)
);