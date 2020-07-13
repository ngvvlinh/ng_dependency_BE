ALTER TABLE fulfillment
    ADD COLUMN coupon TEXT;
ALTER TABLE history.fulfillment
    ADD COLUMN coupon TEXT;

ALTER TABLE fulfillment
    ADD COLUMN insurance_value INT8;
ALTER TABLE history.fulfillment
    ADD COLUMN insurance_value INT8;

