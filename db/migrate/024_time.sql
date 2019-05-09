ALTER TABLE "public"."shop_variant" ALTER COLUMN "created_at" SET NOT NULL;
ALTER TABLE "public"."shop_variant" ALTER COLUMN "updated_at" SET NOT NULL;

UPDATE shop_product sp SET created_at = p.created_at
  FROM product p
 WHERE sp.product_id = p.id
   AND sp.created_at IS NULL;

UPDATE shop_product sp SET updated_at = p.updated_at
  FROM product p
 WHERE sp.product_id = p.id
   AND sp.updated_at IS NULL;

ALTER TABLE "public"."product" ALTER COLUMN "created_at" SET NOT NULL;
ALTER TABLE "public"."product" ALTER COLUMN "updated_at" SET NOT NULL;
ALTER TABLE "public"."shop_product" ALTER COLUMN "created_at" SET NOT NULL;
ALTER TABLE "public"."shop_product" ALTER COLUMN "updated_at" SET NOT NULL;
ALTER TABLE "public"."variant" ALTER COLUMN "created_at" SET NOT NULL;
ALTER TABLE "public"."variant" ALTER COLUMN "updated_at" SET NOT NULL;
