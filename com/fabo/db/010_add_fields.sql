ALTER TABLE fb_external_comment
    ADD COLUMN external_parent JSONB;
ALTER TABLE history.fb_external_comment
    ADD COLUMN external_parent JSONB;

ALTER TABLE fb_customer_conversation
    ADD COLUMN external_comment_attachment JSONB;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN external_comment_attachment JSONB;
ALTER TABLE fb_customer_conversation
    RENAME COLUMN post_attachments TO external_post_attachments;
ALTER TABLE history.fb_customer_conversation
    RENAME COLUMN post_attachments TO external_post_attachments;
ALTER TABLE fb_customer_conversation
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN external_page_id TEXT;

update fb_customer_conversation
set external_page_id = fb_customer_conversation.external_id
from fb_external_page
where fb_external_page.id = fb_customer_conversation.fb_page_id;

ALTER TABLE fb_external_post
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_external_post
    ADD COLUMN external_page_id TEXT;

update fb_external_post
set external_page_id = fb_external_page.external_id
from fb_external_page
where fb_external_post.fb_page_id = fb_external_page.id;

ALTER TABLE fb_external_comment
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_external_comment
    ADD COLUMN external_page_id TEXT;
ALTER TABLE fb_external_comment
    ADD COLUMN external_post_id TEXT;
ALTER TABLE history.fb_external_comment
    ADD COLUMN external_post_id TEXT;

update fb_external_comment
set external_page_id = fb_external_post.external_page_id,
    external_post_id = fb_external_post.external_id
from fb_external_post
where fb_external_comment.fb_post_id = fb_external_post.id;

ALTER TABLE fb_external_conversation
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_external_conversation
    ADD COLUMN external_page_id TEXT;

update fb_external_conversation
set external_page_id = fb_external_page.external_id
from fb_external_page
where fb_external_page.id = fb_external_conversation.fb_page_id;

ALTER TABLE fb_external_message
    ADD COLUMN external_page_id TEXT;
ALTER TABLE history.fb_external_message
    ADD COLUMN external_page_id TEXT;

update fb_external_message
set external_page_id = fb_external_page.external_id
from fb_external_page
where fb_external_page.id = fb_external_message.fb_page_id;


