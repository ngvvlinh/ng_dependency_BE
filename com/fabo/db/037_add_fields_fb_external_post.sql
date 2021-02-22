ALTER TABLE "order"
    ADD COLUMN external_comment_id TEXT REFERENCES "fb_external_comment"(external_id),
    ADD COLUMN external_post_id TEXT REFERENCES "fb_external_post"(external_id);
ALTER TABLE history."order"
    ADD COLUMN external_comment_id TEXT,
    ADD COLUMN external_post_id TEXT;
