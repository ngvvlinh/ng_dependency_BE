-- this function overwrite 025_status.sql and 030_shipping_fee.sql
-- refer to doc/fulfillment_status.xlsx for more details
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

-- this function overwrite 025_status.sql and 030_shipping_fee.sql
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
