-- this function extends 025_status.sql
--
-- -1:cancelled, 0:default, 1:done, 2:processing
--
--  0: just created, maybe error, can be edited
--  2: processing, can not be edited
--  1: done
-- -1: cancelled or returned
CREATE OR REPLACE FUNCTION coalesce_order_status(
  confirm_status INT2,
  sync_status INT2,
  shipping_status INT2,
  etop_payment_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
  IF (etop_payment_status = 1) THEN
    RETURN shipping_status;
  ELSIF (etop_payment_status = 2 OR shipping_status = 2) THEN
    RETURN 2;
  ELSIF (shipping_status = -1 OR confirm_status = -1) THEN
    RETURN -1;
  ELSIF (confirm_status = 1 AND sync_status = 1) THEN
    RETURN 2;
  ELSE
    RETURN 0;
  END IF;
END;
$$;

-- this function overwrite 025_status.sql
CREATE OR REPLACE FUNCTION coalesce_order_status(
  confirm_status INT2,
  shipping_status INT2,
  etop_payment_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
  IF (etop_payment_status = 1) THEN
    RETURN shipping_status;
  ELSIF (etop_payment_status = 2 OR shipping_status = 2) THEN
    RETURN 2;
  ELSIF (shipping_status = -1 OR confirm_status = -1) THEN
    RETURN -1;
  ELSIF (confirm_status = 1) THEN
    RETURN 2;
  ELSE
    RETURN 0;
  END IF;
END;
$$;

-- this function overwrites 025_status.sql
CREATE OR REPLACE FUNCTION shipping_state_to_shipping_status(state shipping_state)
RETURNS INT2 IMMUTABLE
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

-- this function overwrites 025_status.sql
CREATE OR REPLACE FUNCTION coalesce_shipping_states(_states TEXT[])
RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
  IF (_states IS NULL OR _states <@ '{default,created}'::TEXT[]) THEN
    RETURN  0; -- default, including '{}'
  ELSIF (_states <@ '{cancelled}'::TEXT[]) THEN
    RETURN -1; -- cancelled
  ELSIF (_states <@ '{cancelled,returning,returned}'::TEXT[]) THEN
    RETURN -2; -- return
  ELSIF (_states <@ '{cancelled,returning,returned,delivered}'::TEXT[]) THEN
    RETURN  1; -- delivered
  ELSE
    RETURN  2; -- processing
  END IF;
END;
$$;

-- this function overwrites 025_status.sql
DROP FUNCTION IF EXISTS coalese_status4(INT2[]);

CREATE OR REPLACE FUNCTION coalesce_status4(_statuses INT2[])
RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
  IF (_statuses IS NULL OR _statuses <@ '{0}'::INT2[]) THEN
    RETURN  0; -- default, including '{}'
  ELSIF (_statuses <@ '{-1}'::INT2[]) THEN
    RETURN -1; -- all N (negative)
  ELSIF (_statuses <@ '{-1,0,1}'::INT2[]) THEN
    RETURN  1; -- not include any S
  ELSE
    RETURN  2; -- include any S
  END IF;
END;
$$;

-- this function overwrites 025_status.sql
CREATE OR REPLACE FUNCTION order_update_status() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  -- update fulfillment and payment status from fulfillments
  NEW.fulfillment_shipping_status = coalesce_shipping_states(NEW.fulfillment_shipping_states);
  NEW.etop_payment_status = coalesce_status4(NEW.fulfillment_payment_statuses);

  -- no longer change order.status if it's either done or cancelled
  IF TG_OP='UPDATE' AND (OLD.status = 1 OR OLD.status = -1) THEN	RETURN NEW;	END IF;

  -- update confirm_status from shop or external, prioritize shop_confirm
  IF NEW.shop_confirm != 0 THEN
		NEW.confirm_status = NEW.shop_confirm;
	ELSE
		NEW.confirm_status = LEAST(NEW.shop_confirm, NEW.customer_confirm, NEW.external_confirm);
	END IF;

  NEW.status = coalesce_order_status(
    NEW.confirm_status,
    NEW.fulfillment_shipping_status,
    NEW.etop_payment_status
  );
  RETURN NEW;
END;
$$;

ALTER TABLE fulfillment ADD COLUMN confirm_status INT2;
ALTER TABLE history.fulfillment ADD COLUMN confirm_status INT2;

-- this function overwrites 025_status.sql
CREATE OR REPLACE FUNCTION fulfillment_update_status() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  -- calculate etop_payment_status
  IF (NEW.cod_etop_paid_at IS NOT NULL) THEN
    NEW.etop_payment_status := 1;  -- done
  ELSIF (
    NEW.cod_etop_paid_at IS NULL AND
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

-- update fulfillment.shop_confirm
UPDATE fulfillment SET shop_confirm = o.shop_confirm
FROM "order" AS o WHERE fulfillment.order_id = o.id;

ALTER TABLE fulfillment ALTER COLUMN confirm_status SET NOT NULL;

-- add shipping_fee_* for summary
ALTER TABLE fulfillment
  ADD COLUMN IF NOT EXISTS shipping_fee_main        INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_return      INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_insurance   INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_adjustment  INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_cods        INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_info_change INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_other       INTEGER;

ALTER TABLE history.fulfillment
  ADD COLUMN IF NOT EXISTS shipping_fee_main        INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_return      INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_insurance   INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_adjustment  INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_cods        INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_info_change INTEGER,
  ADD COLUMN IF NOT EXISTS shipping_fee_other       INTEGER;

CREATE OR REPLACE FUNCTION fulfillment_update_shipping_fees() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
 item JSON;
BEGIN
  -- do not update completed fulfillments
  IF (NEW.status != 0 AND NEW.status != 2) THEN RETURN NEW; END IF;

  NEW.shipping_fee_main        = 0;
  NEW.shipping_fee_return      = 0;
  NEW.shipping_fee_insurance   = 0;
  NEW.shipping_fee_adjustment  = 0;
  NEW.shipping_fee_cods        = 0;
  NEW.shipping_fee_info_change = 0;
  NEW.shipping_fee_other       = 0;
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
      ELSE
        NEW.shipping_fee_other = NEW.shipping_fee_other + round((item->>'cost')::INT4);
    END CASE;
  END LOOP;
  RETURN NEW;
END;
$$;

CREATE TRIGGER fulfillment_update_shipping_fees BEFORE INSERT OR UPDATE ON fulfillment
  FOR EACH ROW EXECUTE PROCEDURE fulfillment_update_shipping_fees();

UPDATE fulfillment SET
  shipping_fee_main        = 0, -- trigger the function
  shipping_fee_return      = COALESCE (shipping_fee_return     , 0),
  shipping_fee_insurance   = COALESCE (shipping_fee_insurance  , 0),
  shipping_fee_adjustment  = COALESCE (shipping_fee_adjustment , 0),
  shipping_fee_cods        = COALESCE (shipping_fee_cods       , 0),
  shipping_fee_info_change = COALESCE (shipping_fee_info_change, 0),
  shipping_fee_other       = COALESCE (shipping_fee_other      , 0);

ALTER TABLE fulfillment
  ALTER COLUMN shipping_fee_main        SET NOT NULL,
  ALTER COLUMN shipping_fee_return      SET NOT NULL,
  ALTER COLUMN shipping_fee_insurance   SET NOT NULL,
  ALTER COLUMN shipping_fee_adjustment  SET NOT NULL,
  ALTER COLUMN shipping_fee_cods        SET NOT NULL,
  ALTER COLUMN shipping_fee_info_change SET NOT NULL,
  ALTER COLUMN shipping_fee_other       SET NOT NULL;

ALTER TABLE fulfillment RENAME COLUMN total_amount TO basket_value;
ALTER TABLE history.fulfillment RENAME COLUMN total_amount TO basket_value;

ALTER TABLE fulfillment
  ALTER COLUMN total_items SET NOT NULL,
  ALTER COLUMN total_weight SET NOT NULL,
  ALTER COLUMN shipping_fee_customer SET NOT NULL,
  ALTER COLUMN shipping_fee_shop SET NOT NULL,
  ALTER COLUMN external_shipping_fee SET NOT NULL,
  ALTER COLUMN basket_value SET NOT NULL,
  ALTER COLUMN total_cod_amount SET NOT NULL;

