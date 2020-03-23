create table if not exists shop_brand (
	shop_id int8,
	id int8 primary key,
	description text,
	brand_name text,
	updated_at timestamptz,
	created_at timestamptz,
	deleted_at timestamptz,
	rid int8
);