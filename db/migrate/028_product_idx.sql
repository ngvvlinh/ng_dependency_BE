DROP INDEX "public"."product_ed_code_product_source_id_idx";

CREATE UNIQUE INDEX "product_ed_code_product_source_id_idx" ON "public"."product" USING BTREE ("ed_code","product_source_id") WHERE deleted_at IS NULL;

DROP INDEX "public"."variant_ed_code_product_source_id_idx";

CREATE UNIQUE INDEX "variant_ed_code_product_source_id_idx" ON "public"."variant" USING BTREE ("ed_code","product_source_id") WHERE deleted_at IS NULL;

CREATE INDEX ON shop (product_source_id);
CREATE INDEX ON product (product_source_id);
CREATE INDEX ON variant (product_source_id);

ALTER TABLE product
  ADD COLUMN IF NOT EXISTS unit TEXT,
  DROP COLUMN IF EXISTS subcode,
  DROP COLUMN IF EXISTS sku;
ALTER TABLE history.product
  ADD COLUMN IF NOT EXISTS unit TEXT,
  DROP COLUMN IF EXISTS subcode,
  DROP COLUMN IF EXISTS sku;
