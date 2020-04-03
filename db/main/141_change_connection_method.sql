UPDATE connection SET connection_method = 'builtin' WHERE connection_method = 'topship';

UPDATE fulfillment SET connection_method = 'builtin' WHERE connection_method = 'topship';

ALTER TABLE shipment_service
    ADD COLUMN available_locations JSONB
    , ADD COLUMN blacklist_locations JSONB
    , ADD COLUMN other_condition JSONB;
