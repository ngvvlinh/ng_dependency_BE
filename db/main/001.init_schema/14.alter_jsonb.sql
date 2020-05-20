alter table history.money_transaction_shipping
  alter column bank_account type jsonb
;

alter table history.money_transaction_shipping_external
  alter column bank_account type jsonb
;

alter table public.money_transaction_shipping
  alter column bank_account type jsonb
;

alter table public.money_transaction_shipping_external
  alter column bank_account type jsonb
;
