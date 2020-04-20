CREATE TABLE fb_user_internal (
    id INT8 PRIMARY KEY REFERENCES "fb_user"(id),
    token TEXT NOT NULL,
    expires_in INT8,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_user_internal', '{id}');