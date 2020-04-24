CREATE TABLE fb_post (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    fb_page_id INT8 NOT NULL REFERENCES fb_page(id),
    parent_id INT8,
    external_parent_id TEXT,
    external_from JSONB,
    external_picture TEXT,
    external_icon TEXT,
    external_message TEXT,
    external_attachments JSONB,
    external_created_time TIMESTAMP WITH TIME ZONE,
    external_updated_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_page', '{id}');