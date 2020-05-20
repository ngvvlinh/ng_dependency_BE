alter table fulfillment add column if not exists provider_shipping_fee_lines jsonb;
alter table history.fulfillment add column if not exists provider_shipping_fee_lines jsonb;

alter table fulfillment drop column external_shipping_returned_fee;
alter table history.fulfillment drop column external_shipping_returned_fee;

alter table fulfillment add column if not exists shipping_fee_shop_lines jsonb;
alter table history.fulfillment add column if not exists shipping_fee_shop_lines jsonb;
