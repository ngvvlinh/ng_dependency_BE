CREATE TABLE fb_conversation (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    fb_page_id INT8 NOT NULL REFERENCES fb_page(id),
    external_link TEXT,
    external_message_count INT8,
    external_updated_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deteled_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_conversation', '{id}');