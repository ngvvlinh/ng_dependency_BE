create table "inventory_variant" (
  "shop_id" bigint,
  "variant_id" bigint,
  "quantity_on_hand" int,
  "quantity_picked" int,
  "purchase_price" int,
  "created_at" timestamptz,
  "updated_at" timestamptz
);

CREATE UNIQUE INDEX ON "inventory_variant" (shop_id, variant_id);

CREATE TYPE "inventory_voucher_type" AS ENUM ('in', 'out');

CREATE TABLE "inventory_voucher" (
  "title" varchar,
  "shop_id" bigint,
  "id" bigint,
  "created_by" bigint,
  "updated_by" bigint,
  "status" int,
  "note" varchar,
  "trader_id" bigint,
  "total_amount" int,
  "type" inventory_voucher_type,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "confirmed_at" timestamptz,
  "cancelled_at" timestamptz,
  "cancelled_reason" varchar,
  "lines" jsonb
 );

alter table "inventory_voucher" add foreign key ("trader_id") references "shop_trader" ("id");

alter table shop add inventory_overstock boolean;

alter table history.shop add inventory_overstock boolean;



