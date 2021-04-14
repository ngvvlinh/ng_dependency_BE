ALTER TABLE payment
    ADD COLUMN shop_id INT8 REFERENCES shop(id);
ALTER TABLE history.payment
    ADD COLUMN shop_id INT8;
