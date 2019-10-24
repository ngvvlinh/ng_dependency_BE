create table customer_policy_group
(
    id         bigint primary key,
    supply_id  bigint      not null,
    name       text,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

alter table supply_commission_setting
    add column customer_policy_group_id bigint,
    add column "group"                  text;
alter table shop_order_product_history
    add column customer_policy_group_id bigint;
alter table order_commission_setting
    add column "customer_policy_group_id" bigint,
    add column "group"                    text;