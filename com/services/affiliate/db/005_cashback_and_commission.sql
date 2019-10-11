create table shop_cashback
(
    id                      bigint primary key,
    shop_id                 bigint      not null,
    order_id                bigint      not null,
    amount                  int         not null,
    order_created_notify_id bigint      not null,
    description             text,
    status                  int2        not null,
    valid_at                timestamptz not null,
    created_at              timestamptz not null,
    updated_at              timestamptz not null
);

create table seller_commission
(
    id             bigint      not null primary key,
    seller_id      bigint      not null,
    from_seller_id bigint,
    product_id     bigint      not null,
    shop_id        bigint      not null,
    supply_id      bigint      not null,
    amount         int8        not null,
    description    text,
    note           text,
    order_id       bigint      not null,
    status         int2        not null,
    type           commission_type,
    o_value        int         not null,
    o_base_value   int         not null,
    valid_at       timestamptz,
    created_at     timestamptz not null,
    updated_at     timestamptz not null
);
