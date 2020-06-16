ALTER TABLE shipment_price
    ADD COLUMN additional_fees JSONB;

ALTER TABLE history.shipment_price
    ADD COLUMN additional_fees JSONB;

ALTER TABLE shipment_price_list
    RENAME COLUMN is_active TO is_default;

ALTER TABLE history.shipment_price_list
    RENAME COLUMN is_active TO is_default;
