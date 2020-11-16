ALTER TABLE shop
    ADD COLUMN is_prior_money_transaction BOOLEAN;

ALTER TABLE history.shop
    ADD COLUMN is_prior_money_transaction BOOLEAN;
