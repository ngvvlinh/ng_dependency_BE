DROP INDEX "public"."shop_customer_shop_id_email_idx";
DROP INDEX "public"."shop_customer_shop_id_phone_idx";

CREATE INDEX ON shop_customer(shop_id,email);
CREATE INDEX ON shop_customer(shop_id,phone);
