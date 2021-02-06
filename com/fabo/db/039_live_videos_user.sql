ALTER TABLE fb_external_post
    ADD COLUMN "type" INT,
    ADD COLUMN external_user_id TEXT;
ALTER TABLE history.fb_external_post
    ADD COLUMN "type" INT,
    ADD COLUMN external_user_id TEXT;

ALTER TABLE fb_external_comment
    ADD COLUMN post_type INT,
    ADD COLUMN external_owner_post_id TEXT;
ALTER TABLE history.fb_external_comment
    ADD COLUMN post_type INT,
    ADD COLUMN external_owner_post_id TEXT;

ALTER TABLE fb_customer_conversation
    ADD COLUMN external_owner_post_id TEXT;
ALTER TABLE history.fb_customer_conversation
    ADD COLUMN external_owner_post_id TEXT;
