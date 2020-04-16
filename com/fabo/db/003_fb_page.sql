CREATE TABLE fb_page {
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    shop_id INT8 NOT NULL REFERENCES "shop"(id),
    user_id INT8 NOT NULL REFERENCES "user"(id),
    name TEXT,
    tasks []TEXT,
    category TEXT,
    category_list jsonb,
    status INT2,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
};

SELECT init_history('fb_page', '{id}');