create table shop_supplier_variant (
	shop_id int8,
	supplier_id int8,
	variant_id int8,
	created_at timestamptz,
	updated_at timestamptz
);

CREATE UNIQUE INDEX  ON "shop_supplier_variant" (shop_id, supplier_id, variant_id);
CREATE INDEX ON "shop_supplier_variant" (supplier_id);
CREATE INDEX ON "shop_supplier_variant" (variant_id);
