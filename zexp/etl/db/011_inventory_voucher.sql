DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'inventory_voucher_type') THEN
        create type inventory_voucher_type as enum (
            'in',
            'out'
        );
    END IF;
END
$$;

CREATE TABLE if not exists inventory_voucher (
  id bigint primary key,
  shop_id bigint,
  title varchar,
  created_by bigint,
  updated_by bigint,
  status int,
  trader_id bigint,
  total_amount int,
  type inventory_voucher_type,
  created_at timestamptz,
  updated_at timestamptz,
  confirmed_at timestamptz,
  cancelled_at timestamptz,
  cancel_reason varchar,
  lines jsonb,
  ref_id bigint,
  ref_type text,
  ref_code text,
  code text,
  product_ids int8[],
  variant_ids int8[],
  trader jsonb,
  rid int8
);