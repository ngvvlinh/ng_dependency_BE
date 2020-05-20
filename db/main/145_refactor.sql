drop index order_shop_external_id_idx;

CREATE UNIQUE INDEX order_shop_external_id_idx ON "order" (shop_id, external_order_id)
  WHERE external_order_id IS NOT NULL AND partner_id IS NULL AND status != -1 AND fulfillment_shipping_status != -2
