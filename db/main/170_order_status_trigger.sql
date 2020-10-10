-- this function overwrites 120_connection.sql
CREATE OR REPLACE FUNCTION coalesce_order_status_v2(
    confirm_status INT2,
    payment_status INT2,
    fulfillment_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
    IF fulfillment_status = 2 THEN
        RETURN 2;
    END IF;
    IF fulfillment_status = 1 OR fulfillment_status = -2 THEN
        IF payment_status = 1 THEN
            RETURN fulfillment_status;
        ELSE
            RETURN 2;
        END IF;
    END IF;
    IF fulfillment_status = -1 THEN
        IF confirm_status = -1 THEN
            RETURN -1;
        ELSE
            RETURN 2;
        END IF;
    END IF;
    RETURN coalesce_order_status_v2(confirm_status, payment_status);
END;
$$;
