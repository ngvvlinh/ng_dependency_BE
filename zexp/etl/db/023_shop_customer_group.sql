create table if not exists shop_customer_group (
	id bigint primary key,
	name text,
	partner_id bigint,
	shop_id bigint,
	created_at timestamp with time zone,
	updated_at timestamp with time zone,
    rid bigint
);