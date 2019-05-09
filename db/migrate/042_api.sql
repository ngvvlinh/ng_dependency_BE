ALTER TYPE order_source_type ADD VALUE 'api' AFTER 'import';

CREATE TYPE try_on AS ENUM ('none', 'open', 'try');

ALTER TABLE fulfillment ADD COLUMN try_on try_on;
ALTER TABLE history.fulfillment ADD COLUMN try_on try_on;
ALTER TABLE "order" ADD COLUMN try_on try_on;
ALTER TABLE history."order" ADD COLUMN try_on try_on;
ALTER TABLE "shop" ADD COLUMN try_on try_on;
ALTER TABLE history."shop" ADD COLUMN try_on try_on;

CREATE INDEX order_external_id_idx ON "order" (external_order_id) WHERE external_order_id IS NOT NULL;

CREATE UNIQUE INDEX order_partner_external_id_idx ON "order" (partner_id, external_order_id)
  WHERE external_order_id IS NOT NULL AND partner_id IS NOT NULL AND status != -1;

CREATE UNIQUE INDEX order_shop_external_id_idx ON "order" (shop_id, external_order_id)
  WHERE external_order_id IS NOT NULL AND partner_id IS NULL AND status != -1;
