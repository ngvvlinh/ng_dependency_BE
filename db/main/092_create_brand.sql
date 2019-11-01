create table shop_brand (
	shop_id int8,
	id int8,
	description text,
	brand_name text,
	updated_at timestamptz,
	created_at timestamptz,
	deleted_at timestamptz
	);

CREATE UNIQUE INDEX ON "shop_brand" (shop_id, id);

alter table shop_product add column brand_id int8;
alter table history.shop_product add column brand_id int8;
