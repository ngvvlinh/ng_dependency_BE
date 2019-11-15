ALTER TABLE inventory_voucher
    ADD COLUMN code_norm int4,
    ADD COLUMN code text;

CREATE UNIQUE INDEX ON inventory_voucher(shop_id, code);
