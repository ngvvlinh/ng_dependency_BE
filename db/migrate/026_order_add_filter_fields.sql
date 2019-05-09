ALTER TABLE "order"
  ADD COLUMN customer_name TEXT,
  ADD COLUMN customer_name_norm tsvector,
  ADD COLUMN product_name_norm tsvector;

ALTER TABLE history."order"
  ADD COLUMN customer_name TEXT,
  ADD COLUMN customer_name_norm tsvector,
  ADD COLUMN product_name_norm tsvector;

ALTER TABLE fulfillment
  ADD COLUMN address_to_province_code TEXT,
  ADD COLUMN address_to_district_code TEXT,
  ADD COLUMN address_to_ward_code TEXT;

ALTER TABLE history.fulfillment
  ADD COLUMN address_to_province_code TEXT,
  ADD COLUMN address_to_district_code TEXT,
  ADD COLUMN address_to_ward_code TEXT;

CREATE INDEX ON "order" USING gin(customer_name_norm);
CREATE INDEX ON "order" USING gin(product_name_norm);
CREATE INDEX ON "order" USING gin(fulfillment_shipping_states);
CREATE INDEX ON "order" USING gin(fulfillment_shipping_codes);
CREATE INDEX ON "order" (shop_id);
CREATE INDEX ON "order" (customer_phone);
CREATE INDEX ON "order" (status);
CREATE INDEX ON "order" (code);
CREATE INDEX ON "order" (total_amount);
CREATE INDEX ON "order" (confirm_status);
CREATE INDEX ON "order" (fulfillment_shipping_status);
CREATE INDEX ON "order" (etop_payment_status);

CREATE INDEX ON fulfillment (shipping_fee_shop);
CREATE INDEX ON fulfillment (total_cod_amount);
CREATE INDEX ON fulfillment (shipping_code);
CREATE INDEX ON fulfillment (money_transaction_id);
CREATE INDEX ON fulfillment (shipping_state);
CREATE INDEX ON fulfillment (address_to_district_code);
CREATE INDEX ON fulfillment (address_to_province_code);
CREATE INDEX ON fulfillment (address_to_ward_code);
