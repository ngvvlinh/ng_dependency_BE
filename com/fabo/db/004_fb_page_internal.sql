CREATE TABLE fb_page_internal (
    id INT8 PRIMARY KEY REFERENCES "fb_page"(id),
    token TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('fb_page_internal', '{id}');