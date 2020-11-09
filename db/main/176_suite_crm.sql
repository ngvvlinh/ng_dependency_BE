INSERT INTO connection(id, name, status, created_at, updated_at, driver, connection_type, connection_method, connection_provider, code)
VALUES (1000853692605949190, 'CRM - SuiteCRM', 1, now(), now(), 'crm/builtin/suitecrm', 'crm', 'builtin', 'suitecrm', 'PN02');

INSERT INTO shop_connection(connection_id, token, status, is_global, created_at, updated_at)
VALUES (1000853692605949190, '', 1, true, now(), now());
