ALTER TABLE "order"
    ADD COLUMN external_comment_id TEXT REFERENCES "fb_external_comment"(external_id),
    ADD COLUMN external_post_id TEXT REFERENCES "fb_external_post"(external_id);
ALTER TABLE history."order"
    ADD COLUMN external_comment_id TEXT,
    ADD COLUMN external_post_id TEXT;

ALTER TABLE fb_external_post
    ADD COLUMN status_type INT,
    ADD COLUMN total_comments INT,
    ADD COLUMN total_reactions INT;
ALTER TABLE history.fb_external_post
    ADD COLUMN status_type INT,
    ADD COLUMN total_comments INT,
    ADD COLUMN total_reactions INT;
