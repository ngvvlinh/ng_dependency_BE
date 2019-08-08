ALTER TABLE order_line DROP CONSTRAINT payment_price;
ALTER TABLE order_line ADD CONSTRAINT payment_price
  CHECK (payment_price >= 0 AND payment_price <= retail_price);
