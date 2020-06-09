ALTER TABLE fb_external_user
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_external_user
    ADD COLUMN external_page_id TEXT;

UPDATE fb_external_user
SET external_page_id = (
	SELECT external_page_id
	FROM fb_customer_conversation
	WHERE external_user_id = fb_external_user.external_id
	LIMIT 1
);
