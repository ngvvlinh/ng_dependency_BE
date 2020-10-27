ALTER TABLE shipnow_fulfillment
    ADD COLUMN driver_name TEXT
    , ADD COLUMN  driver_phone TEXT;

ALTER TABLE history.shipnow_fulfillment
    ADD COLUMN driver_name TEXT
    , ADD COLUMN  driver_phone TEXT;
