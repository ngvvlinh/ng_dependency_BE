ALTER TYPE shipping_provider ADD VALUE 'vtpost';
ALTER TABLE fulfillment ADD COLUMN IF NOT EXISTS external_shipping_state_code TEXT;
ALTER TABLE history.fulfillment ADD COLUMN IF NOT EXISTS external_shipping_state_code TEXT;

ALTER TABLE money_transaction_shipping_external_line ADD COLUMN IF NOT EXISTS external_total_shipping_fee INTEGER;

ALTER TABLE shipping_source_state RENAME TO shipping_source_internal;
ALTER TABLE shipping_source_internal 
	ADD COLUMN IF NOT EXISTS secret jsonb,
	ADD COLUMN IF NOT EXISTS access_token text,
	ADD COLUMN IF NOT EXISTS expires_at timestamp with time zone;
