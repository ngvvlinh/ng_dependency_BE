create table shop_product_collection(
	product_id bigint not null,
	collection_id bigint not null,
	shop_id bigint,
	status int2,
	created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
 	FOREIGN KEY (product_id) REFERENCES shop_product(product_id),
 	FOREIGN KEY (shop_id) REFERENCES shop(id),
 	FOREIGN KEY (collection_id) REFERENCES shop_collection(id)
);

create unique index on shop_product_collection(product_id,collection_id) where deleted_at is null;
alter table shop_product_collection add constraint shop_product_collection_constraint primary key (product_id, collection_id);
