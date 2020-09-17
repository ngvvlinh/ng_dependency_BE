-- DHL - topship
INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_subtype, connection_method, connection_provider, etop_affiliate_account, code)
VALUES (1000853619759088766, 'Topship - DHL', 1, now(), now(), 'shipping/shipment/builtin/dhl', 'shipping', 'shipment', 'builtin', 'dhl', '{"user_id": "", "secret_key": "", "shop_id": ""}', 'DN95');

INSERT INTO shop_connection(connection_id, token, status, is_global, external_data, token_expires_at, created_at, updated_at)
VALUES (1000853619759088766, 'token_auto_gen', 1, true, '{"shop_id": ""}', now(), now(), now());

ALTER TYPE shipping_provider ADD VALUE 'dhl';
