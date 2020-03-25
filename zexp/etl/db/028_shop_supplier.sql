create table if not exists shop_supplier (
  id int8 primary key
, shop_id int8
, full_name text
, note text
, status int2
, created_at timestamptz
, updated_at timestamptz
, phone text
, email text
, company_name text
, tax_number text
, headquater_address text
, code text
, rid int8
);