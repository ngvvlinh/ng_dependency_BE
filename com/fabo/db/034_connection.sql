ALTER TABLE shop_connection
    ADD COLUMN last_sync_at timestamptz,
    ADD COLUMN telecom_data JSONB;
ALTER TABLE history.shop_connection
    ADD COLUMN last_sync_at timestamptz,
    ADD COLUMN telecom_data JSONB;
