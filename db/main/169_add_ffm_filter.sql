ALTER TABLE fulfillment ADD COLUMN address_to_phone text;
ALTER TABLE history.fulfillment ADD COLUMN address_to_phone text;
CREATE INDEX ON fulfillment USING btree(address_to_phone);

ALTER TABLE fulfillment ADD COLUMN address_to_full_name_norm tsvector;
ALTER TABLE history.fulfillment ADD COLUMN address_to_full_name_norm tsvector;
CREATE INDEX ON fulfillment USING gin(address_to_full_name_norm);
