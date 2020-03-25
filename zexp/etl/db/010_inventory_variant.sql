create table if not exists "inventory_variant" (
    shop_id bigint,
    variant_id bigint,
    quantity_on_hand int,
    quantity_picked int,
    cost_price int,
    created_at timestamptz,
    updated_at timestamptz,
    rid int8,
    primary key (shop_id, variant_id)
);