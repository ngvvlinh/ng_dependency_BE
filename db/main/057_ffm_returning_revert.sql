/*
Treat returning as -2 is more useful because most of our logic groups
"returning" and "returned" as "Trả hàng".

 - Summary can query by status=2 (đơn chưa đối soát)
 - Order code can ignore status=-1 and status=-2 (huỷ và trả hàng)
*/
-- Override 055_status.sql --
CREATE OR REPLACE FUNCTION shipping_state_to_shipping_status(state shipping_state) RETURNS INT2
LANGUAGE plpgsql AS $$
BEGIN
  IF (state = 'default' OR state = 'created') THEN
    RETURN  0;
  ELSIF (state = 'cancelled') THEN
    RETURN -1;
  ELSIF (state = 'returning' OR state = 'returned') THEN
    RETURN -2;
  ELSIF (state = 'delivered') THEN
    RETURN  1;
  ELSE
    RETURN  2;
  END IF;
END;
$$;

-- update old rows
SELECT * FROM fulfillment WHERE shipping_status != shipping_state_to_shipping_status(shipping_state) ORDER BY created_at DESC;
UPDATE fulfillment SET shipping_status = shipping_state_to_shipping_status(shipping_state)
  WHERE shipping_status != shipping_state_to_shipping_status(shipping_state);

-- overwrite 048_order_ed_code_partner
-- ignore unique when status=-1 and status=-2

-- ed_code from shop is unique
DROP INDEX order_shop_id_ed_code_idx;
CREATE UNIQUE INDEX order_shop_id_ed_code_idx ON "order" (shop_id, ed_code)
  WHERE partner_id IS NULL AND status != -1 AND fulfillment_shipping_status != -2;

-- external_id from shop is unique
DROP INDEX order_shop_external_id_idx;
CREATE UNIQUE INDEX order_shop_external_id_idx ON "order" (shop_id, external_order_id)
  WHERE external_order_id IS NOT NULL AND partner_id IS NULL AND status != -1 AND fulfillment_shipping_status != -2
    -- workaround for goldship
    AND (shop_id != 1057792338951722956 OR created_at > '2019-2-14 09:00:00 ICT');

-- ed_code from the same partner for each shop is unique
DROP INDEX order_partner_shop_id_external_code_idx;
CREATE UNIQUE INDEX order_partner_shop_id_external_code_idx ON "order" (shop_id, ed_code, partner_id)
  WHERE partner_id IS NOT NULL AND status != -1 AND fulfillment_shipping_status != -2;
