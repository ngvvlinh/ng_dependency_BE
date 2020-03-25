CREATE TABLE IF NOT EXISTS "user" (
    id bigint primary key,
    full_name text,
    short_name text,
    email text,
    phone text,
    status int2,

    created_at timestamp with time zone,
    updated_at timestamp with time zone,

    agreed_tos_at timestamp with time zone,
    agreed_email_info_at timestamp with time zone,
    email_verified_at timestamp with time zone,
    phone_verified_at timestamp with time zone,

    email_verification_sent_at timestamp with time zone,
    phone_verification_sent_at timestamp with time zone,

    rid bigint
);