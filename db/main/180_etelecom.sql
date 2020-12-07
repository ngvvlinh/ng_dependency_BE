ALTER TABLE shop_connection
    ADD COLUMN last_sync_at timestamptz,
    ADD COLUMN telecom_data JSONB;
ALTER TABLE history.shop_connection
    ADD COLUMN last_sync_at timestamptz,
    ADD COLUMN telecom_data JSONB;

INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_method, connection_provider, code)
VALUES (100085369475949390, 'VHT - Telecom', 1, now(), now(), 'crm/_/builtin/vht', 'telecom', 'builtin', 'vht', 'LN20');

INSERT INTO shop_connection(connection_id, token, status, is_global, created_at, updated_at, telecom_data, token_expires_at, last_sync_at)
VALUES (100085369475949390, 'default_token', 0, true, now(), now(), '{"username": "", "password": "", "tenant_host": "", "tenant_token": ""}', now(), now());
