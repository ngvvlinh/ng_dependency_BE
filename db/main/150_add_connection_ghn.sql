ALTER TABLE connection
    ADD COLUMN "version" TEXT;
ALTER TABLE history.connection
    ADD COLUMN "version" TEXT;

-- GHN(v2) - direct
INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_subtype, connection_method, connection_provider, etop_affiliate_account, code, "version")
VALUES (1000805467932228894, 'GHN - v2', 1, now(), now(), 'shipping/shipment/direct/ghn', 'shipping', 'shipment', 'direct', 'ghn', '{ "token": "", "user_id": "", "shop_id": ""}', 'QN97', 'v2');

-- change email in external_data to identifier
-- identifier can be email or phone
UPDATE shop_connection
SET external_data = replace(external_data::text, '"email"', '"identifier"')::jsonb;

-- GHN(v2) - topship
INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_subtype, connection_method, connection_provider, etop_affiliate_account, code, "version")
VALUES (1000805467932228995, 'Topship - GHN - v2', 1, now(), now(), 'shipping/shipment/builtin/ghn', 'shipping', 'shipment', 'builtin', 'ghn', '{ "token": "", "user_id": "", "shop_id": ""}', 'CV99', 'v2');

INSERT INTO shop_connection(connection_id, token, status, is_global, external_data, created_at, updated_at)
VALUES (1000805467932228995, '', 1, true, '{"token": "", "user_id": "", "shop_id": ""}', now(), now());
