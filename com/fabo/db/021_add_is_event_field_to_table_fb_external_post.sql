ALTER TABLE fb_external_post ADD COLUMN feed_type INT;
ALTER TABLE history.fb_external_post ADD COLUMN feed_type INT;

UPDATE fb_external_post
SET feed_type = 233 -- 233 is POST, 278 is EVENT
WHERE feed_type is NULL;
