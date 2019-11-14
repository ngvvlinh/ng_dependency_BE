ALTER TABLE shop_supplier
    ADD COLUMN code_norm int4,
    ADD COLUMN code text;

ALTER TABLE history."shop_supplier"
    ADD COLUMN code_norm int4,
    ADD COLUMN code text;

CREATE UNIQUE INDEX ON shop_supplier(shop_id, code);
