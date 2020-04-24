CREATE TABLE fb_comment (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    parent_id INT8,
    external_parent_id TEXT,
    fb_post_id INT8 NOT NULl REFERENCES fb_post(id),
    external_message TEXT,
    external_comment_count INT8,
    external_from JSONB,
    external_attachment JSONB
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_comment', '{id}');