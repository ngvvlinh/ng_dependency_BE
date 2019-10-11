create table shop_order_product_history
(
    user_id    bigint      not null,
    shop_id    bigint      not null,
    order_id   bigint      not null,
    supply_id  bigint      not null,
    product_id bigint      not null,
    created_at timestamptz not null,
    updated_at timestamptz not null,
    primary key (user_id, order_id, product_id)
);