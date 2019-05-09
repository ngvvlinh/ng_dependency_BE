ALTER TABLE "order" ADD COLUMN total_fee INT4;
ALTER TABLE "order" ADD COLUMN fee_lines JSONB;

ALTER TABLE history."order" ADD COLUMN total_fee INT4;
ALTER TABLE history."order" ADD COLUMN fee_lines JSONB;

ALTER TABLE "fulfillment" ADD COLUMN external_shipping_name TEXT;
ALTER TABLE "fulfillment" ADD COLUMN shipping_service_fee INT4;
ALTER TABLE "fulfillment" ADD COLUMN original_cod_amount INT4;
ALTER TABLE "fulfillment" ADD COLUMN address_return JSONB;
ALTER TABLE "fulfillment" ADD COLUMN include_insurance BOOLEAN;

ALTER TABLE history."fulfillment" ADD COLUMN external_shipping_name TEXT;
ALTER TABLE history."fulfillment" ADD COLUMN shipping_service_fee INT4;
ALTER TABLE history."fulfillment" ADD COLUMN original_cod_amount INT4;
ALTER TABLE history."fulfillment" ADD COLUMN address_return JSONB;
ALTER TABLE history."fulfillment" ADD COLUMN include_insurance BOOLEAN;

ALTER TABLE "order" DROP CONSTRAINT total_amount;
ALTER TABLE "order" ADD CONSTRAINT total_amount
   CHECK (total_amount = basket_value - total_discount +
      COALESCE(total_fee, shop_shipping_fee));
