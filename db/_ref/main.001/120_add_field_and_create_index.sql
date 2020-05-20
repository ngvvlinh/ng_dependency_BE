alter table "fulfillment" add column created_by int8;
alter table history."fulfillment" add column created_by int8;

alter table inventory_voucher add column product_ids int8[];
alter table shop_stocktake add column product_ids int8[];

alter table receipt add column trader_type trader_type;
alter table receipt add column trader_phone_norm tsvector;
alter table history.receipt add column trader_phone_norm tsvector;
alter table history.receipt add column trader_type trader_type;

alter table shop_supplier add column company_name_norm tsvector;
alter table history.shop_supplier add column company_name_norm tsvector;


create index on shop_product using gin(name_norm);
create index on shop_product(list_price);

create index on shop_supplier using gin(company_name_norm);
create index on receipt using gin(trader_full_name_norm);
create index on receipt using gin(trader_phone_norm);
create index on receipt(type);
create index on receipt(trader_type);
create index on receipt(created_at);
create index on receipt(paid_at);

create index on inventory_voucher(status);
create index on inventory_voucher(ref_name);
create index on inventory_voucher(ref_code);
create index on inventory_voucher(created_at);

create index on shop_stocktake(code);
create index on shop_stocktake(status);
