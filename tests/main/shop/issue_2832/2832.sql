INSERT INTO "account" (id, name, type, url_slug, rid) VALUES
    (101, 'Etop', 'etop', 'default', 62820);

INSERT INTO "user" (id, rid, status, created_at, updated_at, full_name, short_name, email, phone, is_test, identifying, full_name_norm) VALUES
    ('1000101010101010101', 36831, 1, NOW(), NOW(), 'Etop System', 'System Short Name', 'admin@etop.vn', '0101010101', 0, 'full', to_tsvector('Etop System'));


INSERT INTO "user" (id, rid, status, created_at, updated_at, full_name, short_name, email, phone, is_test, identifying, full_name_norm) VALUES
    ('1013297053027163753', 1649, 1, NOW(), NOW(), 'test system', 'test short name', 'test1@etop.vn', '0987654321', 0, 'full', to_tsvector('test system'));

INSERT INTO "user_internal" (id, hashpwd, updated_at, rid) VALUES
    (1000101010101010101, '73d564fd83e242522e9ab67944829503b21f2de159cbf01c4c8278399523804fd5c4e2c8', now(), 62820);

INSERT INTO "account_user" (account_id, user_id, status, roles, created_at, updated_at) VALUES
    (101, 1000101010101010101, 1, '{admin}', now(), now());

INSERT INTO "account" ("id", "name", "type", "owner_id") VALUES
    ('1000015764575267699', 'eTop Trading', 'shop', '1000101010101010101');

INSERT INTO "shop" ("id", "name", "owner_id", "status", "created_at", "updated_at", "is_test", "try_on") VALUES
    ('1000015764575267699', 'eTop Trading', '1000101010101010101', 1, NOW(), NOW(), 1, 'open');
