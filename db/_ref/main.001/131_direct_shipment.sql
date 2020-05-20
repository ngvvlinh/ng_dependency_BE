ALTER TYPE shipping_provider ADD VALUE 'partner';
CREATE UNIQUE INDEX ON "connection"(name);

-- Fulfillment trigger --
CREATE OR REPLACE FUNCTION coalesce_fulfillment_status(
  confirm_status INT2,
  shipping_status INT2,
  sync_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
    IF (confirm_status = -1) THEN
        RETURN -1;
    ELSIF (shipping_status != 0) THEN
        RETURN shipping_status;
    ELSIF confirm_status = 1 AND sync_status = 1 THEN
        RETURN 2;
    ELSE
        RETURN 0;
    END IF;
END;
$$;

-- this function overwrites 120_connection.sql
CREATE OR REPLACE FUNCTION fulfillment_update_status() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.confirm_status := NEW.shop_confirm;
    NEW.shipping_status = shipping_state_to_shipping_status(NEW.shipping_state);

    -- shipping_type 10: shipment
    IF (NEW.shipping_type IS NULL OR NEW.shipping_type = 10) THEN
        IF (NEW.connection_method = 'direct') THEN
            -- trường hợp giao qua NVC nhưng không đối soát với ETOP
            NEW.status = coalesce_fulfillment_status(
                NEW.confirm_status, NEW.shipping_status, NEW.sync_status
            );
        ELSE
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

            NEW.status = coalesce_order_status(
                NEW.confirm_status, NEW.sync_status, NEW.shipping_status, NEW.etop_payment_status);
		END IF;

    ELSE
        -- trường hợp tự giao
        NEW.status = coalesce_fulfillment_status(
            NEW.confirm_status, NEW.shipping_status
        );
    END IF;
    RETURN NEW;
END;
$$;
-- End fulfillment trigger --
