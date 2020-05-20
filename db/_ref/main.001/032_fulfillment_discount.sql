-- total_amount and total_discount is equal to order when supplier_id is null
ALTER TABLE fulfillment
  ADD COLUMN total_discount INT4,
  ADD COLUMN total_amount INT4;

ALTER TABLE history.fulfillment
  ADD COLUMN total_discount INT4,
  ADD COLUMN total_amount INT4;

UPDATE fulfillment
  SET total_amount = "order".total_amount,
      total_discount = "order".total_discount
  FROM "order"
  WHERE fulfillment.order_id = "order".id;

ALTER TABLE fulfillment
  ALTER COLUMN total_discount SET NOT NULL,
  ALTER COLUMN total_amount SET NOT NULL;

UPDATE money_transaction_shipping SET etop_transfered_at = updated_at
  WHERE etop_transfered_at IS NULL AND status = 1;
