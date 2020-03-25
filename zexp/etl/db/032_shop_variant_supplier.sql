create table if not exists shop_variant_supplier (
	shop_id int8,
	supplier_id int8,
	variant_id int8,
	created_at timestamptz,
	updated_at timestamptz,
	rid int8,
	primary key (supplier_id, variant_id)
);
