CREATE TABLE fb_page_internal {
    id INT8 PRIMARY KEY REFERENCES "fb_page"(id),
    token TEXT NOT NULL,
    expires_id INT8,
    updated_at TIMESTAMP WITH TIME ZONE,
};

SELECT init_history('fb_page_internal', '{id}');