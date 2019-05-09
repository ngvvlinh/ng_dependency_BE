ALTER TYPE shipping_provider ADD VALUE 'ghtk';
ALTER TABLE fulfillment ADD COLUMN IF NOT EXISTS provider_service_id TEXT;
ALTER TABLE history.fulfillment ADD COLUMN IF NOT EXISTS provider_service_id TEXT;
