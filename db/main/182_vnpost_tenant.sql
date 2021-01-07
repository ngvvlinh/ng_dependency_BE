INSERT INTO "public"."connection" ("id", "name", "status", "partner_id", "created_at", "updated_at", "deleted_at", "driver_config", "driver", "connection_type", "connection_subtype", "connection_method", "connection_provider", "etop_affiliate_account", "code") VALUES ('1000632092361806111', 'Etelecom', '1', NULL, NOW(), NOW(), NULL, NULL, 'crm/_/direct/portsip', 'telecom', NULL, 'direct', 'portsip', NULL, 'F57J');

-- shop connection for vnpost shop
INSERT INTO shop_connection(owner_id, connection_id, token, status, is_global, created_at, updated_at, telecom_data, token_expires_at, last_sync_at)
VALUES (1156517387638351348, 1000632092361806111, 'default_token', 1, false, now(), now(), '{"password": "", "username": "", "tenant_host": "", "tenant_domain": "", "tenant_token": ""}', now(), now());

-- create hotline for vnpost tenant
