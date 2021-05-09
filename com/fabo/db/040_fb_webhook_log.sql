ALTER TABLE fb_webhook_log
    ADD COLUMN external_user_id TEXT;
    ADD COLUMN external_page_id TEXT;

UPDATE fb_webhook_log
SET external_page_id = page_id;

ALTER TABLE fb_external_log
    DROP COLUMN page_id;
