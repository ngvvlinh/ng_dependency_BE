CREATE TABLE fb_external_conversation (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_user_id TEXT,
    external_user_name TEXT,
    external_link TEXT,
    external_message_count INT8,
    external_updated_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    fb_page_id INT8 NOT NULL REFERENCES fb_external_page(id)
);

SELECT init_history('fb_external_conversation', '{id}');