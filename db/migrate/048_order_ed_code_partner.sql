-- ed_code from shop is unique
DROP INDEX order_shop_id_ed_code_idx;
CREATE UNIQUE INDEX order_shop_id_ed_code_idx ON "order" (shop_id, ed_code)
  WHERE partner_id IS NULL and status != -1;

-- external_id from shop is unique
DROP INDEX order_shop_external_id_idx;
CREATE UNIQUE INDEX order_shop_external_id_idx ON "order" (shop_id, external_order_id)
  WHERE external_order_id IS NOT NULL AND partner_id IS NULL AND status != -1;

-- ed_code from the same partner for each shop is unique
CREATE UNIQUE INDEX order_partner_shop_id_external_code_idx ON "order" (shop_id, ed_code, partner_id)
  WHERE partner_id IS NOT NULL AND status != -1;
