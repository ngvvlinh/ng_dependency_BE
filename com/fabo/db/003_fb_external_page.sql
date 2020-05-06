CREATE TABLE fb_external_page (
    id INT8 PRIMARY KEY,
    external_id TEXT NOT NULL,
    external_name TEXT,
    external_tasks TEXT[],
    external_category TEXT,
    external_category_list jsonb,
    external_permissions TEXT[],
    external_image_url TEXT,
    connection_status INT2,
    status INT2,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    fb_user_id INT8 NOT NULL REFERENCES fb_external_user(id),
    shop_id INT8 NOT NULL REFERENCES "shop"(id),
    user_id INT8 NOT NULL REFERENCES "user"(id)
);

SELECT init_history('fb_external_page', '{id}');