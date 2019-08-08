ALTER TABLE money_transaction_shipping
    ADD COLUMN type TEXT;

ALTER TABLE history.money_transaction_shipping
    ADD COLUMN type TEXT;
