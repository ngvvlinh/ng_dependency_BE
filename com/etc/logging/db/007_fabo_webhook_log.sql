create table fb_webhook_log
(
    id          BIGINT PRIMARY KEY,
    data        JSONB,
    type        VARCHAR(20),
    page_id     TEXT NULL,
    error       TEXT NULL,
    external_id TEXT NULL,
    created_at  TIMESTAMP WITH TIME ZONE,
    updated_at  TIMESTAMP WITH TIME ZONE
);
