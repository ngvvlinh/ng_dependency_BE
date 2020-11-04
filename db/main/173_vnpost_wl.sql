insert into "user" (
    id, status, created_at, updated_at, full_name, short_name, email, phone,
    agreed_tos_at, email_verified_at, phone_verified_at, is_test, identifying)
values (
    1156517387638351348, 1, now(), now(), 'VNPost Owner', 'VNPost Owner', 'hi@vnpost.vn', '0900000002',
    now(), now(), now(), 0, 'full')
on conflict do nothing;

insert into account (id, owner_id, name, type) values
    (1156518020386448488, 1156517387638351348, 'VNPost', 'partner');

insert into account_user (account_id, user_id, status, roles, created_at, updated_at) values
    (1156518020386448488, 1156517387638351348, 1, '{owner}', now(), now());

insert into partner (id, name, public_name, owner_id, status, is_test, created_at, updated_at, white_label_key) values
    (1156518020386448488, 'VNPost', 'VNPost', 1156517387638351348, 1, 0, now(), now(), 'vnpost');
