ALTER TABLE shop_customer
    ADD COLUMN code_norm int4;
ALTER TABLE history."shop_customer"
    ADD COLUMN code_norm int4;

CREATE UNIQUE INDEX ON shop_customer(shop_id, code);