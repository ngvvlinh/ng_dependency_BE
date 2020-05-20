-- this function overwrites 030_shipping_fee.sql
CREATE OR REPLACE FUNCTION coalesce_shipping_states(_states TEXT[])
RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
  IF (_states IS NULL OR _states <@ '{default,created}'::TEXT[]) THEN
    RETURN  0; -- default, including '{}'
  ELSIF (_states <@ '{cancelled}'::TEXT[]) THEN
    RETURN -1; -- cancelled
  ELSIF (_states <@ '{cancelled,returning,returned,undeliverable}'::TEXT[]) THEN
    RETURN -2; -- return
  ELSIF (_states <@ '{cancelled,returning,returned,delivered}'::TEXT[]) THEN
    RETURN  1; -- delivered
  ELSE
    RETURN  2; -- processing
  END IF;
END;
$$;

-- Override 057_ffm_returning_revert.sql --
CREATE OR REPLACE FUNCTION shipping_state_to_shipping_status(state shipping_state) RETURNS INT2
LANGUAGE plpgsql AS $$
BEGIN
  IF (state = 'default' OR state = 'created') THEN
    RETURN  0;
  ELSIF (state = 'cancelled') THEN
    RETURN -1;
  ELSIF (state = 'returning' OR state = 'returned' OR state = 'undeliverable') THEN
    RETURN -2;
  ELSIF (state = 'delivered') THEN
    RETURN  1;
  ELSE
    RETURN  2;
  END IF;
END;
$$;

ALTER TABLE fulfillment ADD COLUMN actual_compensation_amount integer;
ALTER TABLE history.fulfillment ADD COLUMN actual_compensation_amount integer;
