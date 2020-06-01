ALTER TABLE fb_customer_conversation
    ADD COLUMN last_message_external_id TEXT;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN last_message_external_id TEXT;