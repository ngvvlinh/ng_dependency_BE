CREATE INDEX ON fb_external_conversation (external_page_id);

CREATE INDEX ON fb_external_post (external_page_id);

CREATE INDEX ON fb_external_message (external_conversation_id);
CREATE INDEX ON fb_external_message (external_page_id);

CREATE INDEX ON fb_external_comment (external_post_id);
CREATE INDEX ON fb_external_comment (external_created_time desc, id asc);
CREATE INDEX ON fb_external_comment (external_page_id);

CREATE INDEX ON fb_customer_conversation (external_page_id);

CREATE INDEX ON fb_external_user (external_page_id);
