ALTER TABLE fulfillment
  ADD COLUMN admin_note TEXT,
  ADD COLUMN is_partial_delivery BOOLEAN;
ALTER TABLE history.fulfillment
  ADD COLUMN admin_note TEXT,
  ADD COLUMN is_partial_delivery BOOLEAN;
