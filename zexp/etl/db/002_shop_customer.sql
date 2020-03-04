DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_type') THEN
        create type gender_type as enum ('male', 'female', 'other');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'customer_type') THEN
        create type customer_type as enum ('individual', 'organization', 'anonymous', 'independent');
    END IF;
END
$$;

create table if not exists shop_customer (
    id bigint primary key,
    shop_id bigint,
    code text,
    full_name text,
    gender gender_type,
    type customer_type,
    birthday date,
    note text,
    phone text,
    email text,
    status int2,
    code_norm int4,
    full_name_norm tsvector,
    phone_norm tsvector,
    external_id TEXT,
    external_code TEXT,
    partner_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    rid bigint
);