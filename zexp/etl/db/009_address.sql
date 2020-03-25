DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'address_type') THEN
        create type address_type as enum (
            'billing',
            'shipping',
            'general',
            'warehouse',
            'shipfrom',
            'shipto'
        );
    END IF;
END
$$;

CREATE TABLE if not exists address (
    id bigint primary key,
    country text,
    province_code text,
    province text,
    district_code text,
    district text,
    ward text,
    ward_code text,
    address1 text,
    address2 text,
    is_default boolean,
    type address_type,
    account_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    full_name text,
    first_name text,
    last_name text,
    email text,
    position text,
    city text,
    zip text,
    phone text,
    company text,
    coordinates jsonb,
    notes jsonb,
    rid int8
);