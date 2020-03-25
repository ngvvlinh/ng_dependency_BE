create table if not exists shop_product_collection(
	shop_id bigint,
	product_id bigint,
	collection_id bigint,
	created_at timestamp with time zone,
    updated_at timestamp with time zone,
    rid bigint,
    primary key (product_id, collection_id)
);