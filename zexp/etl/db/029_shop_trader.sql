DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'trader_type') THEN
        create type trader_type as enum('customer', 'vendor', 'carrier');
    END IF;
END
$$;

create table if not exists shop_trader (
  id int8 primary key
, shop_id int8 not null
, type trader_type
, rid int8
);