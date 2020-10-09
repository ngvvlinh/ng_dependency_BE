ALTER TABLE shipnow_fulfillment ADD COLUMN address_to_phone tsvector;
ALTER TABLE history.shipnow_fulfillment ADD COLUMN address_to_phone tsvector;
CREATE INDEX ON shipnow_fulfillment USING gin(address_to_phone);

ALTER TABLE shipnow_fulfillment ADD COLUMN address_to_full_name_norm tsvector;
ALTER TABLE history.shipnow_fulfillment ADD COLUMN address_to_full_name_norm tsvector;
CREATE INDEX ON shipnow_fulfillment USING gin(address_to_full_name_norm);
