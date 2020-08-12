ALTER TABLE fb_external_message
    ADD COLUMN external_from_id TEXT;
ALTER TABLE history.fb_external_message
    ADD COLUMN external_from_id TEXT;

CREATE INDEX ON fb_external_message (external_from_id);

UPDATE fb_external_message
SET external_from_id = external_from -> 'id'
WHERE external_from IS NOT NULL;

