ALTER TABLE "order"
    ADD COLUMN customer_id INT8 REFERENCES shop_customer(id);
ALTER TABLE history."order"
    ADD COLUMN customer_id INT8 REFERENCES shop_customer(id);

ALTER TABLE "shop_trader_address"
    ADD COLUMN is_default BOOLEAN;
ALTER TABLE history."shop_trader_address"
    ADD COLUMN is_default BOOLEAN;
