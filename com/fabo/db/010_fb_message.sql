CREATE TABLE fb_message (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    fb_conversation_id INT8 NOT NULL REFERENCES fb_conversation(id),
    external_message TEXT,
    external_to JSONB,
    external_from JSONB,
    external_attachments JSONB,
    external_created_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_message', '{id}');