CREATE TABLE fb_customer_conversation_state (
    id INT8 PRIMARY KEY REFERENCES fb_customer_conversation(id),
    is_read bool,
    external_page_id TEXT,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_customer_conversation_state', '{id}');

INSERT INTO fb_customer_conversation_state(id, is_read, external_page_id, updated_at)
SELECT id, is_read, external_page_id, updated_at
FROM fb_customer_conversation;

ALTER TABLE fb_customer_conversation
DROP COLUMN is_read;

ALTER TABLE history.fb_customer_conversation
DROP COLUMN is_read;