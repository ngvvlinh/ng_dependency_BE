create type contact_type as enum ('phone', 'email');
create type gender_type as enum ('male', 'female', 'other');
create type customer_type as enum ('individual', 'organization');

create table shop_trader (
  id int8 primary key
, shop_id int8 not null
);

create table shop_customer (
  id int8 primary key references shop_trader(id)
, shop_id int8 not null references shop(id)
, code text
, full_name text
, gender gender_type
, type customer_type
, birthday date
, note text
, phone text
, email text
, status int2 not null
, created_at timestamptz
, updated_at timestamptz
, deleted_at timestamptz
);

create unique index on shop_customer(shop_id, phone)
    where deleted_at is null and phone is not null;

create unique index on shop_customer(shop_id, email)
    where deleted_at is null and email is not null;

create table shop_vendor (
  id int8 primary key references shop_trader(id)
, shop_id int8 not null references shop(id)
, full_name text
, note text
, status int2
, created_at timestamptz
, updated_at timestamptz
, deleted_at timestamptz
);

create table shop_trader_address (
  id int8 primary key
, shop_id int8 not null references shop(id)
, trader_id int8 not null references shop_trader(id)
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
, status int2 not null
, coordinates jsonb
, created_at timestamptz
, updated_at timestamptz
, deleted_at timestamptz
);

select init_history('shop_trader', '{id,shop_id}');
select init_history('shop_customer', '{id,shop_id}');
select init_history('shop_vendor', '{id,shop_id}');
select init_history('shop_trader_address', '{id,shop_id,trader_id}');
