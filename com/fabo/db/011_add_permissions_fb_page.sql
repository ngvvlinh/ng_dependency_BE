ALTER TABLE fb_page
    ADD COLUMN external_permissions TEXT[];
ALTER TABLE history.fb_page
    ADD COLUMN external_permissions TEXT[];