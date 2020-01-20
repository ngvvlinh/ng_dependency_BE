ALTER TYPE customer_type ADD VALUE 'anonymous';

update "shop_customer"
set deleted_at = now()
where "type"='independent';

INSERT INTO "public"."shop_trader" ("id", "type") VALUES ('1000570234798935240', 'customer');
INSERT INTO "public"."shop_customer" ("id", "full_name", "status", "type") VALUES ('1000570234798935240', 'Khách lẻ', 1, 'anonymous');
