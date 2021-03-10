ALTER TABLE fulfillment
    ADD COLUMN external_sort_code TEXT;

ALTER TABLE history.fulfillment
    ADD COLUMN external_sort_code TEXT;
