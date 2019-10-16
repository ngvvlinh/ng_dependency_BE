alter table shop_customer add column phone_norm tsvector;
alter table history.shop_customer add column phone_norm tsvector;

CREATE INDEX ON shop_customer USING gin(phone_norm);
