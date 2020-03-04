DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_source_type') THEN
        create type order_source_type as enum (
            'unknown',
            'self',
            'import',
            'api',
            'etop_pos',
            'etop_pxs',
            'etop_cmx',
            'ts_app',
            'etop_app',
            'haravan'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ghn_note_code') THEN
        create type ghn_note_code as enum (
            'CHOTHUHANG',
            'CHOXEMHANGKHONGTHU',
            'KHONGCHOXEMHANG'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'try_on') THEN
        create type try_on as enum (
            'none',
            'open',
            'try'
        );
    END IF;
END
$$;

CREATE TABLE if not exists "order" (
    id bigint primary key,
    rid bigint,
    shop_id bigint,
    shop_name text,
    code text,
    product_ids bigint[],
    variant_ids bigint[],
    supplier_ids bigint[],
    currency text,
    payment_method text,
    customer jsonb,
    customer_address jsonb,
    billing_address jsonb,
    shipping_address jsonb,
    shop_address jsonb,
    customer_phone text,
    customer_email text,
    created_at timestamp with time zone,
    processed_at timestamp with time zone,
    updated_at timestamp with time zone,
    closed_at timestamp with time zone,
    confirmed_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    cancel_reason text,
    customer_confirm smallint,
    external_confirm smallint,
    shop_confirm smallint,
    confirm_status smallint,
    fulfillment_shipping_status smallint,
    etop_payment_status SMALLINT,
    status smallint,
    lines jsonb,
    discounts jsonb,
    total_items integer,
    basket_value integer,
    total_weight integer,
    total_tax integer,
    total_discount integer,
    total_amount integer,
    fulfillment_shipping_states  TEXT[],
    fulfillment_payment_statuses INT2[],
    fulfillment_statuses INT2[],
    order_discount INT4,

    order_note text,
    shop_note text,
    shipping_note text,
    order_source_id int8,
    order_source_type order_source_type,
    external_order_id text,
    ed_code text,
    external_url text,
    shop_shipping jsonb,

    customer_name text,
    shop_shipping_fee int4,
    total_fee int4,
    fee_lines JSONB,
    shop_cod int4,
    reference_url text,
    is_outside_etop boolean,
    ghn_note_code ghn_note_code,
    try_on try_on,
    customer_name_norm tsvector,
    product_name_norm tsvector,
    fulfillment_type int2,
    fulfillment_ids int8[],
    external_meta JSONB,
    trading_shop_id int8,
    payment_status int2,
    payment_id int8,
    referral_meta JSONB,
    customer_id int8,
    created_by int8,

    partner_id int8
);