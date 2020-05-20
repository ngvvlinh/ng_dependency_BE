
alter table history.inventory_voucher
  alter column status type smallint
;

alter table history.purchase_refund
  alter column status type smallint
;

alter table history.refund
  alter column status type smallint
;

alter table history.shipnow_fulfillment
  alter column status type smallint
,  alter column sync_status type smallint
,  alter column confirm_status type smallint
, alter column shipping_status type smallint
,  alter column etop_payment_status type smallint
;

alter table history.shop_stocktake
  alter column status type smallint
;

alter table public.inventory_voucher
  alter column status type smallint
;

alter table public.purchase_refund
  alter column status type smallint
;

alter table public.refund
  alter column status type smallint
;

alter table public.shipnow_fulfillment
  alter column status type smallint
, alter column sync_status type smallint
, alter column confirm_status type smallint
, alter column shipping_status type smallint
, alter column etop_payment_status type smallint
;

alter table public.shop_stocktake
  alter column status type smallint
;
