ALTER TABLE device DROP CONSTRAINT device_device_id_key;
ALTER TABLE device DROP CONSTRAINT device_external_device_id_key;
ALTER TABLE device
    ADD COLUMN user_id INT8,
    ADD COLUMN config JSONB,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

CREATE UNIQUE INDEX ON device (external_device_id, external_service_id, user_id);
