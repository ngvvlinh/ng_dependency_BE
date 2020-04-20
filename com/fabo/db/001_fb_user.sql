CREATE TABLE fb_user (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    user_id INT8 NOT NULL REFERENCES "user"(id),
    external_info jsonb,
    status INT2,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_user', '{id}');