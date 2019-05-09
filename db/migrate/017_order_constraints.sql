ALTER TABLE "order"
  ADD COLUMN order_discount INT4;
ALTER TABLE history."order"
  ADD COLUMN order_discount INT4;

ALTER TABLE order_line
  ALTER COLUMN code SET NOT NULL,
  ALTER COLUMN status SET NOT NULL,
  ALTER COLUMN weight SET NOT NULL,
  ALTER COLUMN is_outside_etop SET NOT NULL;

ALTER TABLE "order"
  ADD CONSTRAINT payment_price
    CHECK (COALESCE(order_discount,0) >= 0 AND COALESCE(order_discount,0) <= total_discount);

ALTER TABLE "order"
  ADD CONSTRAINT total_amount
    CHECK (total_amount = basket_value - total_discount + shop_shipping_fee);

ALTER TABLE order_line
  ADD COLUMN is_free BOOLEAN;
ALTER TABLE history.order_line
  ADD COLUMN is_free BOOLEAN;

ALTER TABLE order_line
  ADD CONSTRAINT quantity CHECK (quantity > 0);

ALTER TABLE order_line
  ADD CONSTRAINT payment_price
    CHECK (payment_price >= 0 AND payment_price <= retail_price AND (payment_price != 0) != COALESCE(is_free,FALSE));

ALTER TABLE order_line
  ADD CONSTRAINT total_discount
    CHECK (total_discount >= 0 AND total_discount = (retail_price - payment_price)*quantity);
