-- Overide 075_shipping_fee.sql
CREATE OR REPLACE FUNCTION fulfillment_update_shipping_fees() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
 item JSON;
BEGIN
  -- do not update completed fulfillments
  IF (NEW.status != 0 AND NEW.status != 2) THEN RETURN NEW; END IF;
  NEW.shipping_fee_shop        = 0;
  NEW.shipping_fee_main        = 0;
  NEW.shipping_fee_return      = 0;
  NEW.shipping_fee_insurance   = 0;
  NEW.shipping_fee_adjustment  = 0;
  NEW.shipping_fee_cods        = 0;
  NEW.shipping_fee_info_change = 0;
  NEW.shipping_fee_other       = 0;
  NEW.shipping_fee_discount    = 0;
  FOR item IN SELECT * FROM jsonb_array_elements_text(NEW.shipping_fee_shop_lines)
  LOOP
    CASE item->>'shipping_fee_type'
      WHEN 'main' THEN
        NEW.shipping_fee_main        = NEW.shipping_fee_main        + (item->>'cost')::INT4;
      WHEN 'return' THEN
        NEW.shipping_fee_return      = NEW.shipping_fee_return      + (item->>'cost')::INT4;
      WHEN 'insurance' THEN
        NEW.shipping_fee_insurance   = NEW.shipping_fee_insurance   + (item->>'cost')::INT4;
      WHEN 'adjustment' THEN
        NEW.shipping_fee_adjustment  = NEW.shipping_fee_adjustment  + (item->>'cost')::INT4;
      WHEN 'cods' THEN
        NEW.shipping_fee_cods        = NEW.shipping_fee_cods        + (item->>'cost')::INT4;
      WHEN 'address_change' THEN
        NEW.shipping_fee_info_change = NEW.shipping_fee_info_change + (item->>'cost')::INT4;
      WHEN 'discount' THEN
        NEW.shipping_fee_discount = NEW.shipping_fee_discount + (item->>'cost')::INT4;
      ELSE
        NEW.shipping_fee_other = NEW.shipping_fee_other + (item->>'cost')::INT4;
    END CASE;
  END LOOP;
  NEW.shipping_fee_shop =
      NEW.shipping_fee_main
    + NEW.shipping_fee_return
    + NEW.shipping_fee_insurance
    + NEW.shipping_fee_adjustment
    + NEW.shipping_fee_cods
    + NEW.shipping_fee_info_change
    + NEW.shipping_fee_other
    + NEW.shipping_fee_discount
    + COALESCE(NEW.etop_fee_adjustment, 0)
    - COALESCE(NEW.etop_discount, 0)
    ;
  RETURN NEW;
END;
$$;
