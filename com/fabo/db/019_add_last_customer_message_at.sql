ALTER TABLE fb_customer_conversation
    ADD COLUMN last_customer_message_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN last_customer_message_at TIMESTAMP WITH TIME ZONE;

-- migrate last_customer_message_at for message
-- update last_customer_message_at where external_from.id <> b.external_page_id
UPDATE fb_customer_conversation as a
SET last_customer_message_at = (
		SELECT external_created_time
		FROM fb_external_message as b
		where a.external_id = b.external_conversation_id and b.external_from->>'id' <> b.external_page_id
		order by b.external_created_time desc
		limit 1
	)
WHERE a.type = 872;
