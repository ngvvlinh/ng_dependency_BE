create type commission_type as enum ('direct', 'indirect');
create type unit_type as enum ('vnd', 'percent');
create type commission_setting_type as enum ('shop', 'affiliate');

create table commission_setting
(
    product_id bigint not null,
    account_id bigint not null,
    unit unit_type not null,
    amount int not null,
    type commission_setting_type not null,
    created_at timestamptz not null ,
    updated_at timestamptz not null,
    primary key (product_id, account_id)
);

create table commission
(
    id bigint not null primary key,
    affiliate_id bigint not null,
    value int8 default 0,
    unit unit_type default 'vnd'::unit_type,
    description text,
    note text,
    order_id bigint not null,
    status int2 default 0 not null,
    type commission_type,
    created_at timestamptz not null,
    updated_at timestamptz not null
);
