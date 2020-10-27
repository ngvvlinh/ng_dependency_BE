create table fb_message_template (
    id int8 primary key,
    shop_id int8 not null references shop(id),
    template text,
    short_code text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
