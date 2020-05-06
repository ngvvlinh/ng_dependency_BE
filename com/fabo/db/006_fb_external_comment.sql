CREATE TABLE fb_external_comment (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_user_id TEXT,
    external_parent_id TEXT,
    external_parent_user_id TEXT,
    external_message TEXT,
    external_comment_count INT8,
    external_from JSONB,
    external_attachment JSONB,
    external_created_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    fb_post_id INT8 NOT NULl REFERENCES fb_external_post(id),
    fb_page_id INT8 NOT NULL REFERENCES fb_external_page(id)
);

SELECT init_history('fb_external_comment', '{id}');