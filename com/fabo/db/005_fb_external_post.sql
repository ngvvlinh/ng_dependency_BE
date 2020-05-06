CREATE TABLE fb_external_post (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_parent_id TEXT,
    external_from JSONB,
    external_picture TEXT,
    external_icon TEXT,
    external_message TEXT,
    external_attachments JSONB,
    external_created_time TIMESTAMP WITH TIME ZONE,
    external_updated_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    fb_page_id INT8 NOT NULL REFERENCES fb_external_page(id)
);

SELECT init_history('fb_external_post', '{id}');