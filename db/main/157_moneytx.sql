ALTER TABLE money_transaction_shipping_external
    ADD COLUMN connection_id INT8 REFERENCES connection(id)
    , ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE money_transaction_shipping
    ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE history.money_transaction_shipping
	ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE money_transaction_shipping_etop
    ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);
