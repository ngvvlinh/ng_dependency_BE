CREATE OR REPLACE FUNCTION fulfillment_update_shipping_code()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  IF (NEW.external_shipping_code IS NOT NULL) THEN
    NEW.shipping_code = NEW.external_shipping_code;
    IF NEW.shipping_provider = 'ghtk' THEN
      NEW.shipping_code = REVERSE(split_part(REVERSE(NEW.shipping_code), '.', 1));
    END IF;
  END IF;
  RETURN NEW;
END
$$;

update fulfillment set shipping_fee_other = 0 where shipping_provider = 'ghtk';
-- Remove this trigger --
DROP TRIGGER fulfillment_update_shipping_code ON fulfillment;
