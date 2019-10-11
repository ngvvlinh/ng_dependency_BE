create type commission_type as enum ('direct', 'indirect');
create type unit_type as enum ('vnd', 'percent');
create type commission_setting_type as enum ('shop', 'affiliate');
create type product_promotion_type as enum ('cashback', 'discount');

create table commission_setting
(
    product_id bigint                  not null,
    account_id bigint                  not null,
    unit       unit_type               not null,
    amount     int                     not null,
    type       commission_setting_type not null,
    created_at timestamptz             not null,
    updated_at timestamptz             not null,
    primary key (product_id, account_id)
);

create table product_promotion
(
    id          bigint                 not null primary key,
    product_id  bigint                 not null,
    shop_id     bigint                 not null,
    amount      int                    not null,
    unit        unit_type              not null,
    code        text,
    description text,
    note        text,
    type        product_promotion_type not null,
    status      int2                   not null,
    created_at  timestamptz            not null,
    updated_at  timestamptz            not null
);
create index on product_promotion (product_id);

create table affiliate_referral_code
(
    id           bigint primary key,
    code         text        not null,
    affiliate_id bigint      not null,
    created_at   timestamptz not null,
    updated_at   timestamptz not null,
    deleted_at   timestamptz
);

create unique index on affiliate_referral_code (code);

create table user_referral
(
    user_id            bigint primary key,
    referral_id        bigint,
    referral_code      text,
    sale_referral_id   bigint,
    sale_referral_code text,
    referral_at        timestamptz,
    sale_referral_at   timestamptz,
    created_at         timestamptz not null,
    updated_at         timestamptz not null
)