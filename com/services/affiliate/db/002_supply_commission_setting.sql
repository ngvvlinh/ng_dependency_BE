create type supply_commission_setting_depend_type as enum ('product', 'customer');

create table supply_commission_setting (
    shop_id bigint not null,
    product_id bigint not null,
    level1_direct_commission int not null,
    level1_indirect_commission int not null,
    level2_direct_commission int not null,
    level2_indirect_commission int not null,
    depend_on supply_commission_setting_depend_type not null,
    level1_limit_count int not null,
    level1_limit_duration bigint not null,
    m_level1_limit_duration jsonb,
    lifetime_duration bigint not null,
    m_lifetime_duration jsonb,
    created_at timestamptz not null,
    updated_at timestamptz not null
)