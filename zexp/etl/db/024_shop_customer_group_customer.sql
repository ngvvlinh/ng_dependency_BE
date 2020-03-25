create table if not exists shop_customer_group_customer (
	customer_id bigint,
	group_id bigint,

	created_at timestamp with time zone,
	updated_at timestamp with time zone,
	rid bigint,
	primary key (customer_id, group_id)
);