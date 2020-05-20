ALTER TABLE fulfillment ALTER COLUMN shipping_provider SET NOT NULL;

CREATE OR REPLACE FUNCTION fulfillment_update_shipping_code()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  IF (NEW.external_shipping_code IS NOT NULL) THEN
    NEW.shipping_code = NEW.external_shipping_code;
  END IF;
  RETURN NEW;
END
$$;

CREATE TRIGGER fulfillment_update_shipping_code BEFORE INSERT OR UPDATE ON fulfillment
  FOR EACH ROW EXECUTE PROCEDURE fulfillment_update_shipping_code();
