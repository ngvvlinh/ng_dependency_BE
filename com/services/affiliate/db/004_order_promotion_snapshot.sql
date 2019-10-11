create type order_promotion_source as enum ('etop', 'seller');

create table order_promotion
(
    id                      bigint primary key,
    product_id              bigint                 not null,
    order_id                bigint                 not null,
    base_value              int                    not null,
    amount                  int                    not null,
    unit                    unit_type              not null,
    type                    product_promotion_type not null,
    order_created_notify_id bigint,
    description             text,
    src                     order_promotion_source not null,
    created_at              timestamptz            not null,
    updated_at              timestamptz            not null,
    deleted_at              timestamptz
);

create table order_commission_setting
(
    order_id                   bigint                                not null,
    supply_id                  bigint                                not null,
    product_id                 bigint                                not null,
    level1_direct_commission   int,
    level1_indirect_commission int,
    level2_direct_commission   int,
    level2_indirect_commission int,
    depend_on                  supply_commission_setting_depend_type not null,
    level1_limit_count         int,
    level1_limit_duration      bigint,
    lifetime_duration          bigint,
    created_at                 timestamptz                           not null,
    updated_at                 timestamptz                           not null,
    primary key (order_id, supply_id, product_id)
);

create table order_created_notify
(
    id                         bigint         not null primary key,
    order_id                   bigint         not null unique,
    shop_user_id               bigint,
    seller_id                  bigint,
    shop_id                    bigint,
    supply_id                  bigint,
    referral_code              text,
    promotion_snapshot_status  int2,
    promotion_snapshot_err     text,
    commission_snapshot_status int2,
    commission_snapshot_err    text,
    cashback_process_status    int2,
    cashback_process_err       text,
    commission_process_status  int2,
    commission_process_err     text,
    payment_status             int2,
    status                     int2 default 0 not null,
    completed_at               timestamptz,
    created_at                 timestamptz    not null,
    updated_at                 timestamptz    not null
);