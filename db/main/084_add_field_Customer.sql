alter table shop_customer add column full_name_norm tsvector;

CREATE INDEX ON shop_customer USING gin(full_name_norm);
CREATE INDEX ON shop_customer (phone);
