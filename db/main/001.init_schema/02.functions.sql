CREATE FUNCTION public.coalesce_fulfillment_status(confirm_status smallint, shipping_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_fulfillment_status(confirm_status smallint, shipping_status smallint, sync_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_order_status(confirm_status smallint, shipping_status smallint, etop_payment_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_order_status(confirm_status smallint, sync_status smallint, shipping_status smallint, etop_payment_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_order_status_v2(confirm_status smallint, payment_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_order_status_v2(confirm_status smallint, payment_status smallint, fulfillment_status smallint) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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
  RETURN 0;
END;
$$;

CREATE FUNCTION public.coalesce_shipping_states(_states text[]) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_status4(_statuses smallint[]) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.coalesce_status5(_statuses smallint[]) RETURNS smallint
  LANGUAGE plpgsql IMMUTABLE
AS $$
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

CREATE FUNCTION public.ffm_active_supplier(supplier_id bigint, status smallint) RETURNS bigint
  LANGUAGE plpgsql IMMUTABLE
AS $$
BEGIN
  IF (status = -1) THEN RETURN NULL; END IF;
  RETURN supplier_id;
END;
$$;

CREATE FUNCTION public.fulfillment_expected_pick_at(created_at timestamp with time zone) RETURNS timestamp with time zone
  LANGUAGE plpgsql IMMUTABLE
AS $$
DECLARE
  created_hour INT;
begin
  created_hour := date_part('hour', created_at at time zone 'ict');
  IF (created_hour < 10) THEN
    return date_trunc('hour', created_at + (13 - created_hour) * interval '1 hour'); -- pick at 13h
  ELSIF (created_hour < 16) THEN
    return date_trunc('hour', created_at + (21 - created_hour) * interval '1 hour'); -- pick at 21h
  ELSE
    return date_trunc('hour', created_at + (37 - created_hour) * interval '1 hour'); -- pick at 13h tomorrow
  END IF;
END
$$;

CREATE FUNCTION public.fulfillment_update_order_status() RETURNS trigger
  LANGUAGE plpgsql
AS $$
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

CREATE FUNCTION public.fulfillment_update_shipping_fees() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  item JSON;
BEGIN
  -- do not update completed fulfillments
  IF (NEW.status != 0 AND NEW.status != 2) THEN RETURN NEW; END IF;
  NEW.shipping_fee_shop        = 0;
  NEW.shipping_fee_main        = 0;
  NEW.shipping_fee_return      = 0;
  NEW.shipping_fee_insurance   = 0;
  NEW.shipping_fee_adjustment  = 0;
  NEW.shipping_fee_cods        = 0;
  NEW.shipping_fee_info_change = 0;
  NEW.shipping_fee_other       = 0;
  NEW.shipping_fee_discount    = 0;
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
        NEW.shipping_fee_discount = NEW.shipping_fee_discount + (item->>'cost')::INT4;
      ELSE
        NEW.shipping_fee_other = NEW.shipping_fee_other + (item->>'cost')::INT4;
      END CASE;
  END LOOP;
  NEW.shipping_fee_shop =
              NEW.shipping_fee_main
              + NEW.shipping_fee_return
              + NEW.shipping_fee_insurance
              + NEW.shipping_fee_adjustment
              + NEW.shipping_fee_cods
              + NEW.shipping_fee_info_change
              + NEW.shipping_fee_other
              + NEW.shipping_fee_discount
              + COALESCE(NEW.etop_fee_adjustment, 0)
            - COALESCE(NEW.etop_discount, 0)
  ;
  RETURN NEW;
END;
$$;

CREATE FUNCTION public.fulfillment_update_status() RETURNS trigger
  LANGUAGE plpgsql
AS $$
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

CREATE FUNCTION public.order_update_status() RETURNS trigger
  LANGUAGE plpgsql
AS $$
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

CREATE FUNCTION public.order_update_status_from_fulfillment(_order_id bigint) RETURNS void
  LANGUAGE plpgsql
AS $$
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

CREATE FUNCTION public.product_update() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  supplier_status INT2;
  external_status INT2;
  etop_status     INT2;
  price_status    INT2;
BEGIN
  IF TG_OP = 'UPDATE' THEN
    supplier_status := COALESCE(NEW.supplier_status, OLD.supplier_status, 0);
    external_status := COALESCE(NEW.external_status, OLD.external_status, 0);
    etop_status     := COALESCE(NEW.etop_status,     OLD.etop_status,     0);
  ELSE
    supplier_status := COALESCE(NEW.supplier_status, 0);
    external_status := COALESCE(NEW.external_status, 0);
    etop_status     := COALESCE(NEW.etop_status,     0);
  END IF;

  IF NEW.wholesale_price > 0 AND NEW.list_price > 0 THEN
    price_status := 1;
  ELSE
    price_status := 0;
  END IF;
  NEW.status := LEAST(supplier_status, external_status, etop_status, price_status);
  RETURN NEW;
END
$$;

CREATE FUNCTION public.save_history() RETURNS trigger
  LANGUAGE plpgsql
AS $_$
DECLARE
  tbl   TEXT;
  cols  TEXT[];
  rold  hstore;
  rnew  hstore;
  diff  hstore;
  query TEXT;
BEGIN
  tbl  := TG_ARGV[0];
  cols := TG_ARGV[1];

  IF (TG_OP = 'INSERT') THEN
    diff := hstore(NEW);
  ELSEIF (TG_OP = 'DELETE') THEN
    OLD.rid := nextval('history_' || tbl || '_seq');
    diff := hstore(OLD);
  ELSE
    rold := hstore(OLD);
    rnew := hstore(NEW);

    -- ignore empty change
    IF (rnew - rold - '{rid,updated_at}'::TEXT[] = '') THEN
      RETURN NULL;
    END IF;

    diff := rnew - (rold - '{rid}'::TEXT[] - cols);
  END IF;

  SELECT 'INSERT INTO history.' || quote_ident(tbl) || '(' ||
         string_agg(quote_ident(key), ',') ||
         ', _time, _op) VALUES (' ||
         string_agg(quote_nullable(value), ',') ||
         ', now(), $1);'
  INTO query
  FROM each(diff);

  EXECUTE query USING TG_OP::tg_op_type;
  RETURN NULL;
END
$_$;

CREATE FUNCTION public.shipping_state_to_shipping_status(state public.shipping_state) RETURNS smallint
  LANGUAGE plpgsql
AS $$
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

CREATE FUNCTION public.supplier_rules_update() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF (NEW.rules != OLD.rules) THEN
    PERFORM pg_notify('supplier_rules_update', NEW.id::TEXT);
  END IF;
  RETURN NEW;
END;
$$;

CREATE FUNCTION public.update_rid() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  NEW.rid = nextval(TG_ARGV[0]);
  RETURN NEW;
END
$$;

CREATE FUNCTION public.update_to_account() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE account SET
                   name = NEW.name,
                   owner_id = NEW.owner_id,
                   image_url = NEW.image_url,
                   deleted_at = NEW.deleted_at
  WHERE id = NEW.id;
  RETURN NEW;
END;
$$;

CREATE FUNCTION public.variant_update() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  supplier_status INT2;
  external_status INT2;
  etop_status     INT2;
  product_status 	INT2;
BEGIN
  SELECT ve.external_status INTO external_status FROM variant_external as ve WHERE id = NEW.id;
  IF TG_OP = 'UPDATE' THEN
    supplier_status := COALESCE(NEW.ed_status, OLD.ed_status, 0);
    etop_status     := COALESCE(NEW.etop_status, OLD.etop_status, 0);
  ELSE
    supplier_status := COALESCE(NEW.ed_status, 0);
    etop_status     := COALESCE(NEW.etop_status,     0);
  END IF;

  NEW.status := LEAST(supplier_status, external_status, etop_status);

  SELECT MAX(v.status) INTO product_status FROM variant as v WHERE product_id = NEW.product_id AND id != NEW.id;
  product_status := GREATEST(product_status, NEW.status);
  UPDATE product SET status = product_status WHERE id = NEW.product_id;
  RETURN NEW;
END
$$;

CREATE FUNCTION public.init_history(tbl text, cols text[], idx text[] DEFAULT NULL::text[]) RETURNS text
  LANGUAGE plpgsql
AS $_$
DECLARE
  key TEXT[];
BEGIN
  EXECUTE format(
          'ALTER TABLE %1$I ADD COLUMN IF NOT EXISTS rid INT8; '
            'CREATE TABLE history.%1$I AS SELECT * FROM %1$I LIMIT 0; '
            'ALTER TABLE history.%1$I'
            '  ADD UNIQUE (rid),'
            '  ALTER COLUMN rid SET NOT NULL,'
            '  ADD COLUMN _time timestamptz NOT NULL,'
            '  ADD COLUMN _op tg_op_type NOT NULL; '
            'CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON %1$I'
            '  FOR EACH ROW EXECUTE PROCEDURE update_rid(%2$L); '
            'CREATE TRIGGER save_history AFTER INSERT OR UPDATE OR DELETE ON %1$I'
            '  FOR EACH ROW EXECUTE PROCEDURE save_history(%1$L, %3$L); ',
          tbl, 'history_' || tbl || '_seq', cols
    );

  EXECUTE format(
          'CREATE SEQUENCE IF NOT EXISTS history_%1$s_seq START WITH %2$s; ',
          tbl, latest_rid(tbl) + 1
    );

  EXECUTE format(
          'UPDATE %I SET rid = nextval(%L) WHERE rid IS NULL',
          tbl, 'history_' || tbl || '_seq'
    );

  EXECUTE format(
          'ALTER TABLE %1$I ALTER COLUMN rid SET NOT NULL; ',
          tbl
    );

  IF (idx IS NULL) THEN
    idx := cols;
  END IF;
  IF (idx != '{}') THEN
    EXECUTE format(
            'CREATE INDEX ON history.%1$I (%2$s); ',
            tbl, '"' || array_to_string(idx, '","') || '"'
      );
  END IF;
  RETURN tbl;
END
$_$;

CREATE FUNCTION public.latest_rid(_tbl regclass, OUT result bigint) RETURNS bigint
  LANGUAGE plpgsql
AS $$
BEGIN
  EXECUTE format('SELECT rid FROM %s WHERE rid IS NOT NULL ORDER BY rid DESC LIMIT 1::int', _tbl)
    INTO result;
  result := COALESCE(result, 0);
END
$$;

CREATE FUNCTION public.notify_pgrid_id() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':' || NEW.rid
            || ':' || NEW._op
            || ':' || NEW.id
    );
  RETURN NULL;
END$$;

CREATE FUNCTION public.notify_pgrid_shop_customer_group_customer() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'  || NEW.rid
            || ':'  || NEW._op
            || ':c' || NEW.customer_id
            || ':g' || NEW.group_id
    );
  RETURN NULL;
END$$;

CREATE FUNCTION public.notify_pgrid_shop_product() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'   || NEW.rid
            || ':'   || NEW._op
            || ':sh' || NEW.shop_id
            || ':p'  || NEW.product_id
    );
  RETURN NULL;
END$$;

CREATE FUNCTION public.notify_pgrid_shop_product_collection() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'  || NEW.rid
            || ':'  || NEW._op
            || ':p' || NEW.product_id
            || ':c' || NEW.collection_id
    );
  RETURN NULL;
END$$;

CREATE FUNCTION public.notify_pgrid_shop_variant() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'   || NEW.rid
            || ':'   || NEW._op
            || ':sh' || NEW.shop_id
            || ':va' || NEW.variant_id
    );
  RETURN NULL;
END$$;
