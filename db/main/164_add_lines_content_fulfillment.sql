
ALTER TABLE fulfillment
    ADD COLUMN lines_content TEXT;
ALTER TABLE history.fulfillment
    ADD COLUMN lines_content TEXT;

ALTER TABLE fulfillment
    ALTER COLUMN order_id drop not null;
ALTER TABLE fulfillment
    ALTER COLUMN lines drop not null;
ALTER TABLE fulfillment
    ALTER COLUMN variant_ids drop not null;

ALTER TABLE fulfillment
    ADD COLUMN ed_code TEXT;
ALTER TABLE history.fulfillment
    ADD COLUMN ed_code TEXT;

CREATE UNIQUE INDEX ON fulfillment(ed_code, shop_id) WHERE ed_code <> '' and status <> -1;
