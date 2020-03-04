DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_type') THEN
        create type account_type as enum (
            'etop',
            'supplier',
            'shop',
            'partner',
            'affiliate'
        );
    END IF;
END
$$;

create table if not exists account (
    id bigint primary key,
    name text,
    type account_type,
    deleted_at timestamp with time zone,
    image_url text,
    url_slug TEXT,
    owner_id bigint,
    rid bigint
);