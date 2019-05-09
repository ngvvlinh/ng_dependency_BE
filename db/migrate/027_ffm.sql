ALTER TABLE fulfillment ADD COLUMN IF NOT EXISTS expected_pick_at TIMESTAMPTZ;
ALTER TABLE history.fulfillment ADD COLUMN IF NOT EXISTS expected_pick_at TIMESTAMPTZ;
