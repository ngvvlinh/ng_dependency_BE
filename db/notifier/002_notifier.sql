ALTER TABLE notification ADD COLUMN
    synced_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE notification ADD COLUMN
    send_notification BOOLEAN;
