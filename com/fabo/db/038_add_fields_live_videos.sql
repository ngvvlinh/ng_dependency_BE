ALTER TABLE fb_external_post
    ADD COLUMN is_live_video bool,
    ADD COLUMN external_live_video_status TEXT,
    ADD COLUMN live_video_status TEXT;
ALTER TABLE history.fb_external_post
    ADD COLUMN is_live_video bool,
    ADD COLUMN external_live_video_status TEXT,
    ADD COLUMN live_video_status TEXT;

CREATE INDEX ON fb_customer_conversation(type);
