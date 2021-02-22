ALTER TABLE "order"
    ADD COLUMN external_comment_id TEXT,
    ADD COLUMN external_post_id TEXT;
ALTER TABLE history."order"
    ADD COLUMN external_comment_id TEXT,
    ADD COLUMN external_post_id TEXT;
