CREATE TABLE fb_customer_conversation (
    id INT8 PRIMARY KEY,
    fb_page_id INT8 NOT NULL REFERENCES fb_external_page(id),
    external_id TEXT, -- post_id || conversation_id
    external_user_id TEXT,
    external_user_name TEXT,
    post_attachments JSONB,
    type INT2,
    is_read bool,
    last_message TEXT,
    last_message_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_customer_conversation', '{id}');