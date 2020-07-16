ALTER TABLE fb_external_comment
    ADD COLUMN source int;
ALTER TABLE history.fb_external_comment
    ADD COLUMN source int;

-- Delete all customer_conversation has type "comment" and external_user_id = external_page_id
DELETE FROM fb_customer_conversation_state
WHERE id in (
    SELECT id
    FROM fb_customer_conversation
    WHERE "type" = 90 AND external_user_id = external_page_id
);

DELETE FROM fb_customer_conversation
WHERE "type" = 90 AND external_user_id = external_page_id;
