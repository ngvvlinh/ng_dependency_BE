DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'fulfillment_endpoint') THEN
        create type fulfillment_endpoint as enum (
            'supplier',
            'shop',
            'customer'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'shipping_provider') THEN
        create type shipping_provider as enum (
            'ghn',
            'ghtk',
            'vtpost'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'shipping_state') THEN
        create type shipping_state as enum (
            'default',       -- have not created shipping order on GHN yet
            'unknown',
            'created',       -- created shipping order on GHN
            'confirmed',     -- unused
            'picking',
            'processing',    -- unused
            'holding',
            'returning',
            'returned',
            'delivering',
            'delivered',
            'undeliverable',
            'cancelled',
            'closed'
        );
    END IF;
END
$$;

CREATE TABLE if not exists fulfillment (
    id bigint primary key,
    order_id bigint,
    lines jsonb,
    variant_ids bigint[],
    type_from fulfillment_endpoint,
    type_to   fulfillment_endpoint,
    address_from jsonb,
    address_to   jsonb,
    shop_id bigint,
    supplier_id bigint,
    supplier_confirm smallint,
    shop_confirm smallint,
    total_items integer,
    total_weight integer,
    total_amount integer,
    total_cod_amount integer,
    shipping_fee_customer integer,
    shipping_fee_shop integer,
    external_shipping_fee integer,
    created_at timestamp with time zone,
--     created_by bigint,
    updated_at timestamp with time zone,
    shipping_delivered_at timestamp with time zone,
    shipping_returned_at timestamp with time zone,
--     updated_by bigint,
    shipping_cancelled_at timestamp with time zone,
    cancel_reason text,
--     cancelled_by bigint,
    closed_at timestamp with time zone,
    shipping_provider shipping_provider,
    shipping_code text,
    shipping_note text,
    external_shipping_id text,
    external_shipping_code text,
    external_shipping_service_id text,
    external_shipping_created_at timestamp with time zone,
    external_shipping_updated_at timestamp with time zone,
    external_shipping_cancelled_at timestamp with time zone,
    external_shipping_delivered_at timestamp with time zone,
    external_shipping_returned_at timestamp with time zone,
    external_shipping_closed_at timestamp with time zone,
    external_shipping_state text,
    external_shipping_status smallint,
    external_shipping_data jsonb,
    shipping_state shipping_state,
    status smallint,
    sync_status smallint,
    sync_states jsonb,
    etop_fee_adjustment integer,
    last_sync_at timestamptz,
    expected_delivery_at TIMESTAMP WITH TIME ZONE,
    money_transaction_id bigint,
    cod_etop_transfered_at timestamp with time zone,
    provider_shipping_fee_lines jsonb,
    shipping_fee_shop_lines jsonb,
    etop_discount integer,
    shipping_status smallint,
    provider_service_id text,
    etop_payment_status int2,
    address_to_province_code text,
    address_to_district_code text,
    address_to_ward_code text,
    expected_pick_at timestamptz,
    external_shipping_state_code text,
    confirm_status int2,
    shipping_fee_main integer,
    shipping_fee_return integer,
    shipping_fee_insurance integer,
    shipping_fee_adjustment integer,
    shipping_fee_cods integer,
    shipping_fee_info_change integer,
    shipping_fee_other integer,
    money_transaction_shipping_external_id bigint,
    total_discount int4,
    basket_value int4,
    external_shipping_logs jsonb,
    partner_id int8,
    external_shipping_note text,
    external_shipping_sub_state text,
    try_on try_on,
    admin_note text,
    is_partial_delivery boolean,
    shipping_fee_discount integer,
    shipping_created_at timestamp with time zone,
    shipping_picking_at timestamp with time zone,
    shipping_holding_at timestamp with time zone,
    shipping_delivering_at timestamp with time zone,
    shipping_returning_at timestamp with time zone,
    etop_adjusted_shipping_fee_main int,
    etop_price_rule boolean,
    actual_compensation_amount integer,
    delivery_route text,
    shipping_type int2,
    connection_id int8,
    connection_method text,
    shop_carrier_id int8,
    shipping_service_name text,
    gross_weight int,
    chargeable_weight int,
    length int,
    width int,
    height int,
    external_affiliate_id text,
    original_cod_amount int4,
    shipping_service_fee int4,
    address_return jsonb,
    shipping_fee_shop_transfered_at timestamp with time zone,
    include_insurance boolean,
    external_shipping_name text,
    created_by int8,

    rid bigint
);