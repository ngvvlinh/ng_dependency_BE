alter table purchase_order add column discount_lines jsonb;
alter table purchase_order add column total_fee int;
alter table purchase_order add column fee_lines jsonb;

alter table history.purchase_order add column discount_lines jsonb;
alter table history.purchase_order add column total_fee int;
alter table history.purchase_order add column fee_lines jsonb;

alter table purchase_refund add column adjustment_lines jsonb;
alter table purchase_refund add column total_adjustment int;

alter table refund add column adjustment_lines jsonb;
alter table refund add column total_adjustment int;

alter table refund drop column discount;
alter table purchase_refund drop column discount;
