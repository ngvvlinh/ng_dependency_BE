ALTER TABLE shop_supplier
    ADD COLUMN phone text,
    ADD COLUMN email text,
    ADD COLUMN company_name text,
    ADD COLUMN tax_number text,
    ADD COLUMN headquater_address text,
    ADD COLUMN full_name_norm tsvector,
    ADD COLUMN phone_norm tsvector;

CREATE INDEX ON shop_supplier USING gin(full_name_norm);

CREATE INDEX ON shop_supplier (phone);

CREATE INDEX ON shop_supplier (email);

CREATE INDEX ON shop_supplier USING gin(phone_norm);

-- Drop init history shop_supplier
DROP TRIGGER save_history ON shop_supplier;
DROP TRIGGER update_rid ON shop_supplier;
select init_history('shop_supplier', '{id,shop_id}');
