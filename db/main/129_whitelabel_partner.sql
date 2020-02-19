-- IMPORTANT: do not run this on production

insert into "user" (
    id, status, created_at, updated_at, full_name, short_name, email, phone,
    agreed_tos_at, email_verified_at, phone_verified_at, is_test, identifying)
values (
    1055611186794473931, 1, now(), now(), 'IMGroup Owner', 'IMGroup Owner', 'imgroup@example.com', '0900000001',
    now(), now(), now(), 0, 'full')
on conflict do nothing; -- for local development only

insert into user_internal (id, hashpwd, updated_at) values
    (1055611186794473931, 'b3e761334473758bf69a23fa392c1b6581e4ebd17ec0d50794535c78cc8dc9b4b6048588', now())
on conflict do nothing; -- for local development only

insert into account (id, owner_id, name, type) values
    (1000642057714249201, 1055611186794473931, 'IMGroup', 'partner');

insert into account_user (account_id, user_id, status, roles, created_at, updated_at) values
    (1000642057714249201, 1055611186794473931, 1, '{owner}', now(), now());

insert into partner (id, name, public_name, owner_id, status, is_test, created_at, updated_at, white_label_key) values
    (1000642057714249201, 'IMGroup', 'IMGroup', 1055611186794473931, 1, 0, now(), now(), 'itopx');
