CREATE TABLE fb_external_message (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_conversation_id TEXT,
    external_message TEXT,
    external_to JSONB,
    external_from JSONB,
    external_attachments JSONB,
    external_created_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    fb_conversation_id INT8 NOT NULL REFERENCES fb_external_conversation(id),
    fb_page_id INT8 NOT NUlL REFERENCES fb_external_page(id)
);

SELECT init_history('fb_external_message', '{id}');