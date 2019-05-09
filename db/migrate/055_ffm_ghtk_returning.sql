-- Override 025_status.sql --
CREATE OR REPLACE FUNCTION shipping_state_to_shipping_status(state shipping_state) RETURNS INT2
LANGUAGE plpgsql AS $$
BEGIN
  IF (state = 'default' OR state = 'created') THEN
    RETURN  0;
  ELSIF (state = 'cancelled') THEN
    RETURN -1;
  ELSIF (state = 'returned') THEN
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
