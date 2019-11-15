ALTER TABLE receipt
    ADD COLUMN trader_full_name_norm tsvector;

alter table history.receipt add column trader_full_name_norm tsvector;

CREATE INDEX ON receipt USING gin(trader_full_name_norm);
