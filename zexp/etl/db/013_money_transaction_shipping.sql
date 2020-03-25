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

create table if not exists money_transaction_shipping (
  id bigint primary key,
  shop_id bigint,
  status smallint,
  total_cod integer,
  total_orders integer,
  total_amount int4,
  type TEXT,
  code text,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  closed_at timestamp with time zone,
  provider shipping_provider,
  confirmed_at timestamp with time zone,
  bank_account JSON,
  note TEXT,
  invoice_number TEXT,
  rid int8
);