drop table shop_supplier_variant;

create table shop_variant_supplier (
	shop_id int8 REFERENCES shop(id),
	supplier_id int8 REFERENCES shop_supplier(id),
	variant_id int8 REFERENCES shop_variant(variant_id),
	created_at timestamptz,
	updated_at timestamptz
);

CREATE UNIQUE INDEX  ON "shop_variant_supplier" (shop_id, supplier_id, variant_id);
CREATE INDEX ON "shop_variant_supplier" (supplier_id);
CREATE INDEX ON "shop_variant_supplier" (variant_id);
