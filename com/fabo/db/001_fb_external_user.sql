CREATE TABLE fb_external_user (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_info jsonb,
    status INT2,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    user_id INT8 NOT NULL REFERENCES "user"(id)
);

SELECT init_history('fb_external_user', '{id}');