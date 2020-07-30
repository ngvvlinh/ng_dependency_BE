ALTER TABLE shipping_provider_webhook
    ADD COLUMN connection_id INT8;

CREATE INDEX ON shipping_provider_webhook(connection_id);
