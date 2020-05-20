ALTER TABLE "shop_product" ADD COLUMN meta_fields JSONB;
ALTER TABLE "history"."shop_product" ADD COLUMN meta_fields JSONB;

ALTER TABLE "order_line" ADD COLUMN meta_fields JSONB;
ALTER TABLE "history"."order_line" ADD COLUMN meta_fields JSONB;
