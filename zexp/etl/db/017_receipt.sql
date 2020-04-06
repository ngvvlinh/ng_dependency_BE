DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'receipt_type') THEN
        create type receipt_type as enum (
            'receipt',
            'payment'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'trader_type') THEN
        create type trader_type as enum (
            'customer',
            'vendor',
            'carrier'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'receipt_created_type') THEN
        create type receipt_created_type as enum (
            'manual',
            'auto'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'receipt_ref_type') THEN
        create type receipt_ref_type as enum (
            'order',
            'fulfillment',
            'inventory_voucer',
            'purchase_order'
        );
    END IF;
END
$$;

CREATE TABLE if not exists receipt (
    id INT8 PRIMARY KEY,
    shop_id INT8,
    trader_id INT8,
    created_by INT8,
    code TEXT,
    title TEXT,
    description TEXT,
    amount INT4,
    status INT2,
    type receipt_type,
    lines JSONB,
    ref_ids INT8[],
    trader_type trader_type,
    shop_ledger_id INT8,
    ledger_id INT8,
    created_type receipt_created_type,
    ref_type receipt_ref_type,
    cancelled_reason TEXT,
    trader jsonb,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    rid INT8
);