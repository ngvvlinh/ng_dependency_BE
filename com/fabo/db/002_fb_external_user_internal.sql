CREATE TABLE fb_external_user_internal (
    id INT8 PRIMARY KEY REFERENCES fb_external_user(id),
    token TEXT NOT NULL,
    expires_in INT8,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_external_user_internal', '{id}');