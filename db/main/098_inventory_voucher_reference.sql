alter table inventory_voucher add column "ref_name" text;
alter table inventory_voucher add column "ref_id" bigint;
alter table inventory_voucher add column "ref_type" text;

alter table inventory_voucher rename column "cancelled_reason" to "cancel_reason";

create unique index on inventory_voucher (id);
