alter table public.payment
  alter column created_at type timestamp with time zone
, alter column updated_at type timestamp with time zone
;

alter table history.payment
  alter column created_at type timestamp with time zone
, alter column updated_at type timestamp with time zone
;
