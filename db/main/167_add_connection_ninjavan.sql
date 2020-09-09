-- NinjaVan - topship
INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_subtype, connection_method, connection_provider, etop_affiliate_account, code)
VALUES (1000853619605947764, 'Topship - NinjaVan', 1, now(), now(), 'shipping/shipment/builtin/ninjavan', 'shipping', 'shipment', 'builtin', 'ninjavan', '{"user_id": "", "secret_key": ""}', 'NL21');

INSERT INTO shop_connection(connection_id, token, status, is_global, external_data, token_expires_at, created_at, updated_at)
VALUES (1000853619605947764, '', 1, true, '{"user_id": ""}', now(), now(), now());

ALTER TYPE shipping_provider ADD VALUE 'ninjavan';
