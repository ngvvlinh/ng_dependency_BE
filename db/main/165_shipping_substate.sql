ALTER TABLE fulfillment
    ADD COLUMN shipping_substate TEXT;

ALTER TABLE history.fulfillment
    ADD COLUMN shipping_substate TEXT;
