CREATE TABLE fb_external_user_connected (
    external_id TEXT PRIMARY KEY,
    external_info jsonb,
    status INT2,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    shop_id INT8 REFERENCES shop(id)
);

SELECT init_history('fb_external_user_connected', '{external_id}');

-- copy fb_external_users have records in fb_external_user_internals
INSERT INTO fb_external_user_connected(external_id, external_info, status, created_at, updated_at)
SELECT fu.external_id, fu.external_info, fu.status, fu.created_at, fu.updated_at
FROM fb_external_user fu
JOIN fb_external_user_internal fui ON fu.external_id = fui.external_id;

ALTER TABLE fb_external_user_internal
    DROP CONSTRAINT fb_external_user_internal_external_id_fkey;

ALTER TABLE fb_external_user_internal
    ADD FOREIGN KEY (external_id) REFERENCES fb_external_user_connected(external_id);

AlTER TABLE fb_external_page
    ADD COLUMN external_user_id TEXT REFERENCES fb_external_user_connected(external_id);

ALTER TABLE history.fb_external_page
    ADD COLUMN external_user_id TEXT;
