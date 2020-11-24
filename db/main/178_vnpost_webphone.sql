-- create default shop for vnpost
INSERT INTO user_internal (id, hashpwd, updated_at) VALUES (1156517387638351348, '', NOW());

INSERT INTO account (id, owner_id, name, type) VALUES
    (1134164111674536521, 1156517387638351348, 'VNPost', 'shop');

INSERT INTO shop (id, owner_id, name, status, created_at, updated_at, is_test, try_on) VALUES (1134164111674536521, 1156517387638351348, 'VietNam Post', 1, NOW(), NOW(), 0, 'open');

INSERT INTO account_user ("account_id", "user_id", "status", "created_at", "updated_at", "roles") VALUES ('1134164111674536521', '1156517387638351348', '1', NOW(), NOW(), '{owner}');
