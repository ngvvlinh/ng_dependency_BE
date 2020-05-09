ALTER TABLE fb_customer_conversation
    ADD COLUMN external_message_attachments JSONB;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN external_message_attachments JSONB;

ALTER TABLE fb_external_message
    ADD COLUMN external_sticker TEXT;
ALTER TABLE history.fb_external_message
    ADD COLUMN external_sticker TEXT;

ALTER TABLE fb_customer_conversation
    ADD COLUMN external_from JSONB;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN external_from JSONB;