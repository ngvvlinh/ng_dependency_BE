ALTER TABLE fb_external_post
    ADD COLUMN status_type INT,
    ADD COLUMN total_comments INT,
    ADD COLUMN total_reactions INT;
ALTER TABLE history.fb_external_post
    ADD COLUMN status_type INT,
    ADD COLUMN total_comments INT,
    ADD COLUMN total_reactions INT;
