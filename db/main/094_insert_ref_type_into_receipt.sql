create type receipt_ref_type as enum('order', 'fulfillment', 'inventory_voucher');

alter table receipt add column ref_type receipt_ref_type;
alter table "history".receipt add column ref_type receipt_ref_type;

update receipt
set ref_type = 'order'
where shop_id <> 1000030662086749358;

update receipt
set ref_type = 'fulfillment'
where shop_id = 1000030662086749358;
