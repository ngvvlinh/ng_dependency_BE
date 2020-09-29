create table fb_customer_conversation_search (
    id int8 primary key,
    external_page_id text,
    external_user_name_norm tsvector,
    created_at timestamp with time zone
);

create table fb_external_message_search (
    id int8 primary key,
    external_page_id text,
    external_conversation_id text,
    external_message_norm tsvector,
    created_at timestamp with time zone
);

create table fb_external_comment_search (
    id int8 primary key,
    external_page_id text,
    external_post_id text,
    external_user_id text,
    external_message_norm tsvector,
    created_at timestamp with time zone
);

create index idx_fb_customer_conversation_search on fb_customer_conversation_search using gin(external_user_name_norm);
create index idx_fb_external_message_search on fb_external_message_search using gin(external_message_norm);
create index idx_fb_external_comment_search on fb_external_comment_search using gin(external_message_norm);
