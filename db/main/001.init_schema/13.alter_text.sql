alter table history.inventory_voucher
  alter column title type text
, alter column cancel_reason type text
;

alter table public.inventory_voucher
  alter column title type text
, alter column cancel_reason type text
;
