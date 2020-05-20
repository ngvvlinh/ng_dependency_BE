ALTER TABLE fulfillment 
	ADD COLUMN IF NOT EXISTS shipping_created_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_picking_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_holding_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_delivering_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_returning_at timestamp with time zone;

ALTER TABLE fulfillment
  RENAME delivered_at TO shipping_delivered_at;
ALTER TABLE fulfillment
  RENAME returned_at TO shipping_returned_at;
ALTER TABLE fulfillment
  RENAME cancelled_at TO shipping_cancelled_at;

ALTER TABLE history.fulfillment 
	ADD COLUMN IF NOT EXISTS shipping_created_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_holding_at timestamp with time zone,
  ADD COLUMN IF NOT EXISTS shipping_picking_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_delivering_at timestamp with time zone,
	ADD COLUMN IF NOT EXISTS shipping_returning_at timestamp with time zone;

ALTER TABLE history.fulfillment
  RENAME delivered_at TO shipping_delivered_at;
ALTER TABLE history.fulfillment
  RENAME returned_at TO shipping_returned_at;
ALTER TABLE history.fulfillment
  RENAME cancelled_at TO shipping_cancelled_at;
