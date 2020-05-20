alter table account_user alter column deleted_at TYPE timestamp with time zone USING deleted_at::timestamp with time zone;

alter table account_user drop constraint account_user_pkey;

create unique index on account_user(account_id, user_id) where deleted_at is null;