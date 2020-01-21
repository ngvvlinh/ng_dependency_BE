CREATE TABLE connection (
    id INT8 PRIMARY KEY
    , name TEXT
    , status INT2
    , partner_id INT8 REFERENCES partner(id)
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , driver_config JSONB
    , driver TEXT
    , connection_type TEXT
    , connection_subtype TEXT
    , connection_method TEXT
    , connection_provider TEXT
    , etop_affiliate_account JSONB
);

CREATE TABLE shop_connection (
    shop_id INT8 REFERENCES shop(id)
    , connection_id INT8 REFERENCES connection(id)
    , token TEXT
    , token_expires_at TIMESTAMPTZ
    , status INT2
    , connection_states JSONB
    , is_global BOOLEAN
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , external_data JSONB
);

CREATE UNIQUE INDEX ON shop_connection (shop_id, connection_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX ON shop_connection(connection_id) WHERE is_global IS TRUE;

ALTER TABLE fulfillment
    ADD COLUMN shipping_type INT2
    , ADD COLUMN connection_id INT8 REFERENCES connection(id)
    , ADD COLUMN connection_method TEXT
    , ADD COLUMN shop_carrier_id INT8 REFERENCES shop_carrier(id)
    , ADD COLUMN shipping_service_name TEXT
    , ADD COLUMN gross_weight INT
    , ADD COLUMN chargeable_weight INT
    , ADD COLUMN length INT
    , ADD COLUMN width INT
    , ADD COLUMN height INT;

ALTER TABLE history.fulfillment
   ADD COLUMN shipping_type INT2
    , ADD COLUMN connection_id INT8
    , ADD COLUMN connection_method TEXT
    , ADD COLUMN shop_carrier_id INT8
    , ADD COLUMN shipping_service_name TEXT
    , ADD COLUMN gross_weight INT
    , ADD COLUMN chargeable_weight INT
    , ADD COLUMN length INT
    , ADD COLUMN width INT
    , ADD COLUMN height INT;

ALTER TABLE connection
    ADD COLUMN code text
    , ADD COLUMN image_url text;

select init_history('connection', '{id}');
select init_history('shop_connection', '{shop_id,connection_id}');

ALTER TYPE code_type ADD VALUE 'connection';
CREATE UNIQUE INDEX ON connection(code);

INSERT INTO code ("code", "type", "created_at") VALUES
    ('50BS', 'connection', NOW()),
    ('448C', 'connection', NOW()),
    ('E7QC', 'connection', NOW());

INSERT INTO "connection" ("id", "name", "status", "created_at", "updated_at", "driver", "connection_type", "connection_subtype", "connection_method", "connection_provider", "code") VALUES
('1000803215822389663', 'TopShip - GHN', '1', NOW(), NOW(), 'shipping/shipment/topship/ghn', 'shipping', 'shipment', 'topship', 'ghn', '50BS'),
('1000804010396750738', 'TopShip - GHTK', '1', NOW(), NOW(), 'shipping/shipment/topship/ghtk', 'shipping', 'shipment', 'topship', 'ghtk', '448C'),
('1000804104889339180', 'TopShip - VTP', '1', NOW(), NOW(), 'shipping/shipment/topship/vtpost', 'shipping', 'shipment', 'topship', 'vtpost', 'E7QC');

INSERT INTO "shop_connection" ("connection_id", "status", "created_at", "updated_at", "is_global") VALUES
('1000803215822389663','1', NOW(), NOW(), 't'),
('1000804010396750738','1', NOW(), NOW(), 't'),
('1000804104889339180','1', NOW(), NOW(), 't');

-- Fulfillment trigger --
CREATE OR REPLACE FUNCTION coalesce_fulfillment_status(
  confirm_status INT2,
  shipping_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
    IF (confirm_status = -1) THEN
        RETURN -1;
    ELSIF (shipping_status != 0) THEN
        RETURN shipping_status;
    ELSIF confirm_status = 1 THEN
        RETURN 2;
    ELSE
        RETURN 0;
    END IF;
END;
$$;

-- this function overwrites 032_fix_money_transaction.sql
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
                NEW.confirm_status, NEW.shipping_status
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

-- Order trigger --
ALTER TABLE "order"
    ADD COLUMN fulfillment_statuses INT2[];
ALTER TABLE history."order"
    ADD COLUMN fulfillment_statuses INT2[];

CREATE OR REPLACE FUNCTION coalesce_status5(_statuses INT2[])
RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
    IF (_statuses IS NULL OR _statuses <@ '{0}'::INT2[]) THEN
        RETURN  0; -- default, including '{}'
    ELSIF (_statuses <@ '{-1,0}'::INT2[]) THEN
        RETURN -1; -- all N (negative)
    ELSIF (_statuses <@ '{-1,0,-2}') THEN
        RETURN -2;
    ELSIF (_statuses <@ '{-1,0,-2,1}') THEN
        RETURN 1;
    ELSE
        RETURN 2;
    END IF;
END;
$$;

-- this function overwrites 025_status.sql
CREATE OR REPLACE FUNCTION order_update_status_from_fulfillment(_order_id INT8) RETURNS void
LANGUAGE plpgsql AS $$
DECLARE
  _fulfillment_shipping_states  TEXT[];
  _fulfillment_payment_statuses INT2[];
  _fulfillment_shipping_codes   TEXT[];
  _fulfillment_sync_statuses    INT2[];
  _fulfillment_statuses         INT2[];
BEGIN
  SELECT
    array_agg(shipping_state),
    array_agg(etop_payment_status),
    array_agg(shipping_code),
    array_agg(sync_status),
    array_agg(status)
  INTO
    _fulfillment_shipping_states,
    _fulfillment_payment_statuses,
    _fulfillment_shipping_codes,
    _fulfillment_sync_statuses,
    _fulfillment_statuses
  FROM fulfillment WHERE order_id = _order_id;

  UPDATE "order" SET
    fulfillment_shipping_states  = _fulfillment_shipping_states,
    fulfillment_payment_statuses = _fulfillment_payment_statuses,
    fulfillment_shipping_codes   = _fulfillment_shipping_codes,
    fulfillment_sync_statuses    = _fulfillment_sync_statuses,
    fulfillment_statuses         = _fulfillment_statuses
  WHERE id = _order_id;
END;
$$;

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
    RETURN 2;
END;
$$;

CREATE OR REPLACE FUNCTION coalesce_order_status_v2(
  confirm_status INT2,
  payment_status INT2
) RETURNS INT2 IMMUTABLE
LANGUAGE plpgsql AS $$
BEGIN
    IF confirm_status = 1 THEN
        IF payment_status = 1 THEN
            RETURN 1;
        ELSE
            RETURN 2;
        END IF;
    ELSE
        return confirm_status;
    END IF;
END;
$$;

-- this function overwrites 030_shipping_fee.sql
CREATE OR REPLACE FUNCTION order_update_status() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
    fulfillment_status INT2;
BEGIN
    -- update fulfillment and payment status from fulfillments
    NEW.fulfillment_shipping_status = coalesce_shipping_states(NEW.fulfillment_shipping_states);
    NEW.etop_payment_status = coalesce_status4(NEW.fulfillment_payment_statuses);

    -- no longer change order.status if it's either done or cancelled
    IF TG_OP='UPDATE' AND (OLD.status = 1 OR OLD.status = -1) THEN
        RETURN NEW;
    END IF;
    IF NEW.status = 1 OR NEW.status = -1 THEN
        RETURN NEW;
    END IF;

     -- update confirm_status from shop or external, prioritize shop_confirm
    IF NEW.shop_confirm != 0 THEN
		NEW.confirm_status = NEW.shop_confirm;
	ELSE
		NEW.confirm_status = LEAST(NEW.shop_confirm, NEW.customer_confirm, NEW.external_confirm);
	END IF;

    -- Trường hợp không có đơn giao hàng
    IF (NEW.fulfillment_type IS NULL OR NEW.fulfillment_type = 0) AND NEW.fulfillment_statuses IS NULL THEN
        NEW.status = coalesce_order_status_v2(
            NEW.confirm_status, NEW.payment_status
        );
        RETURN NEW;
    END IF;

    -- Trường hợp có đơn giao hàng
    IF NEW.fulfillment_statuses IS NULL THEN
        RETURN NEW;
    END IF;
    fulfillment_status = coalesce_status5(NEW.fulfillment_statuses);
    NEW.status = coalesce_order_status_v2(
        NEW.confirm_status, NEW.payment_status, fulfillment_status
    );

    RETURN NEW;
END;
$$;
-- End order trigger --
