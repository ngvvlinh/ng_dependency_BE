ALTER TABLE money_transaction ALTER COLUMN provider DROP NOT NULL;

ALTER TABLE fulfillment ADD COLUMN money_transaction_shipping_external_id BIGINT REFERENCES money_transaction_shipping_external(id);
ALTER TABLE history.fulfillment ADD COLUMN money_transaction_shipping_external_id BIGINT;

ALTER TABLE money_transaction RENAME TO money_transaction_shipping;
ALTER TABLE history.money_transaction RENAME TO money_transaction_shipping;
ALTER TYPE code_type ADD VALUE 'money_transaction_shipping';

ALTER TABLE money_transaction_shipping RENAME COLUMN etop_paid_at TO etop_transfered_at;
ALTER TABLE money_transaction_shipping ADD COLUMN confirmed_at timestamp with time zone;
ALTER TABLE history.money_transaction_shipping ADD COLUMN confirmed_at timestamp with time zone;
ALTER TABLE history.money_transaction_shipping RENAME COLUMN etop_paid_at TO etop_transfered_at;

DROP TRIGGER save_history ON money_transaction_shipping;
CREATE TRIGGER save_history AFTER INSERT OR UPDATE OR DELETE ON money_transaction_shipping FOR EACH ROW EXECUTE PROCEDURE save_history('money_transaction_shipping', '{id}');

ALTER TABLE fulfillment RENAME COLUMN cod_etop_paid_at TO cod_etop_transfered_at;
ALTER TABLE fulfillment RENAME COLUMN shipping_fee_shop_paid_at TO shipping_fee_shop_transfered_at;
ALTER TABLE history.fulfillment RENAME COLUMN cod_etop_paid_at TO cod_etop_transfered_at;
ALTER TABLE history.fulfillment RENAME COLUMN shipping_fee_shop_paid_at TO shipping_fee_shop_transfered_at;

-- this function overwrites 030_shipping_fee.sql
CREATE OR REPLACE FUNCTION fulfillment_update_status() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  -- calculate etop_payment_status
  IF (NEW.cod_etop_transfered_at IS NOT NULL) THEN
    NEW.etop_payment_status := 1;  -- done
  ELSIF (
    NEW.cod_etop_transfered_at IS NULL AND
    NOT (NEW.shipping_state IN ('default','created','picking','cancelled'))
  ) THEN
    NEW.etop_payment_status := 2;  -- processing
  ELSE
    NEW.etop_payment_status := 0;
  END IF;

  IF (NEW.shipping_provider = 'manual') THEN
    RETURN NEW;
  END IF;

  IF (NEW.supplier_id IS NULL) THEN
    NEW.confirm_status := NEW.shop_confirm;
  ELSE
    NEW.confirm_status := LEAST(NEW.shop_confirm, NEW.supplier_confirm);
  END IF;

  NEW.shipping_status = shipping_state_to_shipping_status(NEW.shipping_state);
  NEW.status = coalesce_order_status(
    NEW.confirm_status, NEW.sync_status, NEW.shipping_status, NEW.etop_payment_status
  );
  RETURN NEW;
END;
$$;
