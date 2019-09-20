create type product_promotion_type as enum ('cashback', 'discount');

create table product_promotion
(
    id          bigint                    not null primary key,
    product_id  bigint                    not null,
    shop_id     bigint                    not null,
    amount      int                       not null,
    unit        unit_type                 not null,
    code        text,
    description text,
    note        text,
    type        product_promotion_type    not null,
    status      int2                  not null,
    created_at  timestamptz not null,
    updated_at  timestamptz not null
);
create index on product_promotion (product_id);