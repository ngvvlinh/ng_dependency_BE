CREATE SEQUENCE shipping_code START 100001;

ALTER TABLE fulfillment DROP COLUMN code;
ALTER TABLE history.fulfillment DROP COLUMN code;
