-- remove old triggers and functions
DROP TRIGGER ffm_update_status ON fulfillment;
DROP TRIGGER order_update_status ON "order";
DROP FUNCTION ffm_update_status();
DROP FUNCTION order_update_status();

-- remove unused columns
ALTER TABLE "order"
  DROP COLUMN IF EXISTS fulfillment_states,
  DROP COLUMN IF EXISTS payment_states;

ALTER TABLE history."order"
  DROP COLUMN IF EXISTS fulfillment_states,
  DROP COLUMN IF EXISTS payment_states;

ALTER TABLE "order"
  RENAME COLUMN fulfillment_status TO fulfillment_shipping_status;

ALTER TABLE history."order"
  RENAME COLUMN fulfillment_status TO fulfillment_shipping_status;

ALTER TABLE "order"
  ADD COLUMN fulfillment_shipping_states  TEXT[],
  ADD COLUMN fulfillment_payment_statuses INT2[],
  ADD COLUMN fulfillment_shipping_codes   TEXT[],
  ADD COLUMN fulfillment_sync_statuses    INT2[];

ALTER TABLE history."order"
  ADD COLUMN fulfillment_shipping_states  TEXT[],
  ADD COLUMN fulfillment_payment_statuses INT2[],
  ADD COLUMN fulfillment_shipping_codes   TEXT[],
  ADD COLUMN fulfillment_sync_statuses    INT2[];

ALTER TABLE fulfillment
  ADD COLUMN etop_payment_status INT2;

ALTER TABLE history.fulfillment
  ADD COLUMN etop_payment_status INT2;

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

-- if all states is default, created then status := 0
-- if all states is cancelled then status := -1
-- if all states is cancelled, returning, returned then status := -2
-- if all states is cancelled, returning, returned and delivered then status := 1
-- otherwise, status := 2
CREATE OR REPLACE FUNCTION coalesce_shipping_states(_states TEXT[]) RETURNS INT2
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

CREATE OR REPLACE FUNCTION coalesce_status4(_statuses INT2[]) RETURNS INT2
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

-- this function populates order shipping_states and payment_statuses from fulfillments
CREATE OR REPLACE FUNCTION order_update_status_from_fulfillment(_order_id INT8) RETURNS void
LANGUAGE plpgsql AS $$
DECLARE
  _fulfillment_shipping_states  TEXT[];
  _fulfillment_payment_statuses INT2[];
  _fulfillment_shipping_codes   TEXT[];
  _fulfillment_sync_statuses    INT2[];
BEGIN
  SELECT
    array_agg(shipping_state),
    array_agg(etop_payment_status),
    array_agg(shipping_code),
    array_agg(sync_status)
  INTO
    _fulfillment_shipping_states,
    _fulfillment_payment_statuses,
    _fulfillment_shipping_codes,
    _fulfillment_sync_statuses
  FROM fulfillment WHERE order_id = _order_id;

  UPDATE "order" SET
    fulfillment_shipping_states  = _fulfillment_shipping_states,
    fulfillment_payment_statuses = _fulfillment_payment_statuses,
    fulfillment_shipping_codes   = _fulfillment_shipping_codes,
    fulfillment_sync_statuses    = _fulfillment_sync_statuses
  WHERE id = _order_id;
END;
$$;

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

  NEW.shipping_status = shipping_state_to_shipping_status(NEW.shipping_state);
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION fulfillment_update_order_status() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  IF (
    NEW.shipping_state      != OLD.shipping_state      OR
    NEW.shipping_status     != OLD.shipping_status     OR
    NEW.sync_status         != OLD.sync_status         OR
    NEW.etop_payment_status != OLD.etop_payment_status OR
    NEW.shipping_code       != OLD.shipping_code
  ) THEN
    EXECUTE order_update_status_from_fulfillment(NEW.order_id);
  END IF;

  RETURN NULL;
END;
$$;

/*
Trigger process:

[fulfillment BEFORE UPDATE] fulfillment_update_status()
[fulfillment AFTER  UPDATE] fulfillment_update_order_status()
                         -> order_update_status_from_fulfilment(order_id)
[order       BEFORE UPDATE] order_update_status()
 */

CREATE TRIGGER fulfillment_update_status
  BEFORE INSERT OR UPDATE ON fulfillment
    FOR EACH ROW EXECUTE PROCEDURE fulfillment_update_status();

CREATE TRIGGER fulfillment_update_order_status
  AFTER UPDATE ON fulfillment
    FOR EACH ROW EXECUTE PROCEDURE fulfillment_update_order_status();

CREATE TRIGGER order_update_status
  BEFORE INSERT OR UPDATE ON "order"
    FOR EACH ROW EXECUTE PROCEDURE order_update_status();

-- update existing orders' status from fulfillments
SELECT id FROM "order", LATERAL order_update_status_from_fulfillment("order".id) s;

-- update existing records, this will trigger above functions and set to correct value
UPDATE fulfillment SET etop_payment_status = 0 WHERE etop_payment_status IS NULL;
UPDATE "order" SET fulfillment_shipping_status = 0 WHERE fulfillment_shipping_status IS NULL;
UPDATE "order" SET etop_payment_status = 0 WHERE etop_payment_status IS NULL;

ALTER TABLE fulfillment ALTER COLUMN etop_payment_status SET NOT NULL;
ALTER TABLE "order"
  ALTER COLUMN fulfillment_shipping_status SET NOT NULL,
  ALTER COLUMN etop_payment_status SET NOT NULL,
  ALTER COLUMN status SET NOT NULL;
