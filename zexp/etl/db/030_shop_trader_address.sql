create table if not exists shop_trader_address (
  id int8 primary key
, shop_id int8
, trader_id int8
, full_name text
, phone text
, email text
, company text
, district_code text
, ward_code text
, city text
, address1 text
, address2 text
, position text
, note text
, "primary" boolean
, status int2
, coordinates jsonb
, is_default boolean
, created_at timestamptz
, updated_at timestamptz
, rid int8
);