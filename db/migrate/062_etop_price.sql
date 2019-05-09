-- etop_adjusted_shipping_fee_main: store etop price after apply price rule
-- etop_price_rule: true|false (apply eTop price or not)
ALTER TABLE fulfillment 
  ADD COLUMN etop_adjusted_shipping_fee_main INT,
  ADD COLUMN etop_price_rule BOOLEAN;
ALTER TABLE history.fulfillment 
	ADD COLUMN etop_adjusted_shipping_fee_main INT,
  	ADD COLUMN etop_price_rule BOOLEAN;
