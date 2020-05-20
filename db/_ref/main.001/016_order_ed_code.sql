CREATE UNIQUE INDEX ON "order" (shop_id, ed_code) WHERE shop_confirm != -1;
ALTER TABLE "order" ALTER COLUMN "shop_confirm" SET NOT NULL;
