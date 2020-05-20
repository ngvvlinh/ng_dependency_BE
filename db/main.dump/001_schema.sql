--
-- Name: history; Type: SCHEMA; Schema: -; Owner: etop
--

CREATE SCHEMA history;


ALTER SCHEMA history OWNER TO etop;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: hstore; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;


--
-- Name: EXTENSION hstore; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION hstore IS 'data type for storing sets of (key, value) pairs';


--
-- Name: account_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.account_type AS ENUM (
  'etop',
  'supplier',
  'shop',
  'partner',
  'affiliate'
  );


ALTER TYPE public.account_type OWNER TO etop;

--
-- Name: address_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.address_type AS ENUM (
  'billing',
  'shipping',
  'general',
  'warehouse',
  'shipfrom',
  'shipto'
  );


ALTER TYPE public.address_type OWNER TO etop;

--
-- Name: code_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.code_type AS ENUM (
  'order',
  'money_transaction',
  'shop',
  'money_transaction_external',
  'money_transaction_shipping',
  'money_transaction_shipping_etop',
  'connection'
  );


ALTER TYPE public.code_type OWNER TO etop;

--
-- Name: contact_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.contact_type AS ENUM (
  'phone',
  'email'
  );


ALTER TYPE public.contact_type OWNER TO etop;

--
-- Name: customer_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.customer_type AS ENUM (
  'individual',
  'organization',
  'independent',
  'anonymous'
  );


ALTER TYPE public.customer_type OWNER TO etop;

--
-- Name: fulfillment_endpoint; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.fulfillment_endpoint AS ENUM (
  'supplier',
  'shop',
  'customer'
  );


ALTER TYPE public.fulfillment_endpoint OWNER TO etop;

--
-- Name: gender_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.gender_type AS ENUM (
  'male',
  'female',
  'other'
  );


ALTER TYPE public.gender_type OWNER TO etop;

--
-- Name: ghn_note_code; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.ghn_note_code AS ENUM (
  'CHOTHUHANG',
  'CHOXEMHANGKHONGTHU',
  'KHONGCHOXEMHANG'
  );


ALTER TYPE public.ghn_note_code OWNER TO etop;

--
-- Name: import_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.import_type AS ENUM (
  'other',
  'shop_order',
  'shop_product',
  'etop_money_transaction'
  );


ALTER TYPE public.import_type OWNER TO etop;

--
-- Name: inventory_voucher_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.inventory_voucher_type AS ENUM (
  'in',
  'out'
  );


ALTER TYPE public.inventory_voucher_type OWNER TO etop;

--
-- Name: order_source_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.order_source_type AS ENUM (
  'unknown',
  'self',
  'import',
  'api',
  'etop_pos',
  'etop_pxs',
  'etop_cmx',
  'ts_app',
  'etop_app',
  'haravan'
  );


ALTER TYPE public.order_source_type OWNER TO etop;

--
-- Name: partial_status; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.partial_status AS ENUM (
  'default',
  'partial',
  'done',
  'cancelled'
  );


ALTER TYPE public.partial_status OWNER TO etop;

--
-- Name: payment_method_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.payment_method_type AS ENUM (
  'other',
  'bank',
  'cod'
  );


ALTER TYPE public.payment_method_type OWNER TO etop;

--
-- Name: processing_status; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.processing_status AS ENUM (
  'default',
  'processing',
  'done',
  'cancelled'
  );


ALTER TYPE public.processing_status OWNER TO etop;

--
-- Name: product_source_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.product_source_type AS ENUM (
  'custom',
  'kiotviet'
  );


ALTER TYPE public.product_source_type OWNER TO etop;

--
-- Name: province_region; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.province_region AS ENUM (
  'north',
  'south',
  'middle'
  );


ALTER TYPE public.province_region OWNER TO etop;

--
-- Name: receipt_created_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.receipt_created_type AS ENUM (
  'manual',
  'auto'
  );


ALTER TYPE public.receipt_created_type OWNER TO etop;

--
-- Name: receipt_ref_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.receipt_ref_type AS ENUM (
  'order',
  'fulfillment',
  'inventory_voucher',
  'purchase_order',
  'refund',
  'purchase_refund'
  );


ALTER TYPE public.receipt_ref_type OWNER TO etop;

--
-- Name: receipt_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.receipt_type AS ENUM (
  'receipt',
  'payment'
  );


ALTER TYPE public.receipt_type OWNER TO etop;

--
-- Name: shipping_provider; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.shipping_provider AS ENUM (
  'ghn',
  'manual',
  'ghtk',
  'vtpost',
  'partner'
  );


ALTER TYPE public.shipping_provider OWNER TO etop;

--
-- Name: shipping_state; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.shipping_state AS ENUM (
  'default',
  'unknown',
  'created',
  'confirmed',
  'picking',
  'processing',
  'holding',
  'returning',
  'returned',
  'delivering',
  'delivered',
  'undeliverable',
  'cancelled',
  'closed'
  );


ALTER TYPE public.shipping_state OWNER TO etop;

--
-- Name: shop_ledger_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.shop_ledger_type AS ENUM (
  'cash',
  'bank'
  );


ALTER TYPE public.shop_ledger_type OWNER TO etop;

--
-- Name: subject_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.subject_type AS ENUM (
  'account',
  'user'
  );


ALTER TYPE public.subject_type OWNER TO etop;

--
-- Name: tg_op_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.tg_op_type AS ENUM (
  'INSERT',
  'UPDATE',
  'DELETE'
  );


ALTER TYPE public.tg_op_type OWNER TO etop;

--
-- Name: trader_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.trader_type AS ENUM (
  'customer',
  'vendor',
  'supplier',
  'carrier'
  );


ALTER TYPE public.trader_type OWNER TO etop;

--
-- Name: try_on; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.try_on AS ENUM (
  'none',
  'open',
  'try'
  );


ALTER TYPE public.try_on OWNER TO etop;

--
-- Name: user_identifying_type; Type: TYPE; Schema: public; Owner: etop
--

CREATE TYPE public.user_identifying_type AS ENUM (
  'full',
  'half',
  'stub'
  );


ALTER TYPE public.user_identifying_type OWNER TO etop;

--
-- Name: clean_ghn_web_data(text); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.clean_ghn_web_data(input text) RETURNS text
  LANGUAGE plpgsql
AS $$
BEGIN
  RETURN REPLACE(
          REGEXP_REPLACE(
                  input,
                  '(<[^>]+>)|([\[\]{}"]+)|(: )|(\.000Z)|StateText|Message',
                  '', 'g'
            ),
          'Time',
          ' 📦 '
    );
END
$$;


ALTER FUNCTION public.clean_ghn_web_data(input text) OWNER TO etop;

--
-- Name: coalesce_fulfillment_status(smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_fulfillment_status(confirm_status smallint, shipping_status smallint) OWNER TO etop;

--
-- Name: coalesce_fulfillment_status(smallint, smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_fulfillment_status(confirm_status smallint, shipping_status smallint, sync_status smallint) OWNER TO etop;

--
-- Name: coalesce_order_status(smallint, smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_order_status(confirm_status smallint, shipping_status smallint, etop_payment_status smallint) OWNER TO etop;

--
-- Name: coalesce_order_status(smallint, smallint, smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_order_status(confirm_status smallint, sync_status smallint, shipping_status smallint, etop_payment_status smallint) OWNER TO etop;

--
-- Name: coalesce_order_status_v2(smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_order_status_v2(confirm_status smallint, payment_status smallint) OWNER TO etop;

--
-- Name: coalesce_order_status_v2(smallint, smallint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_order_status_v2(confirm_status smallint, payment_status smallint, fulfillment_status smallint) OWNER TO etop;

--
-- Name: coalesce_shipping_states(text[]); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_shipping_states(_states text[]) OWNER TO etop;

--
-- Name: coalesce_status4(smallint[]); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_status4(_statuses smallint[]) OWNER TO etop;

--
-- Name: coalesce_status5(smallint[]); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.coalesce_status5(_statuses smallint[]) OWNER TO etop;

--
-- Name: convertinterval2hoursminutes(interval); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.convertinterval2hoursminutes(t interval) RETURNS text
  LANGUAGE plpgsql
AS $$
BEGIN
  return concat(EXTRACT(epoch FROM t)::int/3600 , ':', (EXTRACT(epoch FROM t)::int - EXTRACT(epoch FROM t)::int/3600*3600)/60);
END;
$$;


ALTER FUNCTION public.convertinterval2hoursminutes(t interval) OWNER TO etop;

--
-- Name: convertinterval2hoursminutes(timestamp without time zone, timestamp without time zone); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.convertinterval2hoursminutes(t1 timestamp without time zone, t2 timestamp without time zone) RETURNS text
  LANGUAGE plpgsql
AS $$
DECLARE t INTERVAL;
BEGIN
  t = t1::INTERVAL - t2::INTERVAL;
  return concat(EXTRACT(epoch FROM t::INTERVAL)::int/3600 , ':', (EXTRACT(epoch FROM t::INTERVAL)::int - EXTRACT(epoch FROM t::INTERVAL)::int/3600*3600)/60);
END;
$$;


ALTER FUNCTION public.convertinterval2hoursminutes(t1 timestamp without time zone, t2 timestamp without time zone) OWNER TO etop;

--
-- Name: convertinterval2hoursminutes(timestamp with time zone, timestamp with time zone); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.convertinterval2hoursminutes(t1 timestamp with time zone, t2 timestamp with time zone) RETURNS text
  LANGUAGE plpgsql
AS $$
DECLARE t INTERVAL;
BEGIN
  t = t1 - t2;
  return concat(EXTRACT(epoch FROM t::INTERVAL)::int/3600 , ':', (EXTRACT(epoch FROM t::INTERVAL)::int - EXTRACT(epoch FROM t::INTERVAL)::int/3600*3600)/60);
END;
$$;


ALTER FUNCTION public.convertinterval2hoursminutes(t1 timestamp with time zone, t2 timestamp with time zone) OWNER TO etop;

--
-- Name: ffm_active_supplier(bigint, smallint); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.ffm_active_supplier(supplier_id bigint, status smallint) RETURNS bigint
  LANGUAGE plpgsql IMMUTABLE
AS $$
BEGIN
  IF (status = -1) THEN RETURN NULL; END IF;
  RETURN supplier_id;
END;
$$;


ALTER FUNCTION public.ffm_active_supplier(supplier_id bigint, status smallint) OWNER TO etop;

--
-- Name: fulfillment_expected_pick_at(timestamp with time zone); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.fulfillment_expected_pick_at(created_at timestamp with time zone) OWNER TO etop;

--
-- Name: fulfillment_update_order_status(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.fulfillment_update_order_status() OWNER TO etop;

--
-- Name: fulfillment_update_shipping_fees(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.fulfillment_update_shipping_fees() OWNER TO etop;

--
-- Name: fulfillment_update_status(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.fulfillment_update_status() OWNER TO etop;

--
-- Name: ids_not_empty(bigint[]); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.ids_not_empty(ids bigint[]) RETURNS boolean
  LANGUAGE plpgsql IMMUTABLE
AS $$
BEGIN
  RETURN ((array_length(ids, 1) is null) != (0 != ALL(ids)));
END;
$$;


ALTER FUNCTION public.ids_not_empty(ids bigint[]) OWNER TO etop;

--
-- Name: init_history(text, text[], text[]); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.init_history(tbl text, cols text[], idx text[]) OWNER TO etop;

--
-- Name: latest_rid(regclass); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.latest_rid(_tbl regclass, OUT result bigint) RETURNS bigint
  LANGUAGE plpgsql
AS $$
BEGIN
  EXECUTE format('SELECT rid FROM %s WHERE rid IS NOT NULL ORDER BY rid DESC LIMIT 1::int', _tbl)
    INTO result;
  result := COALESCE(result, 0);
END
$$;


ALTER FUNCTION public.latest_rid(_tbl regclass, OUT result bigint) OWNER TO etop;

--
-- Name: notify_pgrid_id(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.notify_pgrid_id() OWNER TO etop;

--
-- Name: notify_pgrid_shop_customer_group_customer(); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.notify_pgrid_shop_customer_group_customer() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'   || NEW.rid
            || ':'   || NEW._op
            || ':c' || NEW.customer_id
            || ':g'  || NEW.group_id
    );
  RETURN NULL;
END$$;


ALTER FUNCTION public.notify_pgrid_shop_customer_group_customer() OWNER TO etop;

--
-- Name: notify_pgrid_shop_product(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.notify_pgrid_shop_product() OWNER TO etop;

--
-- Name: notify_pgrid_shop_product_collection(); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.notify_pgrid_shop_product_collection() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM pg_notify(
          'pgrid',
          TG_TABLE_NAME
            || ':'   || NEW.rid
            || ':'   || NEW._op
            || ':p' || NEW.product_id
            || ':c'  || NEW.collection_id
    );
  RETURN NULL;
END$$;


ALTER FUNCTION public.notify_pgrid_shop_product_collection() OWNER TO etop;

--
-- Name: notify_pgrid_shop_variant(); Type: FUNCTION; Schema: public; Owner: etop
--

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
            || ':va'  || NEW.variant_id
    );
  RETURN NULL;
END$$;


ALTER FUNCTION public.notify_pgrid_shop_variant() OWNER TO etop;

--
-- Name: order_update_status(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.order_update_status() OWNER TO etop;

--
-- Name: order_update_status_from_fulfillment(bigint); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.order_update_status_from_fulfillment(_order_id bigint) OWNER TO etop;

--
-- Name: product_update(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.product_update() OWNER TO etop;

--
-- Name: save_history(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.save_history() OWNER TO etop;

--
-- Name: shipping_state_to_shipping_status(public.shipping_state); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.shipping_state_to_shipping_status(state public.shipping_state) OWNER TO etop;

--
-- Name: supplier_rules_update(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.supplier_rules_update() OWNER TO etop;

--
-- Name: update_rid(); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.update_rid() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  NEW.rid = nextval(TG_ARGV[0]);
  RETURN NEW;
END
$$;


ALTER FUNCTION public.update_rid() OWNER TO etop;

--
-- Name: update_to_account(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.update_to_account() OWNER TO etop;

--
-- Name: variant_external_update(); Type: FUNCTION; Schema: public; Owner: etop
--

CREATE FUNCTION public.variant_external_update() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  supplier_status INT2;
  x_status INT2;
  etop_status     INT2;
  final_status    INT2;
BEGIN
  SELECT v.ed_status, v.etop_status INTO supplier_status, etop_status FROM variant as v WHERE id = NEW.id;
  IF TG_OP = 'UPDATE' THEN
    x_status := COALESCE(NEW.external_status, OLD.external_status, 0);
  ELSE
    x_status := COALESCE(NEW.external_status, 0);
  END IF;

  final_status := LEAST(supplier_status, x_status, etop_status);
  UPDATE variant SET status = final_status where id = NEW.id;
  RETURN NEW;
END
$$;


ALTER FUNCTION public.variant_external_update() OWNER TO etop;

--
-- Name: variant_update(); Type: FUNCTION; Schema: public; Owner: etop
--

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


ALTER FUNCTION public.variant_update() OWNER TO etop;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.account (
                               id bigint,
                               name text,
                               deleted_at timestamp with time zone,
                               image_url text,
                               type public.account_type,
                               rid bigint NOT NULL,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               url_slug text,
                               owner_id bigint
);


ALTER TABLE history.account OWNER TO etop;

--
-- Name: account_auth; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.account_auth (
                                    auth_key text,
                                    account_id bigint,
                                    status smallint,
                                    roles text[],
                                    permissions text[],
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    deleted_at timestamp with time zone,
                                    rid bigint NOT NULL,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL
);


ALTER TABLE history.account_auth OWNER TO etop;

--
-- Name: account_user; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.account_user (
                                    account_id bigint,
                                    user_id bigint,
                                    status smallint,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    roles text[],
                                    permissions text[],
                                    deleted_at text,
                                    rid bigint NOT NULL,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL,
                                    full_name text,
                                    short_name text,
                                    "position" text,
                                    response_status smallint,
                                    invitation_sent_at timestamp with time zone,
                                    invitation_sent_by bigint,
                                    invitation_accepted_at timestamp with time zone,
                                    invitation_rejected_at timestamp with time zone,
                                    disabled_at timestamp with time zone,
                                    disabled_by bigint,
                                    disable_reason text
);


ALTER TABLE history.account_user OWNER TO etop;

--
-- Name: address; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.address (
                               id bigint,
                               country text,
                               province_code text,
                               province text,
                               district_code text,
                               district text,
                               ward text,
                               ward_code text,
                               address1 text,
                               is_default boolean,
                               type public.address_type,
                               account_id bigint,
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               full_name text,
                               first_name text,
                               last_name text,
                               email text,
                               "position" text,
                               city text,
                               zip text,
                               address2 text,
                               phone text,
                               company text,
                               notes jsonb,
                               rid bigint NOT NULL,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               coordinates jsonb
);


ALTER TABLE history.address OWNER TO etop;

--
-- Name: connection; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.connection (
                                  id bigint,
                                  name text,
                                  status smallint,
                                  partner_id bigint,
                                  created_at timestamp with time zone,
                                  updated_at timestamp with time zone,
                                  deleted_at timestamp with time zone,
                                  driver_config jsonb,
                                  driver text,
                                  connection_type text,
                                  connection_subtype text,
                                  connection_method text,
                                  connection_provider text,
                                  etop_affiliate_account jsonb,
                                  code text,
                                  image_url text,
                                  rid bigint NOT NULL,
                                  _time timestamp with time zone NOT NULL,
                                  _op public.tg_op_type NOT NULL
);


ALTER TABLE history.connection OWNER TO etop;

--
-- Name: credit; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.credit (
                              id bigint,
                              rid bigint NOT NULL,
                              amount integer,
                              shop_id bigint,
                              supplier_id bigint,
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              paid_at timestamp with time zone,
                              type public.account_type,
                              status smallint,
                              _time timestamp with time zone NOT NULL,
                              _op public.tg_op_type NOT NULL
);


ALTER TABLE history.credit OWNER TO etop;

--
-- Name: etop_category; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.etop_category (
                                     id bigint,
                                     name text,
                                     status smallint,
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     parent_id bigint,
                                     rid bigint NOT NULL,
                                     _time timestamp with time zone NOT NULL,
                                     _op public.tg_op_type NOT NULL
);


ALTER TABLE history.etop_category OWNER TO etop;

--
-- Name: fulfillment; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.fulfillment (
                                   id bigint,
                                   order_id bigint,
                                   lines jsonb,
                                   variant_ids bigint[],
                                   type_from public.fulfillment_endpoint,
                                   type_to public.fulfillment_endpoint,
                                   address_from jsonb,
                                   address_to jsonb,
                                   shop_id bigint,
                                   supplier_id bigint,
                                   total_items integer,
                                   total_weight integer,
                                   basket_value integer,
                                   total_cod_amount integer,
                                   shipping_fee_customer integer,
                                   shipping_fee_shop integer,
                                   external_shipping_fee integer,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   shipping_cancelled_at timestamp with time zone,
                                   cancel_reason text,
                                   closed_at timestamp with time zone,
                                   shipping_provider public.shipping_provider,
                                   shipping_code text,
                                   shipping_note text,
                                   external_shipping_id text,
                                   external_shipping_code text,
                                   external_shipping_created_at timestamp with time zone,
                                   external_shipping_updated_at timestamp with time zone,
                                   external_shipping_cancelled_at timestamp with time zone,
                                   external_shipping_delivered_at timestamp with time zone,
                                   external_shipping_returned_at timestamp with time zone,
                                   external_shipping_state text,
                                   external_shipping_status smallint,
                                   external_shipping_data jsonb,
                                   shipping_state public.shipping_state,
                                   status smallint,
                                   sync_status smallint,
                                   rid bigint NOT NULL,
                                   sync_states jsonb,
                                   shipping_delivered_at timestamp with time zone,
                                   shipping_returned_at timestamp with time zone,
                                   external_shipping_closed_at timestamp with time zone,
                                   supplier_confirm smallint,
                                   shop_confirm smallint,
                                   last_sync_at timestamp with time zone,
                                   expected_delivery_at timestamp with time zone,
                                   money_transaction_id bigint,
                                   cod_etop_transfered_at timestamp with time zone,
                                   shipping_fee_shop_transfered_at timestamp with time zone,
                                   _time timestamp with time zone NOT NULL,
                                   _op public.tg_op_type NOT NULL,
                                   provider_shipping_fee_lines jsonb,
                                   shipping_fee_shop_lines jsonb,
                                   etop_discount integer,
                                   shipping_status smallint,
                                   etop_fee_adjustment integer,
                                   etop_payment_status smallint,
                                   address_to_province_code text,
                                   address_to_district_code text,
                                   address_to_ward_code text,
                                   provider_service_id text,
                                   expected_pick_at timestamp with time zone,
                                   confirm_status smallint,
                                   shipping_fee_main integer,
                                   shipping_fee_return integer,
                                   shipping_fee_insurance integer,
                                   shipping_fee_adjustment integer,
                                   shipping_fee_cods integer,
                                   shipping_fee_info_change integer,
                                   shipping_fee_other integer,
                                   external_shipping_state_code text,
                                   money_transaction_shipping_external_id bigint,
                                   total_discount integer,
                                   total_amount integer,
                                   external_shipping_logs jsonb,
                                   partner_id bigint,
                                   external_shipping_sub_state text,
                                   external_shipping_note text,
                                   try_on public.try_on,
                                   external_shipping_name text,
                                   shipping_service_fee integer,
                                   original_cod_amount integer,
                                   address_return jsonb,
                                   include_insurance boolean,
                                   admin_note text,
                                   is_partial_delivery boolean,
                                   shipping_fee_discount integer,
                                   shipping_created_at timestamp with time zone,
                                   shipping_holding_at timestamp with time zone,
                                   shipping_picking_at timestamp with time zone,
                                   shipping_delivering_at timestamp with time zone,
                                   shipping_returning_at timestamp with time zone,
                                   etop_adjusted_shipping_fee_main integer,
                                   etop_price_rule boolean,
                                   actual_compensation_amount integer,
                                   delivery_route text,
                                   created_by bigint,
                                   shipping_type smallint,
                                   connection_id bigint,
                                   connection_method text,
                                   shop_carrier_id bigint,
                                   shipping_service_name text,
                                   gross_weight integer,
                                   chargeable_weight integer,
                                   length integer,
                                   width integer,
                                   height integer,
                                   external_affiliate_id text
);


ALTER TABLE history.fulfillment OWNER TO etop;

--
-- Name: import_attempt; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.import_attempt (
                                      id bigint,
                                      user_id bigint,
                                      account_id bigint,
                                      original_file text,
                                      stored_file text,
                                      type public.import_type,
                                      n_created smallint,
                                      n_updated smallint,
                                      n_error smallint,
                                      status smallint,
                                      error_type text,
                                      errors jsonb,
                                      duration_ms integer,
                                      created_at timestamp with time zone,
                                      rid bigint NOT NULL,
                                      _time timestamp with time zone NOT NULL,
                                      _op public.tg_op_type NOT NULL
);


ALTER TABLE history.import_attempt OWNER TO etop;

--
-- Name: inventory_variant; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.inventory_variant (
                                         shop_id bigint,
                                         variant_id bigint,
                                         quantity_on_hand integer,
                                         quantity_picked integer,
                                         cost_price integer,
                                         created_at timestamp with time zone,
                                         updated_at timestamp with time zone,
                                         rid bigint NOT NULL,
                                         _time timestamp with time zone NOT NULL,
                                         _op public.tg_op_type NOT NULL
);


ALTER TABLE history.inventory_variant OWNER TO etop;

--
-- Name: inventory_voucher; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.inventory_voucher (
                                         title character varying,
                                         shop_id bigint,
                                         id bigint,
                                         created_by bigint,
                                         updated_by bigint,
                                         status integer,
                                         trader_id bigint,
                                         total_amount integer,
                                         type public.inventory_voucher_type,
                                         created_at timestamp with time zone,
                                         updated_at timestamp with time zone,
                                         confirmed_at timestamp with time zone,
                                         cancelled_at timestamp with time zone,
                                         cancel_reason character varying,
                                         lines jsonb,
                                         ref_id bigint,
                                         ref_type text,
                                         code_norm integer,
                                         code text,
                                         trader jsonb,
                                         variant_ids bigint[],
                                         ref_code text,
                                         product_ids bigint[],
                                         rid bigint NOT NULL,
                                         _time timestamp with time zone NOT NULL,
                                         _op public.tg_op_type NOT NULL
);


ALTER TABLE history.inventory_voucher OWNER TO etop;

--
-- Name: invitation; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.invitation (
                                  id bigint,
                                  account_id bigint,
                                  email text,
                                  roles text[],
                                  token text,
                                  status smallint,
                                  invited_by bigint,
                                  accepted_at timestamp with time zone,
                                  rejected_at timestamp with time zone,
                                  expires_at timestamp with time zone,
                                  created_at timestamp with time zone,
                                  updated_at timestamp with time zone,
                                  deleted_at timestamp with time zone,
                                  rid bigint NOT NULL,
                                  _time timestamp with time zone NOT NULL,
                                  _op public.tg_op_type NOT NULL,
                                  full_name text,
                                  short_name text,
                                  "position" text,
                                  phone text
);


ALTER TABLE history.invitation OWNER TO etop;

--
-- Name: money_transaction_shipping; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.money_transaction_shipping (
                                                  id bigint,
                                                  shop_id bigint,
                                                  supplier_id bigint,
                                                  status smallint,
                                                  total_cod integer,
                                                  total_orders integer,
                                                  code text,
                                                  created_at timestamp with time zone,
                                                  updated_at timestamp with time zone,
                                                  closed_at timestamp with time zone,
                                                  money_transaction_shipping_external_id bigint,
                                                  provider public.shipping_provider,
                                                  etop_transfered_at timestamp with time zone,
                                                  total_amount integer,
                                                  rid bigint NOT NULL,
                                                  _time timestamp with time zone NOT NULL,
                                                  _op public.tg_op_type NOT NULL,
                                                  confirmed_at timestamp with time zone,
                                                  money_transaction_shipping_etop_id bigint,
                                                  bank_account json,
                                                  note text,
                                                  invoice_number text,
                                                  type text
);


ALTER TABLE history.money_transaction_shipping OWNER TO etop;

--
-- Name: order; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history."order" (
                               id bigint,
                               rid bigint NOT NULL,
                               shop_id bigint,
                               code text,
                               product_ids bigint[],
                               variant_ids bigint[],
                               supplier_ids bigint[],
                               currency text,
                               payment_method text,
                               customer jsonb,
                               customer_address jsonb,
                               billing_address jsonb,
                               shipping_address jsonb,
                               customer_phone text,
                               customer_email text,
                               created_at timestamp with time zone,
                               processed_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               closed_at timestamp with time zone,
                               confirmed_at timestamp with time zone,
                               cancelled_at timestamp with time zone,
                               cancel_reason text,
                               customer_confirm smallint,
                               external_confirm smallint,
                               shop_confirm smallint,
                               confirm_status smallint,
                               processing_status public.processing_status,
                               status smallint,
                               lines jsonb,
                               discounts jsonb,
                               total_items integer,
                               basket_value integer,
                               total_weight integer,
                               total_tax integer,
                               total_discount integer,
                               total_amount integer,
                               order_note text,
                               shop_note text,
                               shipping_note text,
                               order_source_id bigint,
                               order_source_type public.order_source_type,
                               external_order_id text,
                               fulfillment_shipping_status smallint,
                               customer_payment_status smallint,
                               shop_shipping_fee integer,
                               reference_url text,
                               external_order_source text,
                               shop_address jsonb,
                               shop_shipping jsonb,
                               shop_cod integer,
                               is_outside_etop boolean,
                               ghn_note_code public.ghn_note_code,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               ed_code text,
                               order_discount integer,
                               etop_payment_status smallint,
                               fulfillment_shipping_states text[],
                               fulfillment_payment_statuses smallint[],
                               fulfillment_shipping_codes text[],
                               customer_name text,
                               customer_name_norm tsvector,
                               product_name_norm tsvector,
                               fulfillment_sync_statuses smallint[],
                               partner_id bigint,
                               try_on public.try_on,
                               total_fee integer,
                               fee_lines jsonb,
                               external_url text,
                               fulfillment_type smallint,
                               fulfillment_ids bigint[],
                               external_meta jsonb,
                               trading_shop_id bigint,
                               payment_status smallint,
                               payment_id bigint,
                               referral_meta jsonb,
                               customer_id bigint,
                               created_by bigint,
                               fulfillment_statuses smallint[]
);


ALTER TABLE history."order" OWNER TO etop;

--
-- Name: order_external; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.order_external (
                                      id bigint,
                                      rid bigint NOT NULL,
                                      order_source_id bigint,
                                      external_order_source_type public.order_source_type,
                                      external_provider text,
                                      external_shop_id text,
                                      external_order_id text,
                                      external_order_code text,
                                      external_user_id text,
                                      external_customer_id text,
                                      external_created_at timestamp with time zone,
                                      external_updated_at timestamp with time zone,
                                      external_processed_at timestamp with time zone,
                                      external_closed_at timestamp with time zone,
                                      external_cancelled_at timestamp with time zone,
                                      external_cancel_reason text,
                                      external_data jsonb,
                                      external_lines jsonb,
                                      external_order_source text,
                                      _time timestamp with time zone NOT NULL,
                                      _op public.tg_op_type NOT NULL
);


ALTER TABLE history.order_external OWNER TO etop;

--
-- Name: order_line; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.order_line (
                                  rid bigint NOT NULL,
                                  order_id bigint,
                                  product_id bigint,
                                  variant_id bigint,
                                  supplier_id bigint,
                                  external_variant_id text,
                                  external_supplier_order_id text,
                                  product_name text,
                                  supplier_name text,
                                  image_url text,
                                  created_at timestamp with time zone,
                                  updated_at timestamp with time zone,
                                  closed_at timestamp with time zone,
                                  confirmed_at timestamp with time zone,
                                  cancelled_at timestamp with time zone,
                                  cancel_reason text,
                                  supplier_confirm smallint,
                                  status smallint,
                                  weight integer,
                                  quantity integer,
                                  wholesale_price_0 integer,
                                  wholesale_price integer,
                                  list_price integer,
                                  retail_price integer,
                                  payment_price integer,
                                  line_amount integer,
                                  total_discount integer,
                                  total_line_amount integer,
                                  requires_shipping boolean,
                                  is_outside_etop boolean,
                                  code text,
                                  shop_id bigint,
                                  _time timestamp with time zone NOT NULL,
                                  _op public.tg_op_type NOT NULL,
                                  is_free boolean,
                                  meta_fields jsonb
);


ALTER TABLE history.order_line OWNER TO etop;

--
-- Name: order_source; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.order_source (
                                    id bigint,
                                    rid bigint NOT NULL,
                                    type public.order_source_type,
                                    name text,
                                    status smallint,
                                    external_key text,
                                    external_info jsonb,
                                    extra_info jsonb,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    shop_id bigint,
                                    external_status smallint,
                                    last_sync_at timestamp with time zone,
                                    sync_state_orders jsonb,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL
);


ALTER TABLE history.order_source OWNER TO etop;

--
-- Name: order_source_internal; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.order_source_internal (
                                             id bigint,
                                             rid bigint NOT NULL,
                                             secret jsonb,
                                             access_token text,
                                             expires_at timestamp with time zone,
                                             created_at timestamp with time zone,
                                             updated_at timestamp with time zone,
                                             _time timestamp with time zone NOT NULL,
                                             _op public.tg_op_type NOT NULL
);


ALTER TABLE history.order_source_internal OWNER TO etop;

--
-- Name: partner; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.partner (
                               id bigint,
                               rid bigint NOT NULL,
                               name text,
                               public_name text,
                               owner_id bigint,
                               status smallint,
                               is_test smallint,
                               phone text,
                               email text,
                               website_url text,
                               image_url text,
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               deleted_at timestamp with time zone,
                               contact_persons jsonb,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               recognized_hosts text[],
                               redirect_urls text[],
                               available_from_etop boolean,
                               available_from_etop_config jsonb,
                               white_label_key text
);


ALTER TABLE history.partner OWNER TO etop;

--
-- Name: partner_relation; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.partner_relation (
                                        partner_id bigint,
                                        subject_id bigint,
                                        subject_type public.subject_type,
                                        external_subject_id text,
                                        nonce bigint,
                                        status smallint,
                                        roles text[],
                                        permissions text[],
                                        created_at timestamp with time zone,
                                        updated_at timestamp with time zone,
                                        deleted_at timestamp with time zone,
                                        rid bigint NOT NULL,
                                        _time timestamp with time zone NOT NULL,
                                        _op public.tg_op_type NOT NULL,
                                        auth_key text
);


ALTER TABLE history.partner_relation OWNER TO etop;

--
-- Name: payment; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.payment (
                               id bigint,
                               amount integer,
                               status smallint,
                               state text,
                               payment_provider text,
                               external_trans_id text,
                               external_data jsonb,
                               created_at timestamp without time zone,
                               updated_at timestamp without time zone,
                               rid bigint NOT NULL,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL
);


ALTER TABLE history.payment OWNER TO etop;

--
-- Name: product; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product (
                               id bigint,
                               rid bigint NOT NULL,
                               product_source_id bigint,
                               supplier_id bigint,
                               product_source_category_id bigint,
                               etop_category_id bigint,
                               product_brand_id bigint,
                               name text,
                               short_desc text,
                               description text,
                               desc_html text,
                               ed_name text,
                               ed_short_desc text,
                               ed_description text,
                               ed_desc_html text,
                               status smallint,
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               deleted_at timestamp with time zone,
                               code text,
                               quantity_available integer,
                               quantity_on_hand integer,
                               quantity_reserved integer,
                               image_urls text[],
                               external_id text,
                               ed_tags text[],
                               name_norm tsvector,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               ed_code text,
                               unit text,
                               name_norm_ua text
);


ALTER TABLE history.product OWNER TO etop;

--
-- Name: product_brand; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_brand (
                                     id bigint,
                                     name text,
                                     description text,
                                     policy text,
                                     supplier_id bigint,
                                     image_urls text[],
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     rid bigint NOT NULL,
                                     _time timestamp with time zone NOT NULL,
                                     _op public.tg_op_type NOT NULL
);


ALTER TABLE history.product_brand OWNER TO etop;

--
-- Name: product_external; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_external (
                                        id bigint,
                                        rid bigint NOT NULL,
                                        product_source_id bigint,
                                        product_source_type public.product_source_type,
                                        external_id text,
                                        external_name text,
                                        external_code text,
                                        external_category_id text,
                                        external_description text,
                                        external_image_urls text[],
                                        external_unit text,
                                        external_data jsonb,
                                        external_status smallint,
                                        external_created_at timestamp with time zone,
                                        external_updated_at timestamp with time zone,
                                        external_deleted_at timestamp with time zone,
                                        last_sync_at timestamp with time zone,
                                        external_units text[],
                                        _time timestamp with time zone NOT NULL,
                                        _op public.tg_op_type NOT NULL
);


ALTER TABLE history.product_external OWNER TO etop;

--
-- Name: product_shop_collection; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_shop_collection (
                                               product_id bigint,
                                               shop_id bigint,
                                               collection_id bigint,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               status smallint,
                                               rid bigint NOT NULL,
                                               _time timestamp with time zone NOT NULL,
                                               _op public.tg_op_type NOT NULL,
                                               partner_id bigint,
                                               external_collection_id text,
                                               external_product_id text,
                                               deleted_at timestamp with time zone
);


ALTER TABLE history.product_shop_collection OWNER TO etop;

--
-- Name: product_source; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_source (
                                      id bigint,
                                      rid bigint NOT NULL,
                                      type public.product_source_type,
                                      name text,
                                      status smallint,
                                      external_key text,
                                      external_info jsonb,
                                      extra_info jsonb,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      supplier_id bigint,
                                      last_sync_at timestamp with time zone,
                                      sync_state_products jsonb,
                                      sync_state_categories jsonb,
                                      external_status smallint,
                                      _time timestamp with time zone NOT NULL,
                                      _op public.tg_op_type NOT NULL
);


ALTER TABLE history.product_source OWNER TO etop;

--
-- Name: product_source_category; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_source_category (
                                               id bigint,
                                               rid bigint NOT NULL,
                                               product_source_id bigint,
                                               product_source_type public.product_source_type,
                                               supplier_id bigint,
                                               parent_id bigint,
                                               name text,
                                               status smallint,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               deleted_at timestamp with time zone,
                                               shop_id bigint,
                                               _time timestamp with time zone NOT NULL,
                                               _op public.tg_op_type NOT NULL,
                                               external_id text,
                                               external_parent_id text,
                                               partner_id bigint
);


ALTER TABLE history.product_source_category OWNER TO etop;

--
-- Name: product_source_category_external; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_source_category_external (
                                                        id bigint,
                                                        rid bigint NOT NULL,
                                                        product_source_id bigint,
                                                        product_source_type public.product_source_type,
                                                        external_id text,
                                                        external_parent_id text,
                                                        external_code text,
                                                        external_name text,
                                                        external_status smallint,
                                                        external_updated_at timestamp with time zone,
                                                        external_created_at timestamp with time zone,
                                                        external_deleted_at timestamp with time zone,
                                                        last_sync_at timestamp with time zone,
                                                        _time timestamp with time zone NOT NULL,
                                                        _op public.tg_op_type NOT NULL
);


ALTER TABLE history.product_source_category_external OWNER TO etop;

--
-- Name: product_source_internal; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.product_source_internal (
                                               id bigint,
                                               rid bigint NOT NULL,
                                               type public.product_source_type,
                                               status smallint,
                                               external_status smallint,
                                               secret jsonb,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               last_sync_at timestamp with time zone,
                                               sync_state_products jsonb,
                                               sync_state_categories jsonb,
                                               access_token text,
                                               expires_at timestamp with time zone,
                                               _time timestamp with time zone NOT NULL,
                                               _op public.tg_op_type NOT NULL
);


ALTER TABLE history.product_source_internal OWNER TO etop;

--
-- Name: purchase_order; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.purchase_order (
                                      id bigint,
                                      shop_id bigint,
                                      supplier_id bigint,
                                      supplier jsonb,
                                      basket_value bigint,
                                      total_discount bigint,
                                      total_amount bigint,
                                      paid_amount bigint,
                                      code text,
                                      code_norm integer,
                                      note text,
                                      status smallint,
                                      variant_ids bigint[],
                                      lines jsonb,
                                      cancelled_reason text,
                                      created_by bigint,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      confirmed_at timestamp with time zone,
                                      cancelled_at timestamp with time zone,
                                      deleted_at timestamp with time zone,
                                      rid bigint NOT NULL,
                                      _time timestamp with time zone NOT NULL,
                                      _op public.tg_op_type NOT NULL,
                                      supplier_full_name_norm text,
                                      supplier_phone_norm text,
                                      discount_lines jsonb,
                                      total_fee integer,
                                      fee_lines jsonb
);


ALTER TABLE history.purchase_order OWNER TO etop;

--
-- Name: purchase_refund; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.purchase_refund (
                                       id bigint,
                                       shop_id bigint,
                                       purchase_order_id bigint,
                                       note text,
                                       code_norm integer,
                                       code text,
                                       supplier_id bigint,
                                       lines jsonb,
                                       created_at timestamp with time zone,
                                       updated_at timestamp with time zone,
                                       cancelled_at timestamp with time zone,
                                       confirmed_at timestamp with time zone,
                                       created_by bigint,
                                       updated_by bigint,
                                       total_amount integer,
                                       basket_value integer,
                                       cancel_reason text,
                                       status integer,
                                       adjustment_lines jsonb,
                                       total_adjustment integer,
                                       rid bigint NOT NULL,
                                       _time timestamp with time zone NOT NULL,
                                       _op public.tg_op_type NOT NULL
);


ALTER TABLE history.purchase_refund OWNER TO etop;

--
-- Name: receipt; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.receipt (
                               id bigint,
                               shop_id bigint,
                               trader_id bigint,
                               created_by bigint,
                               code text,
                               title text,
                               description text,
                               amount integer,
                               status smallint,
                               type public.receipt_type,
                               lines jsonb,
                               ref_ids bigint[],
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               deleted_at timestamp with time zone,
                               rid bigint NOT NULL,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               shop_ledger_id bigint,
                               created_type public.receipt_created_type,
                               ledger_id bigint,
                               cancelled_reason text,
                               code_norm integer,
                               paid_at timestamp with time zone,
                               confirmed_at timestamp with time zone,
                               cancelled_at timestamp with time zone,
                               ref_type public.receipt_ref_type,
                               trader jsonb,
                               trader_full_name_norm tsvector,
                               trader_phone_norm tsvector,
                               trader_type public.trader_type
);


ALTER TABLE history.receipt OWNER TO etop;

--
-- Name: refund; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.refund (
                              id bigint,
                              shop_id bigint,
                              order_id bigint,
                              note text,
                              code_norm integer,
                              code text,
                              customer_id bigint,
                              lines jsonb,
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              cancelled_at timestamp with time zone,
                              confirmed_at timestamp with time zone,
                              created_by bigint,
                              updated_by bigint,
                              total_amount integer,
                              basket_value integer,
                              cancel_reason text,
                              status integer,
                              adjustment_lines jsonb,
                              total_adjustment integer,
                              rid bigint NOT NULL,
                              _time timestamp with time zone NOT NULL,
                              _op public.tg_op_type NOT NULL
);


ALTER TABLE history.refund OWNER TO etop;

--
-- Name: shipnow_fulfillment; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shipnow_fulfillment (
                                           id bigint,
                                           shop_id bigint,
                                           partner_id bigint,
                                           order_ids bigint[],
                                           pickup_address jsonb,
                                           carrier text,
                                           shipping_service_code text,
                                           shipping_service_fee integer,
                                           chargeable_weight integer,
                                           gross_weight integer,
                                           basket_value integer,
                                           cod_amount integer,
                                           shipping_note text,
                                           request_pickup_at timestamp with time zone,
                                           delivery_points jsonb,
                                           status integer,
                                           shipping_state text,
                                           sync_status integer,
                                           sync_states jsonb,
                                           last_sync_at timestamp with time zone,
                                           created_at timestamp with time zone,
                                           updated_at timestamp with time zone,
                                           confirm_status integer,
                                           shipping_status integer,
                                           shipping_code text,
                                           fee_lines jsonb,
                                           carrier_fee_lines jsonb,
                                           total_fee integer,
                                           shipping_created_at timestamp with time zone,
                                           etop_payment_status integer,
                                           cod_etop_transfered_at timestamp with time zone,
                                           shipping_picking_at timestamp with time zone,
                                           shipping_delivering_at timestamp with time zone,
                                           shipping_delivered_at timestamp with time zone,
                                           shipping_cancelled_at timestamp with time zone,
                                           shipping_service_name text,
                                           shipping_service_description text,
                                           cancel_reason text,
                                           shipping_shared_link text,
                                           address_to_province_code text,
                                           address_to_district_code text,
                                           rid bigint NOT NULL,
                                           _time timestamp with time zone NOT NULL,
                                           _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shipnow_fulfillment OWNER TO etop;

--
-- Name: shop; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop (
                            id bigint,
                            rid bigint NOT NULL,
                            name text,
                            owner_id bigint,
                            status smallint,
                            product_source_id bigint,
                            created_at timestamp with time zone,
                            updated_at timestamp with time zone,
                            rules jsonb,
                            is_test smallint,
                            image_url text,
                            phone text,
                            website_url text,
                            email text,
                            deleted_at timestamp with time zone,
                            address_id bigint,
                            bank_account jsonb,
                            contact_persons jsonb,
                            order_source_id bigint,
                            ship_to_address_id bigint,
                            ship_from_address_id bigint,
                            code text,
                            auto_create_ffm boolean,
                            ghn_note_code public.ghn_note_code,
                            _time timestamp with time zone NOT NULL,
                            _op public.tg_op_type NOT NULL,
                            try_on public.try_on,
                            recognized_hosts text[],
                            company_info jsonb,
                            money_transaction_rrule text,
                            survey_info jsonb,
                            shipping_service_select_strategy jsonb,
                            inventory_overstock boolean,
                            wl_partner_id bigint
);


ALTER TABLE history.shop OWNER TO etop;

--
-- Name: shop_brand; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_brand (
                                  shop_id bigint,
                                  id bigint,
                                  description text,
                                  brand_name text,
                                  updated_at timestamp with time zone,
                                  created_at timestamp with time zone,
                                  deleted_at timestamp with time zone,
                                  rid bigint NOT NULL,
                                  _time timestamp with time zone NOT NULL,
                                  _op public.tg_op_type NOT NULL,
                                  external_id text,
                                  partner_id bigint
);


ALTER TABLE history.shop_brand OWNER TO etop;

--
-- Name: shop_carrier; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_carrier (
                                    id bigint,
                                    shop_id bigint,
                                    full_name text,
                                    note text,
                                    status smallint,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    deleted_at timestamp with time zone,
                                    rid bigint NOT NULL,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_carrier OWNER TO etop;

--
-- Name: shop_collection; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_collection (
                                       id bigint,
                                       rid bigint NOT NULL,
                                       shop_id bigint,
                                       name text,
                                       description text,
                                       desc_html text,
                                       short_desc text,
                                       created_at timestamp with time zone,
                                       updated_at timestamp with time zone,
                                       _time timestamp with time zone NOT NULL,
                                       _op public.tg_op_type NOT NULL,
                                       deleted_at timestamp with time zone,
                                       partner_id bigint,
                                       external_id text
);


ALTER TABLE history.shop_collection OWNER TO etop;

--
-- Name: shop_connection; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_connection (
                                       shop_id bigint,
                                       connection_id bigint,
                                       token text,
                                       token_expires_at timestamp with time zone,
                                       status smallint,
                                       connection_states jsonb,
                                       is_global boolean,
                                       created_at timestamp with time zone,
                                       updated_at timestamp with time zone,
                                       deleted_at timestamp with time zone,
                                       external_data jsonb,
                                       rid bigint NOT NULL,
                                       _time timestamp with time zone NOT NULL,
                                       _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_connection OWNER TO etop;

--
-- Name: shop_customer; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_customer (
                                     id bigint,
                                     shop_id bigint,
                                     code text,
                                     full_name text,
                                     gender public.gender_type,
                                     type public.customer_type,
                                     birthday date,
                                     note text,
                                     phone text,
                                     email text,
                                     status smallint,
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     deleted_at timestamp with time zone,
                                     rid bigint NOT NULL,
                                     _time timestamp with time zone NOT NULL,
                                     _op public.tg_op_type NOT NULL,
                                     full_name_norm tsvector,
                                     phone_norm tsvector,
                                     code_norm integer,
                                     external_id text,
                                     external_code text,
                                     partner_id bigint
);


ALTER TABLE history.shop_customer OWNER TO etop;

--
-- Name: shop_customer_group; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_customer_group (
                                           id bigint,
                                           name text,
                                           created_at timestamp with time zone,
                                           updated_at timestamp with time zone,
                                           deleted_at timestamp with time zone,
                                           shop_id bigint,
                                           rid bigint NOT NULL,
                                           _time timestamp with time zone NOT NULL,
                                           _op public.tg_op_type NOT NULL,
                                           partner_id bigint
);


ALTER TABLE history.shop_customer_group OWNER TO etop;

--
-- Name: shop_customer_group_customer; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_customer_group_customer (
                                                    customer_id bigint,
                                                    group_id bigint,
                                                    created_at timestamp with time zone,
                                                    updated_at timestamp with time zone,
                                                    rid bigint NOT NULL,
                                                    _time timestamp with time zone NOT NULL,
                                                    _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_customer_group_customer OWNER TO etop;

--
-- Name: shop_ledger; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_ledger (
                                   id bigint,
                                   shop_id bigint,
                                   name text,
                                   bank_account jsonb,
                                   status smallint,
                                   note text,
                                   type public.shop_ledger_type,
                                   created_by bigint,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   deleted_at timestamp with time zone,
                                   rid bigint NOT NULL,
                                   _time timestamp with time zone NOT NULL,
                                   _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_ledger OWNER TO etop;

--
-- Name: shop_product; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_product (
                                    shop_id bigint,
                                    product_id bigint,
                                    rid bigint NOT NULL,
                                    collection_id bigint,
                                    name text,
                                    description text,
                                    desc_html text,
                                    short_desc text,
                                    retail_price integer,
                                    tags text[],
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    note text,
                                    image_urls text[],
                                    status smallint,
                                    haravan_id text,
                                    name_norm tsvector,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL,
                                    deleted_at timestamp with time zone,
                                    code text,
                                    name_norm_ua text,
                                    category_id bigint,
                                    cost_price integer,
                                    list_price integer,
                                    unit text,
                                    product_type text,
                                    supplier_id bigint,
                                    meta_fields jsonb,
                                    brand_id bigint,
                                    external_id text,
                                    external_code text,
                                    partner_id bigint,
                                    code_norm integer,
                                    external_brand_id text,
                                    external_category_id text
);


ALTER TABLE history.shop_product OWNER TO etop;

--
-- Name: shop_product_collection; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_product_collection (
                                               product_id bigint,
                                               collection_id bigint,
                                               shop_id bigint,
                                               status smallint,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               partner_id bigint,
                                               external_collection_id text,
                                               external_product_id text,
                                               deleted_at timestamp with time zone,
                                               rid bigint NOT NULL,
                                               _time timestamp with time zone NOT NULL,
                                               _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_product_collection OWNER TO etop;

--
-- Name: shop_stocktake; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_stocktake (
                                      id bigint,
                                      shop_id bigint,
                                      created_by bigint,
                                      updated_by bigint,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      confirmed_at timestamp with time zone,
                                      cancelled_at timestamp with time zone,
                                      variant_ids bigint[],
                                      total_quantity integer,
                                      status integer,
                                      lines jsonb,
                                      code text,
                                      code_norm integer,
                                      note text,
                                      cancel_reason text,
                                      product_ids bigint[],
                                      type text,
                                      rid bigint NOT NULL,
                                      _time timestamp with time zone NOT NULL,
                                      _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_stocktake OWNER TO etop;

--
-- Name: shop_supplier; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_supplier (
                                     id bigint,
                                     shop_id bigint,
                                     full_name text,
                                     note text,
                                     status smallint,
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     deleted_at timestamp with time zone,
                                     rid bigint NOT NULL,
                                     phone text,
                                     email text,
                                     company_name text,
                                     tax_number text,
                                     headquater_address text,
                                     full_name_norm tsvector,
                                     phone_norm tsvector,
                                     _time timestamp with time zone NOT NULL,
                                     _op public.tg_op_type NOT NULL,
                                     code_norm integer,
                                     code text,
                                     company_name_norm tsvector
);


ALTER TABLE history.shop_supplier OWNER TO etop;

--
-- Name: shop_trader; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_trader (
                                   id bigint,
                                   shop_id bigint,
                                   rid bigint NOT NULL,
                                   _time timestamp with time zone NOT NULL,
                                   _op public.tg_op_type NOT NULL,
                                   type public.trader_type,
                                   deleted_at timestamp with time zone
);


ALTER TABLE history.shop_trader OWNER TO etop;

--
-- Name: shop_trader_address; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_trader_address (
                                           id bigint,
                                           shop_id bigint,
                                           trader_id bigint,
                                           full_name text,
                                           phone text,
                                           email text,
                                           company text,
                                           district_code text,
                                           ward_code text,
                                           city text,
                                           address1 text,
                                           address2 text,
                                           "position" text,
                                           note text,
                                           "primary" boolean,
                                           status smallint,
                                           coordinates jsonb,
                                           created_at timestamp with time zone,
                                           updated_at timestamp with time zone,
                                           deleted_at timestamp with time zone,
                                           rid bigint NOT NULL,
                                           _time timestamp with time zone NOT NULL,
                                           _op public.tg_op_type NOT NULL,
                                           is_default boolean,
                                           partner_id bigint
);


ALTER TABLE history.shop_trader_address OWNER TO etop;

--
-- Name: shop_variant; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_variant (
                                    shop_id bigint,
                                    variant_id bigint,
                                    rid bigint NOT NULL,
                                    collection_id bigint,
                                    name text,
                                    description text,
                                    desc_html text,
                                    short_desc text,
                                    retail_price integer,
                                    tags text[],
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    note text,
                                    image_urls text[],
                                    status smallint,
                                    haravan_id text,
                                    name_norm tsvector,
                                    product_id bigint,
                                    _time timestamp with time zone NOT NULL,
                                    _op public.tg_op_type NOT NULL,
                                    deleted_at timestamp with time zone,
                                    attr_norm_kv tsvector,
                                    code text,
                                    cost_price integer,
                                    list_price integer,
                                    attributes jsonb,
                                    external_id text,
                                    external_code text,
                                    partner_id bigint,
                                    code_norm integer,
                                    external_product_id text
);


ALTER TABLE history.shop_variant OWNER TO etop;

--
-- Name: shop_variant_supplier; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_variant_supplier (
                                             shop_id bigint,
                                             supplier_id bigint,
                                             variant_id bigint,
                                             created_at timestamp with time zone,
                                             updated_at timestamp with time zone,
                                             rid bigint NOT NULL,
                                             _time timestamp with time zone NOT NULL,
                                             _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_variant_supplier OWNER TO etop;

--
-- Name: shop_vendor; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.shop_vendor (
                                   id bigint,
                                   shop_id bigint,
                                   full_name text,
                                   note text,
                                   status smallint,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   deleted_at timestamp with time zone,
                                   rid bigint NOT NULL,
                                   _time timestamp with time zone NOT NULL,
                                   _op public.tg_op_type NOT NULL
);


ALTER TABLE history.shop_vendor OWNER TO etop;

--
-- Name: supplier; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.supplier (
                                id bigint,
                                rid bigint NOT NULL,
                                status smallint,
                                created_at timestamp with time zone,
                                updated_at timestamp with time zone,
                                name text,
                                owner_id bigint,
                                product_source_id bigint,
                                rules jsonb,
                                is_test smallint,
                                company_info jsonb,
                                warehouse_address_id bigint,
                                bank_account jsonb,
                                contact_persons jsonb,
                                image_url text,
                                deleted_at timestamp with time zone,
                                ship_from_address_id bigint,
                                _time timestamp with time zone NOT NULL,
                                _op public.tg_op_type NOT NULL
);


ALTER TABLE history.supplier OWNER TO etop;

--
-- Name: transaction; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.transaction (
                                   id bigint,
                                   amount integer,
                                   account_id bigint,
                                   status smallint,
                                   type text,
                                   note text,
                                   metadata jsonb,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   rid bigint NOT NULL,
                                   _time timestamp with time zone NOT NULL,
                                   _op public.tg_op_type NOT NULL
);


ALTER TABLE history.transaction OWNER TO etop;

--
-- Name: user; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history."user" (
                              id bigint,
                              rid bigint NOT NULL,
                              status smallint,
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              full_name text,
                              short_name text,
                              email text,
                              phone text,
                              _time timestamp with time zone NOT NULL,
                              _op public.tg_op_type NOT NULL,
                              agreed_tos_at timestamp with time zone,
                              agreed_email_info_at timestamp with time zone,
                              email_verified_at timestamp with time zone,
                              phone_verified_at timestamp with time zone,
                              email_verification_sent_at timestamp with time zone,
                              phone_verification_sent_at timestamp with time zone,
                              is_test smallint,
                              identifying public.user_identifying_type,
                              source text,
                              ref_user_id bigint,
                              ref_sale_id bigint,
                              wl_partner_id bigint
);


ALTER TABLE history."user" OWNER TO etop;

--
-- Name: user_internal; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.user_internal (
                                     id bigint,
                                     hashpwd text,
                                     updated_at timestamp with time zone,
                                     rid bigint NOT NULL,
                                     _time timestamp with time zone NOT NULL,
                                     _op public.tg_op_type NOT NULL
);


ALTER TABLE history.user_internal OWNER TO etop;

--
-- Name: variant; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.variant (
                               id bigint,
                               rid bigint NOT NULL,
                               product_id bigint,
                               product_source_id bigint,
                               supplier_id bigint,
                               name text,
                               short_desc text,
                               description text,
                               desc_html text,
                               ed_name text,
                               ed_short_desc text,
                               ed_description text,
                               ed_desc_html text,
                               desc_norm tsvector,
                               name_norm tsvector,
                               status smallint,
                               etop_status smallint,
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               deleted_at timestamp with time zone,
                               sku text,
                               code text,
                               wholesale_price integer,
                               wholesale_price_0 integer,
                               list_price integer,
                               retail_price_min integer,
                               retail_price_max integer,
                               ed_wholesale_price integer,
                               ed_wholesale_price_0 integer,
                               ed_list_price integer,
                               ed_retail_price_min integer,
                               ed_retail_price_max integer,
                               quantity_available integer,
                               quantity_on_hand integer,
                               quantity_reserved integer,
                               image_urls text[],
                               supplier_meta jsonb,
                               product_source_category_id bigint,
                               etop_category_id bigint,
                               product_brand_id bigint,
                               ed_status smallint,
                               attributes jsonb,
                               cost_price integer,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL,
                               ed_code text,
                               attr_norm_kv tsvector
);


ALTER TABLE history.variant OWNER TO etop;

--
-- Name: variant_external; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.variant_external (
                                        id bigint,
                                        rid bigint NOT NULL,
                                        product_source_id bigint,
                                        product_source_type public.product_source_type,
                                        external_id text,
                                        external_product_id text,
                                        external_category_id text,
                                        external_code text,
                                        external_name text,
                                        external_description text,
                                        external_attributes jsonb,
                                        external_image_urls text[],
                                        external_unit text,
                                        external_conversion_value real,
                                        external_price integer,
                                        external_unit_conv real,
                                        external_base_unit_id text,
                                        external_data jsonb,
                                        external_name_norm tsvector,
                                        external_status smallint,
                                        external_created_at timestamp with time zone,
                                        external_updated_at timestamp with time zone,
                                        external_deleted_at timestamp with time zone,
                                        last_sync_at timestamp with time zone,
                                        external_units text[],
                                        product_id bigint,
                                        _time timestamp with time zone NOT NULL,
                                        _op public.tg_op_type NOT NULL
);


ALTER TABLE history.variant_external OWNER TO etop;

--
-- Name: webhook; Type: TABLE; Schema: history; Owner: etop
--

CREATE TABLE history.webhook (
                               id bigint,
                               account_id bigint,
                               entities text[],
                               fields text[],
                               url text,
                               metadata text,
                               created_at timestamp with time zone,
                               updated_at timestamp with time zone,
                               deleted_at timestamp with time zone,
                               rid bigint NOT NULL,
                               _time timestamp with time zone NOT NULL,
                               _op public.tg_op_type NOT NULL
);


ALTER TABLE history.webhook OWNER TO etop;

--
-- Name: account; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.account (
                              id bigint NOT NULL,
                              name text NOT NULL,
                              deleted_at timestamp with time zone,
                              image_url text,
                              type public.account_type,
                              rid bigint NOT NULL,
                              url_slug text,
                              owner_id bigint
);


ALTER TABLE public.account OWNER TO etop;

--
-- Name: account_auth; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.account_auth (
                                   auth_key text NOT NULL,
                                   account_id bigint NOT NULL,
                                   status smallint,
                                   roles text[],
                                   permissions text[],
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   deleted_at timestamp with time zone,
                                   rid bigint NOT NULL
);


ALTER TABLE public.account_auth OWNER TO etop;

--
-- Name: account_user; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.account_user (
                                   account_id bigint NOT NULL,
                                   user_id bigint NOT NULL,
                                   status smallint NOT NULL,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   roles text[],
                                   permissions text[],
                                   deleted_at timestamp with time zone,
                                   rid bigint NOT NULL,
                                   full_name text,
                                   short_name text,
                                   "position" text,
                                   response_status smallint,
                                   invitation_sent_at timestamp with time zone,
                                   invitation_sent_by bigint,
                                   invitation_accepted_at timestamp with time zone,
                                   invitation_rejected_at timestamp with time zone,
                                   disabled_at timestamp with time zone,
                                   disabled_by bigint,
                                   disable_reason text
);


ALTER TABLE public.account_user OWNER TO etop;

--
-- Name: address; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.address (
                              id bigint NOT NULL,
                              country text DEFAULT 'VN'::text,
                              province_code text,
                              province text,
                              district_code text,
                              district text,
                              ward text,
                              ward_code text,
                              address1 text,
                              is_default boolean DEFAULT false,
                              type public.address_type,
                              account_id bigint,
                              created_at timestamp with time zone NOT NULL,
                              updated_at timestamp with time zone NOT NULL,
                              full_name text,
                              first_name text,
                              last_name text,
                              email text,
                              "position" text,
                              city text,
                              zip text,
                              address2 text,
                              phone text,
                              company text,
                              notes jsonb,
                              rid bigint NOT NULL,
                              coordinates jsonb
);


ALTER TABLE public.address OWNER TO etop;

--
-- Name: affiliate; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.affiliate (
                                id bigint NOT NULL,
                                rid bigint,
                                name text NOT NULL,
                                owner_id bigint NOT NULL,
                                status smallint,
                                is_test smallint DEFAULT '0'::smallint,
                                phone text,
                                email text,
                                created_at timestamp with time zone,
                                updated_at timestamp with time zone,
                                deleted_at timestamp with time zone,
                                bank_account jsonb
);


ALTER TABLE public.affiliate OWNER TO etop;

--
-- Name: code; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.code (
                           code text NOT NULL,
                           type public.code_type NOT NULL,
                           created_at timestamp with time zone
);


ALTER TABLE public.code OWNER TO etop;

--
-- Name: connection; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.connection (
                                 id bigint NOT NULL,
                                 name text,
                                 status smallint,
                                 partner_id bigint,
                                 created_at timestamp with time zone,
                                 updated_at timestamp with time zone,
                                 deleted_at timestamp with time zone,
                                 driver_config jsonb,
                                 driver text,
                                 connection_type text,
                                 connection_subtype text,
                                 connection_method text,
                                 connection_provider text,
                                 etop_affiliate_account jsonb,
                                 code text,
                                 image_url text,
                                 rid bigint NOT NULL
);


ALTER TABLE public.connection OWNER TO etop;

--
-- Name: credit; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.credit (
                             id bigint NOT NULL,
                             rid bigint NOT NULL,
                             amount integer,
                             shop_id bigint,
                             supplier_id bigint,
                             created_at timestamp with time zone,
                             updated_at timestamp with time zone,
                             paid_at timestamp with time zone,
                             type public.account_type,
                             status smallint
);


ALTER TABLE public.credit OWNER TO etop;

--
-- Name: district; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.district (
                               code text NOT NULL,
                               province_code text NOT NULL,
                               urban smallint NOT NULL,
                               name text NOT NULL,
                               ghn_id integer,
                               vtpost_id integer
);


ALTER TABLE public.district OWNER TO etop;

--
-- Name: export_attempt; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.export_attempt (
                                     id text NOT NULL,
                                     user_id bigint NOT NULL,
                                     account_id bigint NOT NULL,
                                     export_type text NOT NULL,
                                     filename text,
                                     stored_file text,
                                     download_url text,
                                     request_query text,
                                     mime_type text,
                                     status smallint NOT NULL,
                                     errors jsonb,
                                     error jsonb,
                                     n_total integer,
                                     n_exported integer,
                                     n_error integer,
                                     created_at timestamp with time zone,
                                     deleted_at timestamp with time zone,
                                     started_at timestamp with time zone,
                                     done_at timestamp with time zone,
                                     expires_at timestamp with time zone
);


ALTER TABLE public.export_attempt OWNER TO etop;

--
-- Name: external_account_ahamove; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.external_account_ahamove (
                                               id bigint NOT NULL,
                                               owner_id bigint NOT NULL,
                                               phone text NOT NULL,
                                               name text NOT NULL,
                                               external_token text,
                                               external_verified boolean,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               external_created_at timestamp with time zone,
                                               last_send_verified_at timestamp with time zone,
                                               external_ticket_id text,
                                               external_id text,
                                               id_card_front_img text,
                                               id_card_back_img text,
                                               portrait_img text,
                                               uploaded_at timestamp with time zone,
                                               fanpage_url text,
                                               business_license_imgs text[],
                                               website_url text,
                                               company_imgs text[],
                                               external_data_verified jsonb
);


ALTER TABLE public.external_account_ahamove OWNER TO etop;

--
-- Name: external_account_ahamove_view; Type: VIEW; Schema: public; Owner: etop
--

CREATE VIEW public.external_account_ahamove_view AS
SELECT external_account_ahamove.id,
       external_account_ahamove.owner_id,
       external_account_ahamove.phone,
       external_account_ahamove.name,
       external_account_ahamove.external_verified,
       external_account_ahamove.created_at,
       external_account_ahamove.updated_at,
       external_account_ahamove.external_created_at,
       external_account_ahamove.last_send_verified_at,
       external_account_ahamove.external_ticket_id,
       external_account_ahamove.external_id,
       (external_account_ahamove.id_card_front_img IS NOT NULL) AS id_card_front_img_uploaded,
       (external_account_ahamove.id_card_back_img IS NOT NULL) AS id_card_back_img_uploaded,
       (external_account_ahamove.portrait_img IS NOT NULL) AS portrait_img_uploaded,
       external_account_ahamove.uploaded_at
FROM public.external_account_ahamove;


ALTER TABLE public.external_account_ahamove_view OWNER TO etop;

--
-- Name: external_account_haravan; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.external_account_haravan (
                                               id bigint NOT NULL,
                                               shop_id bigint NOT NULL,
                                               subdomain text,
                                               external_shop_id integer,
                                               external_carrier_service_id integer,
                                               external_connected_carrier_service_at timestamp with time zone,
                                               access_token text,
                                               expires_at timestamp with time zone,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone
);


ALTER TABLE public.external_account_haravan OWNER TO etop;

--
-- Name: fulfillment; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.fulfillment (
                                  id bigint NOT NULL,
                                  order_id bigint NOT NULL,
                                  lines jsonb,
                                  variant_ids bigint[] NOT NULL,
                                  type_from public.fulfillment_endpoint NOT NULL,
                                  type_to public.fulfillment_endpoint NOT NULL,
                                  address_from jsonb NOT NULL,
                                  address_to jsonb NOT NULL,
                                  shop_id bigint NOT NULL,
                                  supplier_id bigint,
                                  total_items integer NOT NULL,
                                  total_weight integer NOT NULL,
                                  basket_value integer NOT NULL,
                                  total_cod_amount integer NOT NULL,
                                  shipping_fee_customer integer NOT NULL,
                                  shipping_fee_shop integer NOT NULL,
                                  external_shipping_fee integer NOT NULL,
                                  created_at timestamp with time zone NOT NULL,
                                  updated_at timestamp with time zone NOT NULL,
                                  shipping_cancelled_at timestamp with time zone,
                                  cancel_reason text,
                                  closed_at timestamp with time zone,
                                  shipping_provider public.shipping_provider NOT NULL,
                                  shipping_code text,
                                  shipping_note text,
                                  external_shipping_id text,
                                  external_shipping_code text,
                                  external_shipping_created_at timestamp with time zone,
                                  external_shipping_updated_at timestamp with time zone,
                                  external_shipping_cancelled_at timestamp with time zone,
                                  external_shipping_delivered_at timestamp with time zone,
                                  external_shipping_returned_at timestamp with time zone,
                                  external_shipping_state text,
                                  external_shipping_status smallint,
                                  external_shipping_data jsonb,
                                  shipping_state public.shipping_state,
                                  status smallint,
                                  sync_status smallint,
                                  rid bigint NOT NULL,
                                  sync_states jsonb,
                                  shipping_delivered_at timestamp with time zone,
                                  shipping_returned_at timestamp with time zone,
                                  external_shipping_closed_at timestamp with time zone,
                                  supplier_confirm smallint,
                                  shop_confirm smallint,
                                  last_sync_at timestamp with time zone,
                                  expected_delivery_at timestamp with time zone,
                                  money_transaction_id bigint,
                                  cod_etop_transfered_at timestamp with time zone,
                                  shipping_fee_shop_transfered_at timestamp with time zone,
                                  provider_shipping_fee_lines jsonb,
                                  shipping_fee_shop_lines jsonb,
                                  etop_discount integer,
                                  shipping_status smallint,
                                  etop_fee_adjustment integer,
                                  etop_payment_status smallint NOT NULL,
                                  address_to_province_code text,
                                  address_to_district_code text,
                                  address_to_ward_code text,
                                  provider_service_id text,
                                  expected_pick_at timestamp with time zone,
                                  confirm_status smallint NOT NULL,
                                  shipping_fee_main integer NOT NULL,
                                  shipping_fee_return integer NOT NULL,
                                  shipping_fee_insurance integer NOT NULL,
                                  shipping_fee_adjustment integer NOT NULL,
                                  shipping_fee_cods integer NOT NULL,
                                  shipping_fee_info_change integer NOT NULL,
                                  shipping_fee_other integer NOT NULL,
                                  external_shipping_state_code text,
                                  money_transaction_shipping_external_id bigint,
                                  total_discount integer NOT NULL,
                                  total_amount integer NOT NULL,
                                  external_shipping_logs jsonb,
                                  partner_id bigint,
                                  external_shipping_sub_state text,
                                  external_shipping_note text,
                                  try_on public.try_on,
                                  external_shipping_name text,
                                  shipping_service_fee integer,
                                  original_cod_amount integer,
                                  address_return jsonb,
                                  include_insurance boolean,
                                  admin_note text,
                                  is_partial_delivery boolean,
                                  shipping_fee_discount integer,
                                  shipping_created_at timestamp with time zone,
                                  shipping_picking_at timestamp with time zone,
                                  shipping_holding_at timestamp with time zone,
                                  shipping_delivering_at timestamp with time zone,
                                  shipping_returning_at timestamp with time zone,
                                  etop_adjusted_shipping_fee_main integer,
                                  etop_price_rule boolean,
                                  actual_compensation_amount integer,
                                  delivery_route text,
                                  created_by bigint,
                                  shipping_type smallint,
                                  connection_id bigint,
                                  connection_method text,
                                  shop_carrier_id bigint,
                                  shipping_service_name text,
                                  gross_weight integer,
                                  chargeable_weight integer,
                                  length integer,
                                  width integer,
                                  height integer,
                                  external_affiliate_id text,
                                  CONSTRAINT type_from_supplier_id CHECK (((type_from = 'supplier'::public.fulfillment_endpoint) = (supplier_id IS NOT NULL)))
);


ALTER TABLE public.fulfillment OWNER TO etop;

--
-- Name: history_account_auth_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_account_auth_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_account_auth_seq OWNER TO etop;

--
-- Name: history_account_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_account_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_account_seq OWNER TO etop;

--
-- Name: history_account_user_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_account_user_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_account_user_seq OWNER TO etop;

--
-- Name: history_address_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_address_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_address_seq OWNER TO etop;

--
-- Name: history_connection_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_connection_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_connection_seq OWNER TO etop;

--
-- Name: history_credit_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_credit_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_credit_seq OWNER TO etop;

--
-- Name: history_etop_category_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_etop_category_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_etop_category_seq OWNER TO etop;

--
-- Name: history_fulfillment_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_fulfillment_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_fulfillment_seq OWNER TO etop;

--
-- Name: history_import_attempt_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_import_attempt_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_import_attempt_seq OWNER TO etop;

--
-- Name: history_inventory_variant_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_inventory_variant_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_inventory_variant_seq OWNER TO etop;

--
-- Name: history_inventory_voucher_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_inventory_voucher_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_inventory_voucher_seq OWNER TO etop;

--
-- Name: history_invitation_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_invitation_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_invitation_seq OWNER TO etop;

--
-- Name: history_money_transaction_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_money_transaction_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_money_transaction_seq OWNER TO etop;

--
-- Name: history_order_external_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_order_external_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_order_external_seq OWNER TO etop;

--
-- Name: history_order_line_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_order_line_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_order_line_seq OWNER TO etop;

--
-- Name: history_order_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_order_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_order_seq OWNER TO etop;

--
-- Name: history_order_source_internal_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_order_source_internal_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_order_source_internal_seq OWNER TO etop;

--
-- Name: history_order_source_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_order_source_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_order_source_seq OWNER TO etop;

--
-- Name: history_partner_relation_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_partner_relation_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_partner_relation_seq OWNER TO etop;

--
-- Name: history_partner_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_partner_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_partner_seq OWNER TO etop;

--
-- Name: history_payment_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_payment_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_payment_seq OWNER TO etop;

--
-- Name: history_product_brand_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_brand_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_brand_seq OWNER TO etop;

--
-- Name: history_product_external_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_external_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_external_seq OWNER TO etop;

--
-- Name: history_product_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_seq OWNER TO etop;

--
-- Name: history_product_shop_collection_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_shop_collection_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_shop_collection_seq OWNER TO etop;

--
-- Name: history_product_source_category_external_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_source_category_external_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_source_category_external_seq OWNER TO etop;

--
-- Name: history_product_source_category_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_source_category_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_source_category_seq OWNER TO etop;

--
-- Name: history_product_source_internal_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_source_internal_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_source_internal_seq OWNER TO etop;

--
-- Name: history_product_source_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_product_source_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_product_source_seq OWNER TO etop;

--
-- Name: history_purchase_order_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_purchase_order_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_purchase_order_seq OWNER TO etop;

--
-- Name: history_purchase_refund_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_purchase_refund_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_purchase_refund_seq OWNER TO etop;

--
-- Name: history_receipt_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_receipt_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_receipt_seq OWNER TO etop;

--
-- Name: history_refund_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_refund_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_refund_seq OWNER TO etop;

--
-- Name: history_shipnow_fulfillment_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shipnow_fulfillment_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shipnow_fulfillment_seq OWNER TO etop;

--
-- Name: history_shop_brand_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_brand_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_brand_seq OWNER TO etop;

--
-- Name: history_shop_carrier_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_carrier_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_carrier_seq OWNER TO etop;

--
-- Name: history_shop_collection_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_collection_seq
  START WITH 241
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_collection_seq OWNER TO etop;

--
-- Name: history_shop_connection_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_connection_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_connection_seq OWNER TO etop;

--
-- Name: history_shop_customer_group_customer_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_customer_group_customer_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_customer_group_customer_seq OWNER TO etop;

--
-- Name: history_shop_customer_group_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_customer_group_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_customer_group_seq OWNER TO etop;

--
-- Name: history_shop_customer_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_customer_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_customer_seq OWNER TO etop;

--
-- Name: history_shop_ledger_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_ledger_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_ledger_seq OWNER TO etop;

--
-- Name: history_shop_product_collection_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_product_collection_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_product_collection_seq OWNER TO etop;

--
-- Name: history_shop_product_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_product_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_product_seq OWNER TO etop;

--
-- Name: history_shop_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_seq
  START WITH 318
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_seq OWNER TO etop;

--
-- Name: history_shop_stocktake_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_stocktake_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_stocktake_seq OWNER TO etop;

--
-- Name: history_shop_supplier_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_supplier_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_supplier_seq OWNER TO etop;

--
-- Name: history_shop_trader_address_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_trader_address_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_trader_address_seq OWNER TO etop;

--
-- Name: history_shop_trader_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_trader_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_trader_seq OWNER TO etop;

--
-- Name: history_shop_variant_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_variant_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_variant_seq OWNER TO etop;

--
-- Name: history_shop_variant_supplier_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_variant_supplier_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_variant_supplier_seq OWNER TO etop;

--
-- Name: history_shop_vendor_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_shop_vendor_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_shop_vendor_seq OWNER TO etop;

--
-- Name: history_supplier_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_supplier_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_supplier_seq OWNER TO etop;

--
-- Name: history_transaction_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_transaction_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_transaction_seq OWNER TO etop;

--
-- Name: history_user_internal_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_user_internal_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_user_internal_seq OWNER TO etop;

--
-- Name: history_user_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_user_seq
  START WITH 378
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_user_seq OWNER TO etop;

--
-- Name: history_variant_external_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_variant_external_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_variant_external_seq OWNER TO etop;

--
-- Name: history_variant_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_variant_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_variant_seq OWNER TO etop;

--
-- Name: history_webhook_seq; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.history_webhook_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.history_webhook_seq OWNER TO etop;

--
-- Name: ids; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.ids
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.ids OWNER TO etop;

--
-- Name: import_attempt; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.import_attempt (
                                     id bigint NOT NULL,
                                     user_id bigint NOT NULL,
                                     account_id bigint NOT NULL,
                                     original_file text NOT NULL,
                                     stored_file text,
                                     type public.import_type NOT NULL,
                                     n_created smallint NOT NULL,
                                     n_updated smallint NOT NULL,
                                     n_error smallint NOT NULL,
                                     status smallint NOT NULL,
                                     error_type text,
                                     errors jsonb,
                                     duration_ms integer NOT NULL,
                                     created_at timestamp with time zone NOT NULL,
                                     rid bigint NOT NULL
);


ALTER TABLE public.import_attempt OWNER TO etop;

--
-- Name: inventory_variant; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.inventory_variant (
                                        shop_id bigint,
                                        variant_id bigint,
                                        quantity_on_hand integer,
                                        quantity_picked integer,
                                        cost_price integer,
                                        created_at timestamp with time zone,
                                        updated_at timestamp with time zone,
                                        rid bigint NOT NULL
);


ALTER TABLE public.inventory_variant OWNER TO etop;

--
-- Name: inventory_voucher; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.inventory_voucher (
                                        title character varying,
                                        shop_id bigint,
                                        id bigint,
                                        created_by bigint,
                                        updated_by bigint,
                                        status integer,
                                        trader_id bigint,
                                        total_amount integer,
                                        type public.inventory_voucher_type,
                                        created_at timestamp with time zone,
                                        updated_at timestamp with time zone,
                                        confirmed_at timestamp with time zone,
                                        cancelled_at timestamp with time zone,
                                        cancel_reason character varying,
                                        lines jsonb,
                                        ref_id bigint,
                                        ref_type text,
                                        code_norm integer,
                                        code text,
                                        trader jsonb,
                                        variant_ids bigint[],
                                        ref_code text,
                                        product_ids bigint[],
                                        rid bigint NOT NULL
);


ALTER TABLE public.inventory_voucher OWNER TO etop;

--
-- Name: invitation; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.invitation (
                                 id bigint NOT NULL,
                                 account_id bigint NOT NULL,
                                 email text,
                                 roles text[],
                                 token text,
                                 status smallint,
                                 invited_by bigint NOT NULL,
                                 accepted_at timestamp with time zone,
                                 rejected_at timestamp with time zone,
                                 expires_at timestamp with time zone,
                                 created_at timestamp with time zone,
                                 updated_at timestamp with time zone,
                                 deleted_at timestamp with time zone,
                                 rid bigint NOT NULL,
                                 full_name text,
                                 short_name text,
                                 "position" text,
                                 phone text
);


ALTER TABLE public.invitation OWNER TO etop;

--
-- Name: money_transaction_shipping; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.money_transaction_shipping (
                                                 id bigint NOT NULL,
                                                 shop_id bigint,
                                                 supplier_id bigint,
                                                 status smallint,
                                                 total_cod integer NOT NULL,
                                                 total_orders integer NOT NULL,
                                                 code text NOT NULL,
                                                 created_at timestamp with time zone,
                                                 updated_at timestamp with time zone,
                                                 closed_at timestamp with time zone,
                                                 money_transaction_shipping_external_id bigint,
                                                 provider public.shipping_provider,
                                                 etop_transfered_at timestamp with time zone,
                                                 total_amount integer,
                                                 rid bigint NOT NULL,
                                                 confirmed_at timestamp with time zone,
                                                 money_transaction_shipping_etop_id bigint,
                                                 bank_account json,
                                                 note text,
                                                 invoice_number text,
                                                 type text
);


ALTER TABLE public.money_transaction_shipping OWNER TO etop;

--
-- Name: money_transaction_shipping_etop; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.money_transaction_shipping_etop (
                                                      id bigint NOT NULL,
                                                      code text NOT NULL,
                                                      total_cod integer NOT NULL,
                                                      total_orders integer NOT NULL,
                                                      total_amount integer NOT NULL,
                                                      total_fee integer NOT NULL,
                                                      total_money_transaction integer NOT NULL,
                                                      created_at timestamp with time zone,
                                                      updated_at timestamp with time zone,
                                                      confirmed_at timestamp with time zone,
                                                      status smallint,
                                                      bank_account jsonb,
                                                      note text,
                                                      invoice_number text
);


ALTER TABLE public.money_transaction_shipping_etop OWNER TO etop;

--
-- Name: money_transaction_shipping_external; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.money_transaction_shipping_external (
                                                          id bigint NOT NULL,
                                                          code text NOT NULL,
                                                          total_cod integer NOT NULL,
                                                          total_orders integer NOT NULL,
                                                          created_at timestamp with time zone,
                                                          updated_at timestamp with time zone,
                                                          status smallint,
                                                          external_paid_at timestamp with time zone,
                                                          provider public.shipping_provider NOT NULL,
                                                          bank_account json,
                                                          note text,
                                                          invoice_number text
);


ALTER TABLE public.money_transaction_shipping_external OWNER TO etop;

--
-- Name: money_transaction_shipping_external_line; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.money_transaction_shipping_external_line (
                                                               id bigint NOT NULL,
                                                               external_code text NOT NULL,
                                                               external_total_cod integer NOT NULL,
                                                               external_created_at timestamp with time zone,
                                                               external_closed_at timestamp with time zone,
                                                               external_customer text,
                                                               external_address text,
                                                               etop_fulfillment_id_raw text,
                                                               etop_fulfillment_id bigint,
                                                               note text,
                                                               money_transaction_shipping_external_id bigint NOT NULL,
                                                               created_at timestamp with time zone,
                                                               updated_at timestamp with time zone,
                                                               import_error jsonb,
                                                               external_total_shipping_fee integer
);


ALTER TABLE public.money_transaction_shipping_external_line OWNER TO etop;

--
-- Name: order; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public."order" (
                              id bigint NOT NULL,
                              rid bigint NOT NULL,
                              shop_id bigint NOT NULL,
                              code text NOT NULL,
                              product_ids bigint[],
                              variant_ids bigint[],
                              supplier_ids bigint[],
                              currency text,
                              payment_method text,
                              customer jsonb,
                              customer_address jsonb,
                              billing_address jsonb,
                              shipping_address jsonb,
                              customer_phone text,
                              customer_email text,
                              created_at timestamp with time zone NOT NULL,
                              processed_at timestamp with time zone,
                              updated_at timestamp with time zone NOT NULL,
                              closed_at timestamp with time zone,
                              confirmed_at timestamp with time zone,
                              cancelled_at timestamp with time zone,
                              cancel_reason text,
                              customer_confirm smallint,
                              external_confirm smallint,
                              shop_confirm smallint NOT NULL,
                              confirm_status smallint,
                              processing_status public.processing_status,
                              status smallint NOT NULL,
                              lines jsonb,
                              discounts jsonb,
                              total_items integer,
                              basket_value integer,
                              total_weight integer,
                              total_tax integer,
                              total_discount integer,
                              total_amount integer,
                              order_note text,
                              shop_note text,
                              shipping_note text,
                              order_source_id bigint,
                              order_source_type public.order_source_type NOT NULL,
                              external_order_id text,
                              fulfillment_shipping_status smallint NOT NULL,
                              customer_payment_status smallint,
                              shop_shipping_fee integer,
                              reference_url text,
                              external_order_source text,
                              shop_address jsonb,
                              shop_shipping jsonb,
                              shop_cod integer,
                              is_outside_etop boolean DEFAULT false,
                              ghn_note_code public.ghn_note_code,
                              ed_code text,
                              order_discount integer,
                              etop_payment_status smallint NOT NULL,
                              fulfillment_shipping_states text[],
                              fulfillment_payment_statuses smallint[],
                              fulfillment_shipping_codes text[],
                              customer_name text,
                              customer_name_norm tsvector,
                              product_name_norm tsvector,
                              fulfillment_sync_statuses smallint[],
                              partner_id bigint,
                              try_on public.try_on,
                              total_fee integer,
                              fee_lines jsonb,
                              external_url text,
                              fulfillment_type smallint,
                              fulfillment_ids bigint[],
                              external_meta jsonb,
                              trading_shop_id bigint,
                              payment_status smallint,
                              payment_id bigint,
                              referral_meta jsonb,
                              customer_id bigint,
                              created_by bigint,
                              fulfillment_statuses smallint[],
                              CONSTRAINT customer_email_notempty CHECK (((customer_email IS NULL) OR (customer_email <> ''::text))),
                              CONSTRAINT payment_price CHECK (((COALESCE(order_discount, 0) >= 0) AND (COALESCE(order_discount, 0) <= total_discount))),
                              CONSTRAINT total_amount CHECK ((total_amount = ((basket_value - total_discount) + COALESCE(total_fee, shop_shipping_fee))))
);


ALTER TABLE public."order" OWNER TO etop;

--
-- Name: order_line; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.order_line (
                                 rid bigint NOT NULL,
                                 order_id bigint NOT NULL,
                                 product_id bigint,
                                 variant_id bigint,
                                 supplier_id bigint,
                                 external_variant_id text,
                                 external_supplier_order_id text,
                                 product_name text,
                                 supplier_name text,
                                 image_url text,
                                 created_at timestamp with time zone,
                                 updated_at timestamp with time zone,
                                 closed_at timestamp with time zone,
                                 confirmed_at timestamp with time zone,
                                 cancelled_at timestamp with time zone,
                                 cancel_reason text,
                                 supplier_confirm smallint,
                                 status smallint,
                                 weight integer NOT NULL,
                                 quantity integer NOT NULL,
                                 wholesale_price_0 integer,
                                 wholesale_price integer,
                                 list_price integer NOT NULL,
                                 retail_price integer NOT NULL,
                                 payment_price integer NOT NULL,
                                 line_amount integer NOT NULL,
                                 total_discount integer NOT NULL,
                                 total_line_amount integer NOT NULL,
                                 requires_shipping boolean,
                                 is_outside_etop boolean DEFAULT false NOT NULL,
                                 code text NOT NULL,
                                 shop_id bigint NOT NULL,
                                 is_free boolean,
                                 meta_fields jsonb,
                                 CONSTRAINT payment_price CHECK (((payment_price >= 0) AND (payment_price <= retail_price))),
                                 CONSTRAINT quantity CHECK ((quantity > 0)),
                                 CONSTRAINT total_discount CHECK (((total_discount >= 0) AND (total_discount = ((retail_price - payment_price) * quantity))))
);


ALTER TABLE public.order_line OWNER TO etop;

--
-- Name: order_source_internal; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.order_source_internal (
                                            id bigint NOT NULL,
                                            rid bigint NOT NULL,
                                            secret jsonb,
                                            access_token text,
                                            expires_at timestamp with time zone,
                                            created_at timestamp with time zone,
                                            updated_at timestamp with time zone
);


ALTER TABLE public.order_source_internal OWNER TO etop;

--
-- Name: partner; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.partner (
                              id bigint NOT NULL,
                              rid bigint NOT NULL,
                              name text NOT NULL,
                              public_name text NOT NULL,
                              owner_id bigint NOT NULL,
                              status smallint NOT NULL,
                              is_test smallint NOT NULL,
                              phone text,
                              email text,
                              website_url text,
                              image_url text,
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              deleted_at timestamp with time zone,
                              contact_persons jsonb,
                              recognized_hosts text[],
                              redirect_urls text[],
                              available_from_etop boolean,
                              available_from_etop_config jsonb,
                              white_label_key text
);


ALTER TABLE public.partner OWNER TO etop;

--
-- Name: partner_relation; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.partner_relation (
                                       partner_id bigint NOT NULL,
                                       subject_id bigint NOT NULL,
                                       subject_type public.subject_type NOT NULL,
                                       external_subject_id text,
                                       nonce bigint,
                                       status smallint NOT NULL,
                                       roles text[],
                                       permissions text[],
                                       created_at timestamp with time zone,
                                       updated_at timestamp with time zone,
                                       deleted_at timestamp with time zone,
                                       rid bigint NOT NULL,
                                       auth_key text NOT NULL
);


ALTER TABLE public.partner_relation OWNER TO etop;

--
-- Name: partner_relation_view; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.partner_relation_view AS
SELECT partner_relation.partner_id,
       partner_relation.subject_id,
       partner_relation.subject_type,
       partner_relation.status,
       partner_relation.created_at,
       partner_relation.updated_at,
       partner_relation.deleted_at
FROM public.partner_relation;


ALTER TABLE public.partner_relation_view OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.payment (
                              id bigint NOT NULL,
                              amount integer,
                              status smallint,
                              state text,
                              payment_provider text,
                              external_trans_id text,
                              external_data jsonb,
                              created_at timestamp without time zone,
                              updated_at timestamp without time zone,
                              rid bigint NOT NULL
);


ALTER TABLE public.payment OWNER TO etop;

--
-- Name: product_shop_collection; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.product_shop_collection (
                                              product_id bigint NOT NULL,
                                              shop_id bigint,
                                              collection_id bigint NOT NULL,
                                              created_at timestamp with time zone,
                                              updated_at timestamp with time zone,
                                              status smallint,
                                              rid bigint NOT NULL
);


ALTER TABLE public.product_shop_collection OWNER TO etop;

--
-- Name: province; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.province (
                               code text NOT NULL,
                               region public.province_region NOT NULL,
                               name text NOT NULL,
                               vtpost_id integer
);


ALTER TABLE public.province OWNER TO etop;

--
-- Name: purchase_order; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.purchase_order (
                                     id bigint NOT NULL,
                                     shop_id bigint NOT NULL,
                                     supplier_id bigint NOT NULL,
                                     supplier jsonb,
                                     basket_value bigint,
                                     total_discount bigint,
                                     total_amount bigint,
                                     paid_amount bigint,
                                     code text NOT NULL,
                                     code_norm integer,
                                     note text,
                                     status smallint,
                                     variant_ids bigint[],
                                     lines jsonb,
                                     cancelled_reason text,
                                     created_by bigint NOT NULL,
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     confirmed_at timestamp with time zone,
                                     cancelled_at timestamp with time zone,
                                     deleted_at timestamp with time zone,
                                     rid bigint NOT NULL,
                                     supplier_full_name_norm text,
                                     supplier_phone_norm text,
                                     discount_lines jsonb,
                                     total_fee integer,
                                     fee_lines jsonb
);


ALTER TABLE public.purchase_order OWNER TO etop;

--
-- Name: purchase_refund; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.purchase_refund (
                                      id bigint NOT NULL,
                                      shop_id bigint NOT NULL,
                                      purchase_order_id bigint NOT NULL,
                                      note text,
                                      code_norm integer,
                                      code text,
                                      supplier_id bigint,
                                      lines jsonb,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      cancelled_at timestamp with time zone,
                                      confirmed_at timestamp with time zone,
                                      created_by bigint NOT NULL,
                                      updated_by bigint NOT NULL,
                                      total_amount integer,
                                      basket_value integer,
                                      cancel_reason text,
                                      status integer,
                                      adjustment_lines jsonb,
                                      total_adjustment integer,
                                      rid bigint NOT NULL
);


ALTER TABLE public.purchase_refund OWNER TO etop;

--
-- Name: receipt; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.receipt (
                              id bigint NOT NULL,
                              shop_id bigint NOT NULL,
                              trader_id bigint,
                              created_by bigint,
                              code text NOT NULL,
                              title text NOT NULL,
                              description text,
                              amount integer,
                              status smallint,
                              type public.receipt_type NOT NULL,
                              lines jsonb,
                              ref_ids bigint[],
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              deleted_at timestamp with time zone,
                              rid bigint NOT NULL,
                              shop_ledger_id bigint,
                              created_type public.receipt_created_type,
                              ledger_id bigint,
                              cancelled_reason text,
                              code_norm integer,
                              paid_at timestamp with time zone,
                              confirmed_at timestamp with time zone,
                              cancelled_at timestamp with time zone,
                              ref_type public.receipt_ref_type,
                              trader jsonb,
                              trader_full_name_norm tsvector,
                              trader_type public.trader_type,
                              trader_phone_norm tsvector
);


ALTER TABLE public.receipt OWNER TO etop;

--
-- Name: refund; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.refund (
                             id bigint NOT NULL,
                             shop_id bigint NOT NULL,
                             order_id bigint NOT NULL,
                             note text,
                             code_norm integer,
                             code text,
                             customer_id bigint,
                             lines jsonb,
                             created_at timestamp with time zone,
                             updated_at timestamp with time zone,
                             cancelled_at timestamp with time zone,
                             confirmed_at timestamp with time zone,
                             created_by bigint NOT NULL,
                             updated_by bigint NOT NULL,
                             total_amount integer,
                             basket_value integer,
                             cancel_reason text,
                             status integer,
                             adjustment_lines jsonb,
                             total_adjustment integer,
                             rid bigint NOT NULL
);


ALTER TABLE public.refund OWNER TO etop;

--
-- Name: shipnow_fulfillment; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shipnow_fulfillment (
                                          id bigint NOT NULL,
                                          shop_id bigint,
                                          partner_id bigint,
                                          order_ids bigint[],
                                          pickup_address jsonb,
                                          carrier text,
                                          shipping_service_code text,
                                          shipping_service_fee integer,
                                          chargeable_weight integer,
                                          gross_weight integer,
                                          basket_value integer,
                                          cod_amount integer,
                                          shipping_note text,
                                          request_pickup_at timestamp with time zone,
                                          delivery_points jsonb,
                                          status integer,
                                          shipping_state text,
                                          sync_status integer,
                                          sync_states jsonb,
                                          last_sync_at timestamp with time zone,
                                          created_at timestamp with time zone,
                                          updated_at timestamp with time zone,
                                          confirm_status integer,
                                          shipping_status integer,
                                          shipping_code text,
                                          fee_lines jsonb,
                                          carrier_fee_lines jsonb,
                                          total_fee integer,
                                          shipping_created_at timestamp with time zone,
                                          etop_payment_status integer,
                                          cod_etop_transfered_at timestamp with time zone,
                                          shipping_picking_at timestamp with time zone,
                                          shipping_delivering_at timestamp with time zone,
                                          shipping_delivered_at timestamp with time zone,
                                          shipping_cancelled_at timestamp with time zone,
                                          shipping_service_name text,
                                          shipping_service_description text,
                                          cancel_reason text,
                                          shipping_shared_link text,
                                          address_to_province_code text,
                                          address_to_district_code text,
                                          rid bigint NOT NULL
);


ALTER TABLE public.shipnow_fulfillment OWNER TO etop;

--
-- Name: shipping_code; Type: SEQUENCE; Schema: public; Owner: etop
--

CREATE SEQUENCE public.shipping_code
  START WITH 100001
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE public.shipping_code OWNER TO etop;

--
-- Name: shipping_source; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shipping_source (
                                      id bigint NOT NULL,
                                      name text,
                                      type public.shipping_provider NOT NULL,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      rid bigint,
                                      username text
);


ALTER TABLE public.shipping_source OWNER TO etop;

--
-- Name: shipping_source_internal; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shipping_source_internal (
                                               id bigint NOT NULL,
                                               rid bigint,
                                               last_sync_at timestamp with time zone,
                                               created_at timestamp with time zone,
                                               updated_at timestamp with time zone,
                                               secret jsonb,
                                               access_token text,
                                               expires_at timestamp with time zone
);


ALTER TABLE public.shipping_source_internal OWNER TO etop;

--
-- Name: shop; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop (
                           id bigint NOT NULL,
                           rid bigint NOT NULL,
                           name text NOT NULL,
                           owner_id bigint NOT NULL,
                           status smallint NOT NULL,
                           product_source_id bigint,
                           created_at timestamp with time zone DEFAULT date_trunc('second'::text, now()) NOT NULL,
                           updated_at timestamp with time zone NOT NULL,
                           rules jsonb,
                           is_test smallint DEFAULT '0'::smallint NOT NULL,
                           image_url text,
                           phone text,
                           website_url text,
                           email text,
                           deleted_at timestamp with time zone,
                           address_id bigint,
                           bank_account jsonb,
                           contact_persons jsonb,
                           order_source_id bigint,
                           ship_to_address_id bigint,
                           ship_from_address_id bigint,
                           code text,
                           auto_create_ffm boolean DEFAULT false,
                           ghn_note_code public.ghn_note_code,
                           try_on public.try_on NOT NULL,
                           recognized_hosts text[],
                           company_info jsonb,
                           money_transaction_rrule text,
                           survey_info jsonb,
                           shipping_service_select_strategy jsonb,
                           inventory_overstock boolean,
                           wl_partner_id bigint
);


ALTER TABLE public.shop OWNER TO etop;

--
-- Name: shop_brand; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_brand (
                                 shop_id bigint,
                                 id bigint,
                                 description text,
                                 brand_name text,
                                 updated_at timestamp with time zone,
                                 created_at timestamp with time zone,
                                 deleted_at timestamp with time zone,
                                 rid bigint NOT NULL,
                                 external_id text,
                                 partner_id bigint
);


ALTER TABLE public.shop_brand OWNER TO etop;

--
-- Name: shop_carrier; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_carrier (
                                   id bigint NOT NULL,
                                   shop_id bigint,
                                   full_name text,
                                   note text,
                                   status smallint NOT NULL,
                                   created_at timestamp with time zone,
                                   updated_at timestamp with time zone,
                                   deleted_at timestamp with time zone,
                                   rid bigint NOT NULL
);


ALTER TABLE public.shop_carrier OWNER TO etop;

--
-- Name: shop_category; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_category (
                                    id bigint NOT NULL,
                                    rid bigint NOT NULL,
                                    product_source_id bigint,
                                    product_source_type public.product_source_type,
                                    supplier_id bigint,
                                    parent_id bigint,
                                    name text,
                                    status smallint NOT NULL,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    deleted_at timestamp with time zone,
                                    shop_id bigint,
                                    external_id text,
                                    external_parent_id text,
                                    partner_id bigint
);


ALTER TABLE public.shop_category OWNER TO etop;

--
-- Name: shop_collection; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_collection (
                                      id bigint NOT NULL,
                                      rid bigint NOT NULL,
                                      shop_id bigint NOT NULL,
                                      name text,
                                      description text,
                                      desc_html text,
                                      short_desc text,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      deleted_at timestamp with time zone,
                                      partner_id bigint,
                                      external_id text
);


ALTER TABLE public.shop_collection OWNER TO etop;

--
-- Name: shop_connection; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_connection (
                                      shop_id bigint,
                                      connection_id bigint,
                                      token text,
                                      token_expires_at timestamp with time zone,
                                      status smallint,
                                      connection_states jsonb,
                                      is_global boolean,
                                      created_at timestamp with time zone,
                                      updated_at timestamp with time zone,
                                      deleted_at timestamp with time zone,
                                      external_data jsonb,
                                      rid bigint NOT NULL
);


ALTER TABLE public.shop_connection OWNER TO etop;

--
-- Name: shop_customer; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_customer (
                                    id bigint NOT NULL,
                                    shop_id bigint,
                                    code text,
                                    full_name text,
                                    gender public.gender_type,
                                    type public.customer_type,
                                    birthday date,
                                    note text,
                                    phone text,
                                    email text,
                                    status smallint NOT NULL,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    deleted_at timestamp with time zone,
                                    rid bigint NOT NULL,
                                    full_name_norm tsvector,
                                    phone_norm tsvector,
                                    code_norm integer,
                                    external_id text,
                                    external_code text,
                                    partner_id bigint
);


ALTER TABLE public.shop_customer OWNER TO etop;

--
-- Name: shop_customer_group; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_customer_group (
                                          id bigint NOT NULL,
                                          name text,
                                          created_at timestamp with time zone,
                                          updated_at timestamp with time zone,
                                          deleted_at timestamp with time zone,
                                          shop_id bigint,
                                          rid bigint NOT NULL,
                                          partner_id bigint
);


ALTER TABLE public.shop_customer_group OWNER TO etop;

--
-- Name: shop_customer_group_customer; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_customer_group_customer (
                                                   customer_id bigint NOT NULL,
                                                   group_id bigint NOT NULL,
                                                   created_at timestamp with time zone,
                                                   updated_at timestamp with time zone,
                                                   rid bigint NOT NULL
);


ALTER TABLE public.shop_customer_group_customer OWNER TO etop;

--
-- Name: shop_ledger; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_ledger (
                                  id bigint NOT NULL,
                                  shop_id bigint,
                                  name text NOT NULL,
                                  bank_account jsonb,
                                  status smallint,
                                  note text,
                                  type public.shop_ledger_type NOT NULL,
                                  created_by bigint,
                                  created_at timestamp with time zone,
                                  updated_at timestamp with time zone,
                                  deleted_at timestamp with time zone,
                                  rid bigint NOT NULL
);


ALTER TABLE public.shop_ledger OWNER TO etop;

--
-- Name: shop_product; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_product (
                                   shop_id bigint NOT NULL,
                                   product_id bigint NOT NULL,
                                   rid bigint NOT NULL,
                                   collection_id bigint,
                                   name text,
                                   description text,
                                   desc_html text,
                                   short_desc text,
                                   retail_price integer,
                                   tags text[],
                                   created_at timestamp with time zone NOT NULL,
                                   updated_at timestamp with time zone NOT NULL,
                                   note text,
                                   image_urls text[],
                                   status smallint,
                                   haravan_id text,
                                   name_norm tsvector,
                                   deleted_at timestamp with time zone,
                                   code text,
                                   name_norm_ua text,
                                   category_id bigint,
                                   cost_price integer,
                                   list_price integer,
                                   unit text,
                                   product_type text,
                                   meta_fields jsonb,
                                   brand_id bigint,
                                   external_id text,
                                   external_code text,
                                   partner_id bigint,
                                   code_norm integer,
                                   external_brand_id text,
                                   external_category_id text
);


ALTER TABLE public.shop_product OWNER TO etop;

--
-- Name: shop_product_collection; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_product_collection (
                                              product_id bigint NOT NULL,
                                              collection_id bigint NOT NULL,
                                              shop_id bigint,
                                              status smallint,
                                              created_at timestamp with time zone,
                                              updated_at timestamp with time zone,
                                              partner_id bigint,
                                              external_collection_id text,
                                              external_product_id text,
                                              deleted_at timestamp with time zone,
                                              rid bigint NOT NULL
);


ALTER TABLE public.shop_product_collection OWNER TO etop;

--
-- Name: shop_stocktake; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_stocktake (
                                     id bigint NOT NULL,
                                     shop_id bigint,
                                     created_by bigint,
                                     updated_by bigint,
                                     created_at timestamp with time zone,
                                     updated_at timestamp with time zone,
                                     confirmed_at timestamp with time zone,
                                     cancelled_at timestamp with time zone,
                                     variant_ids bigint[],
                                     total_quantity integer,
                                     status integer,
                                     lines jsonb,
                                     code text,
                                     code_norm integer,
                                     note text,
                                     cancel_reason text,
                                     product_ids bigint[],
                                     type text,
                                     rid bigint NOT NULL
);


ALTER TABLE public.shop_stocktake OWNER TO etop;

--
-- Name: shop_supplier; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_supplier (
                                    id bigint NOT NULL,
                                    shop_id bigint NOT NULL,
                                    full_name text,
                                    note text,
                                    status smallint,
                                    created_at timestamp with time zone,
                                    updated_at timestamp with time zone,
                                    deleted_at timestamp with time zone,
                                    rid bigint NOT NULL,
                                    phone text,
                                    email text,
                                    company_name text,
                                    tax_number text,
                                    headquater_address text,
                                    full_name_norm tsvector,
                                    phone_norm tsvector,
                                    code_norm integer,
                                    code text,
                                    company_name_norm tsvector
);


ALTER TABLE public.shop_supplier OWNER TO etop;

--
-- Name: shop_trader; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_trader (
                                  id bigint NOT NULL,
                                  shop_id bigint,
                                  rid bigint NOT NULL,
                                  type public.trader_type NOT NULL,
                                  deleted_at timestamp with time zone
);


ALTER TABLE public.shop_trader OWNER TO etop;

--
-- Name: shop_trader_address; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_trader_address (
                                          id bigint NOT NULL,
                                          shop_id bigint NOT NULL,
                                          trader_id bigint NOT NULL,
                                          full_name text,
                                          phone text,
                                          email text,
                                          company text,
                                          district_code text,
                                          ward_code text,
                                          city text,
                                          address1 text,
                                          address2 text,
                                          "position" text,
                                          note text,
                                          "primary" boolean,
                                          status smallint NOT NULL,
                                          coordinates jsonb,
                                          created_at timestamp with time zone,
                                          updated_at timestamp with time zone,
                                          deleted_at timestamp with time zone,
                                          rid bigint NOT NULL,
                                          is_default boolean,
                                          partner_id bigint
);


ALTER TABLE public.shop_trader_address OWNER TO etop;

--
-- Name: shop_variant; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_variant (
                                   shop_id bigint NOT NULL,
                                   variant_id bigint NOT NULL,
                                   rid bigint NOT NULL,
                                   collection_id bigint,
                                   name text,
                                   description text,
                                   desc_html text,
                                   short_desc text,
                                   retail_price integer,
                                   tags text[],
                                   created_at timestamp with time zone NOT NULL,
                                   updated_at timestamp with time zone NOT NULL,
                                   note text,
                                   image_urls text[],
                                   status smallint,
                                   haravan_id text,
                                   name_norm tsvector,
                                   product_id bigint,
                                   deleted_at timestamp with time zone,
                                   attr_norm_kv tsvector,
                                   code text,
                                   cost_price integer,
                                   list_price integer,
                                   attributes jsonb,
                                   external_id text,
                                   external_code text,
                                   partner_id bigint,
                                   code_norm integer,
                                   external_product_id text
);


ALTER TABLE public.shop_variant OWNER TO etop;

--
-- Name: shop_variant_supplier; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.shop_variant_supplier (
                                            shop_id bigint,
                                            supplier_id bigint,
                                            variant_id bigint,
                                            created_at timestamp with time zone,
                                            updated_at timestamp with time zone,
                                            rid bigint NOT NULL
);


ALTER TABLE public.shop_variant_supplier OWNER TO etop;

--
-- Name: transaction; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.transaction (
                                  id bigint NOT NULL,
                                  amount integer,
                                  account_id bigint,
                                  status smallint,
                                  type text,
                                  note text,
                                  metadata jsonb,
                                  created_at timestamp with time zone,
                                  updated_at timestamp with time zone,
                                  rid bigint NOT NULL
);


ALTER TABLE public.transaction OWNER TO etop;

--
-- Name: user; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public."user" (
                             id bigint NOT NULL,
                             rid bigint NOT NULL,
                             status smallint NOT NULL,
                             created_at timestamp with time zone NOT NULL,
                             updated_at timestamp with time zone NOT NULL,
                             full_name text,
                             short_name text,
                             email text,
                             phone text,
                             agreed_tos_at timestamp with time zone,
                             agreed_email_info_at timestamp with time zone,
                             email_verified_at timestamp with time zone,
                             phone_verified_at timestamp with time zone,
                             email_verification_sent_at timestamp with time zone,
                             phone_verification_sent_at timestamp with time zone,
                             is_test smallint NOT NULL,
                             identifying public.user_identifying_type,
                             source text,
                             ref_user_id bigint,
                             ref_sale_id bigint,
                             wl_partner_id bigint,
                             CONSTRAINT user_identifying CHECK ((((identifying = 'full'::public.user_identifying_type) AND (status <> 0) AND (phone IS NOT NULL) AND (email IS NOT NULL)) OR ((identifying = 'half'::public.user_identifying_type) AND (status <> 0) AND (phone IS NOT NULL)) OR ((identifying = 'stub'::public.user_identifying_type) AND (status = 0) AND ((phone IS NOT NULL) OR (email IS NOT NULL))))),
                             CONSTRAINT user_not_full CHECK (((identifying = 'stub'::public.user_identifying_type) OR ((full_name IS NOT NULL) AND (short_name IS NOT NULL))))
);


ALTER TABLE public."user" OWNER TO etop;

--
-- Name: user_auth; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.user_auth (
                                user_id bigint,
                                auth_type text NOT NULL,
                                auth_key text NOT NULL,
                                created_at timestamp with time zone,
                                updated_at timestamp with time zone
);


ALTER TABLE public.user_auth OWNER TO etop;

--
-- Name: user_internal; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.user_internal (
                                    id bigint NOT NULL,
                                    hashpwd text,
                                    updated_at timestamp with time zone,
                                    rid bigint NOT NULL
);


ALTER TABLE public.user_internal OWNER TO etop;

--
-- Name: webhook; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.webhook (
                              id bigint NOT NULL,
                              account_id bigint NOT NULL,
                              entities text[],
                              fields text[],
                              url text NOT NULL,
                              metadata text,
                              created_at timestamp with time zone,
                              updated_at timestamp with time zone,
                              deleted_at timestamp with time zone,
                              rid bigint NOT NULL
);


ALTER TABLE public.webhook OWNER TO etop;

--
-- Name: webhook_changes; Type: TABLE; Schema: public; Owner: etop
--

CREATE TABLE public.webhook_changes (
                                      id bigint NOT NULL,
                                      webhook_id bigint NOT NULL,
                                      account_id bigint NOT NULL,
                                      created_at timestamp with time zone,
                                      changes jsonb,
                                      result jsonb
);


ALTER TABLE public.webhook_changes OWNER TO etop;

--
-- Name: account_auth account_auth_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.account_auth
  ADD CONSTRAINT account_auth_rid_key UNIQUE (rid);


--
-- Name: account account_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.account
  ADD CONSTRAINT account_rid_key UNIQUE (rid);


--
-- Name: account_user account_user_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.account_user
  ADD CONSTRAINT account_user_rid_key UNIQUE (rid);


--
-- Name: address address_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.address
  ADD CONSTRAINT address_rid_key UNIQUE (rid);


--
-- Name: connection connection_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.connection
  ADD CONSTRAINT connection_rid_key UNIQUE (rid);


--
-- Name: credit credit_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.credit
  ADD CONSTRAINT credit_rid_key UNIQUE (rid);


--
-- Name: etop_category etop_category_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.etop_category
  ADD CONSTRAINT etop_category_rid_key UNIQUE (rid);


--
-- Name: fulfillment fulfillment_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.fulfillment
  ADD CONSTRAINT fulfillment_rid_key UNIQUE (rid);


--
-- Name: import_attempt import_attempt_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.import_attempt
  ADD CONSTRAINT import_attempt_rid_key UNIQUE (rid);


--
-- Name: inventory_variant inventory_variant_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.inventory_variant
  ADD CONSTRAINT inventory_variant_rid_key UNIQUE (rid);


--
-- Name: inventory_voucher inventory_voucher_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.inventory_voucher
  ADD CONSTRAINT inventory_voucher_rid_key UNIQUE (rid);


--
-- Name: invitation invitation_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.invitation
  ADD CONSTRAINT invitation_rid_key UNIQUE (rid);


--
-- Name: money_transaction_shipping money_transaction_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.money_transaction_shipping
  ADD CONSTRAINT money_transaction_rid_key UNIQUE (rid);


--
-- Name: order_external order_external_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.order_external
  ADD CONSTRAINT order_external_rid_key UNIQUE (rid);


--
-- Name: order_line order_line_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.order_line
  ADD CONSTRAINT order_line_rid_key UNIQUE (rid);


--
-- Name: order order_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history."order"
  ADD CONSTRAINT order_rid_key UNIQUE (rid);


--
-- Name: order_source_internal order_source_internal_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.order_source_internal
  ADD CONSTRAINT order_source_internal_rid_key UNIQUE (rid);


--
-- Name: order_source order_source_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.order_source
  ADD CONSTRAINT order_source_rid_key UNIQUE (rid);


--
-- Name: partner_relation partner_relation_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.partner_relation
  ADD CONSTRAINT partner_relation_rid_key UNIQUE (rid);


--
-- Name: partner partner_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.partner
  ADD CONSTRAINT partner_rid_key UNIQUE (rid);


--
-- Name: payment payment_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.payment
  ADD CONSTRAINT payment_rid_key UNIQUE (rid);


--
-- Name: product_brand product_brand_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_brand
  ADD CONSTRAINT product_brand_rid_key UNIQUE (rid);


--
-- Name: product_external product_external_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_external
  ADD CONSTRAINT product_external_rid_key UNIQUE (rid);


--
-- Name: product product_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product
  ADD CONSTRAINT product_rid_key UNIQUE (rid);


--
-- Name: product_shop_collection product_shop_collection_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_shop_collection
  ADD CONSTRAINT product_shop_collection_rid_key UNIQUE (rid);


--
-- Name: product_source_category_external product_source_category_external_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_source_category_external
  ADD CONSTRAINT product_source_category_external_rid_key UNIQUE (rid);


--
-- Name: product_source_category product_source_category_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_source_category
  ADD CONSTRAINT product_source_category_rid_key UNIQUE (rid);


--
-- Name: product_source_internal product_source_internal_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_source_internal
  ADD CONSTRAINT product_source_internal_rid_key UNIQUE (rid);


--
-- Name: product_source product_source_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.product_source
  ADD CONSTRAINT product_source_rid_key UNIQUE (rid);


--
-- Name: purchase_order purchase_order_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.purchase_order
  ADD CONSTRAINT purchase_order_rid_key UNIQUE (rid);


--
-- Name: purchase_refund purchase_refund_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.purchase_refund
  ADD CONSTRAINT purchase_refund_rid_key UNIQUE (rid);


--
-- Name: receipt receipt_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.receipt
  ADD CONSTRAINT receipt_rid_key UNIQUE (rid);


--
-- Name: refund refund_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.refund
  ADD CONSTRAINT refund_rid_key UNIQUE (rid);


--
-- Name: shipnow_fulfillment shipnow_fulfillment_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_rid_key UNIQUE (rid);


--
-- Name: shop_brand shop_brand_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_brand
  ADD CONSTRAINT shop_brand_rid_key UNIQUE (rid);


--
-- Name: shop_carrier shop_carrier_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_carrier
  ADD CONSTRAINT shop_carrier_rid_key UNIQUE (rid);


--
-- Name: shop_collection shop_collection_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_collection
  ADD CONSTRAINT shop_collection_rid_key UNIQUE (rid);


--
-- Name: shop_connection shop_connection_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_connection
  ADD CONSTRAINT shop_connection_rid_key UNIQUE (rid);


--
-- Name: shop_customer_group_customer shop_customer_group_customer_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_rid_key UNIQUE (rid);


--
-- Name: shop_customer_group shop_customer_group_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_customer_group
  ADD CONSTRAINT shop_customer_group_rid_key UNIQUE (rid);


--
-- Name: shop_customer shop_customer_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_customer
  ADD CONSTRAINT shop_customer_rid_key UNIQUE (rid);


--
-- Name: shop_ledger shop_ledger_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_ledger
  ADD CONSTRAINT shop_ledger_rid_key UNIQUE (rid);


--
-- Name: shop_product_collection shop_product_collection_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_product_collection
  ADD CONSTRAINT shop_product_collection_rid_key UNIQUE (rid);


--
-- Name: shop_product shop_product_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_product
  ADD CONSTRAINT shop_product_rid_key UNIQUE (rid);


--
-- Name: shop shop_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop
  ADD CONSTRAINT shop_rid_key UNIQUE (rid);


--
-- Name: shop_stocktake shop_stocktake_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_stocktake
  ADD CONSTRAINT shop_stocktake_rid_key UNIQUE (rid);


--
-- Name: shop_supplier shop_supplier_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_supplier
  ADD CONSTRAINT shop_supplier_rid_key UNIQUE (rid);


--
-- Name: shop_trader_address shop_trader_address_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_trader_address
  ADD CONSTRAINT shop_trader_address_rid_key UNIQUE (rid);


--
-- Name: shop_trader shop_trader_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_trader
  ADD CONSTRAINT shop_trader_rid_key UNIQUE (rid);


--
-- Name: shop_variant shop_variant_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_variant
  ADD CONSTRAINT shop_variant_rid_key UNIQUE (rid);


--
-- Name: shop_variant_supplier shop_variant_supplier_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_rid_key UNIQUE (rid);


--
-- Name: shop_vendor shop_vendor_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.shop_vendor
  ADD CONSTRAINT shop_vendor_rid_key UNIQUE (rid);


--
-- Name: supplier supplier_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.supplier
  ADD CONSTRAINT supplier_rid_key UNIQUE (rid);


--
-- Name: transaction transaction_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.transaction
  ADD CONSTRAINT transaction_rid_key UNIQUE (rid);


--
-- Name: user_internal user_internal_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.user_internal
  ADD CONSTRAINT user_internal_rid_key UNIQUE (rid);


--
-- Name: user user_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history."user"
  ADD CONSTRAINT user_rid_key UNIQUE (rid);


--
-- Name: variant_external variant_external_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.variant_external
  ADD CONSTRAINT variant_external_rid_key UNIQUE (rid);


--
-- Name: variant variant_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.variant
  ADD CONSTRAINT variant_rid_key UNIQUE (rid);


--
-- Name: webhook webhook_rid_key; Type: CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history.webhook
  ADD CONSTRAINT webhook_rid_key UNIQUE (rid);


--
-- Name: account_auth account_auth_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_auth
  ADD CONSTRAINT account_auth_pkey PRIMARY KEY (auth_key);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account
  ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.address
  ADD CONSTRAINT address_pkey PRIMARY KEY (id);


--
-- Name: affiliate affiliate_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_pkey PRIMARY KEY (id);


--
-- Name: code code_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.code
  ADD CONSTRAINT code_pkey PRIMARY KEY (code, type);


--
-- Name: connection connection_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.connection
  ADD CONSTRAINT connection_pkey PRIMARY KEY (id);


--
-- Name: credit credit_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.credit
  ADD CONSTRAINT credit_pkey PRIMARY KEY (id);


--
-- Name: district district_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.district
  ADD CONSTRAINT district_pkey PRIMARY KEY (code);


--
-- Name: export_attempt export_attempt_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_pkey PRIMARY KEY (id);


--
-- Name: external_account_ahamove external_account_ahamove_phone_key; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_phone_key UNIQUE (phone);


--
-- Name: external_account_ahamove external_account_ahamove_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_pkey PRIMARY KEY (id);


--
-- Name: external_account_haravan external_account_haravan_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.external_account_haravan
  ADD CONSTRAINT external_account_haravan_pkey PRIMARY KEY (id);


--
-- Name: fulfillment fulfillment_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_pkey PRIMARY KEY (id);


--
-- Name: import_attempt import_attempt_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_pkey PRIMARY KEY (id);


--
-- Name: invitation invitation_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_pkey PRIMARY KEY (id);


--
-- Name: money_transaction_shipping money_transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_pkey PRIMARY KEY (id);


--
-- Name: money_transaction_shipping_etop money_transaction_shipping_etop_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping_etop
  ADD CONSTRAINT money_transaction_shipping_etop_pkey PRIMARY KEY (id);


--
-- Name: money_transaction_shipping_external_line money_transaction_shipping_external_line_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping_external_line
  ADD CONSTRAINT money_transaction_shipping_external_line_pkey PRIMARY KEY (id);


--
-- Name: money_transaction_shipping_external money_transaction_shipping_external_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping_external
  ADD CONSTRAINT money_transaction_shipping_external_pkey PRIMARY KEY (id);


--
-- Name: order order_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_pkey PRIMARY KEY (id);


--
-- Name: order_source_internal order_source_internal_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.order_source_internal
  ADD CONSTRAINT order_source_internal_pkey PRIMARY KEY (id);


--
-- Name: partner partner_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_pkey PRIMARY KEY (id);


--
-- Name: partner_relation partner_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.partner_relation
  ADD CONSTRAINT partner_relation_pkey PRIMARY KEY (partner_id, subject_id, subject_type);


--
-- Name: payment payment_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.payment
  ADD CONSTRAINT payment_pkey PRIMARY KEY (id);


--
-- Name: product_shop_collection product_shop_collection_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_pkey PRIMARY KEY (product_id, collection_id);


--
-- Name: shop_category product_source_category_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT product_source_category_pkey PRIMARY KEY (id);


--
-- Name: province province_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.province
  ADD CONSTRAINT province_pkey PRIMARY KEY (code);


--
-- Name: purchase_order purchase_order_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_pkey PRIMARY KEY (id);


--
-- Name: purchase_refund purchase_refund_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.purchase_refund
  ADD CONSTRAINT purchase_refund_pkey PRIMARY KEY (id);


--
-- Name: receipt receipt_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_pkey PRIMARY KEY (id);


--
-- Name: refund refund_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.refund
  ADD CONSTRAINT refund_pkey PRIMARY KEY (id);


--
-- Name: shipnow_fulfillment shipnow_fulfillment_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_pkey PRIMARY KEY (id);


--
-- Name: shipping_source shipping_source_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shipping_source
  ADD CONSTRAINT shipping_source_pkey PRIMARY KEY (id);


--
-- Name: shipping_source_internal shipping_source_state_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shipping_source_internal
  ADD CONSTRAINT shipping_source_state_pkey PRIMARY KEY (id);


--
-- Name: shop_carrier shop_carrier_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_pkey PRIMARY KEY (id);


--
-- Name: shop shop_code_key; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_code_key UNIQUE (code);


--
-- Name: shop_collection shop_collection_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_pkey PRIMARY KEY (id);


--
-- Name: shop_customer_group_customer shop_customer_group_customer_constraint; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_constraint PRIMARY KEY (group_id, customer_id);


--
-- Name: shop_customer_group shop_customer_group_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer_group
  ADD CONSTRAINT shop_customer_group_pkey PRIMARY KEY (id);


--
-- Name: shop_customer shop_customer_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_pkey PRIMARY KEY (id);


--
-- Name: shop_ledger shop_ledger_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_pkey PRIMARY KEY (id);


--
-- Name: shop shop_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_pkey PRIMARY KEY (id);


--
-- Name: shop_product_collection shop_product_collection_constraint; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_constraint PRIMARY KEY (product_id, collection_id);


--
-- Name: shop_product shop_product_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_pkey PRIMARY KEY (product_id);


--
-- Name: shop_stocktake shop_stocktake_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_stocktake
  ADD CONSTRAINT shop_stocktake_pkey PRIMARY KEY (id);


--
-- Name: shop_trader_address shop_trader_address_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_pkey PRIMARY KEY (id);


--
-- Name: shop_trader shop_trader_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_trader
  ADD CONSTRAINT shop_trader_pkey PRIMARY KEY (id);


--
-- Name: shop_variant shop_variant_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_pkey PRIMARY KEY (variant_id);


--
-- Name: shop_supplier shop_vendor_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_pkey PRIMARY KEY (id);


--
-- Name: transaction transaction_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.transaction
  ADD CONSTRAINT transaction_pkey PRIMARY KEY (id);


--
-- Name: user_auth user_auth_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.user_auth
  ADD CONSTRAINT user_auth_pkey PRIMARY KEY (auth_type, auth_key);


--
-- Name: user_internal user_internal_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.user_internal
  ADD CONSTRAINT user_internal_pkey PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_pkey PRIMARY KEY (id);

--
-- Name: webhook_changes webhook_changes_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_pkey PRIMARY KEY (id);


--
-- Name: webhook webhook_pkey; Type: CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.webhook
  ADD CONSTRAINT webhook_pkey PRIMARY KEY (id);


--
-- Name: account_auth_key_account_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX account_auth_key_account_id_idx ON history.account_auth USING btree (auth_key, account_id);


--
-- Name: account_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX account_id_idx ON history.account USING btree (id);


--
-- Name: account_user_account_id_user_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX account_user_account_id_user_id_idx ON history.account_user USING btree (account_id, user_id);


--
-- Name: address_account_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX address_account_id_idx ON history.address USING btree (account_id);


--
-- Name: address_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX address_id_idx ON history.address USING btree (id);


--
-- Name: connection_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX connection_id_idx ON history.connection USING btree (id);


--
-- Name: credit_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX credit_id_idx ON history.credit USING btree (id);


--
-- Name: credit_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX credit_shop_id_idx ON history.credit USING btree (shop_id);


--
-- Name: etop_category_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX etop_category_id_idx ON history.etop_category USING btree (id);


--
-- Name: fulfillment_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX fulfillment_id_idx ON history.fulfillment USING btree (id);


--
-- Name: fulfillment_order_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX fulfillment_order_id_idx ON history.fulfillment USING btree (order_id);


--
-- Name: fulfillment_rid_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX fulfillment_rid_id_shop_id_idx ON history.fulfillment USING btree (rid, id, shop_id);


--
-- Name: fulfillment_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX fulfillment_shop_id_idx ON history.fulfillment USING btree (shop_id);


--
-- Name: import_attempt_id_user_id_account_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX import_attempt_id_user_id_account_id_idx ON history.import_attempt USING btree (id, user_id, account_id);


--
-- Name: inventory_variant_variant_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX inventory_variant_variant_id_shop_id_idx ON history.inventory_variant USING btree (variant_id, shop_id);


--
-- Name: inventory_voucher_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX inventory_voucher_id_shop_id_idx ON history.inventory_voucher USING btree (id, shop_id);


--
-- Name: invitation_id_account_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX invitation_id_account_id_idx ON history.invitation USING btree (id, account_id);


--
-- Name: money_transaction_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX money_transaction_id_idx ON history.money_transaction_shipping USING btree (id);


--
-- Name: money_transaction_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX money_transaction_shop_id_idx ON history.money_transaction_shipping USING btree (shop_id);


--
-- Name: order_external_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_external_id_idx ON history.order_external USING btree (id);


--
-- Name: order_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_id_idx ON history."order" USING btree (id);


--
-- Name: order_line_order_id_variant_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_line_order_id_variant_id_idx ON history.order_line USING btree (order_id, variant_id);


--
-- Name: order_line_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_line_shop_id_idx ON history.order_line USING btree (shop_id);


--
-- Name: order_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_shop_id_idx ON history."order" USING btree (shop_id);


--
-- Name: order_source_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_source_id_idx ON history.order_source USING btree (id);


--
-- Name: order_source_internal_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX order_source_internal_id_idx ON history.order_source_internal USING btree (id);


--
-- Name: partner_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX partner_id_idx ON history.partner USING btree (id);


--
-- Name: partner_relation_partner_id_subject_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX partner_relation_partner_id_subject_id_idx ON history.partner_relation USING btree (partner_id, subject_id);


--
-- Name: payment_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX payment_id_idx ON history.payment USING btree (id);


--
-- Name: product_brand_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_brand_id_idx ON history.product_brand USING btree (id);


--
-- Name: product_external_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_external_id_idx ON history.product_external USING btree (id);


--
-- Name: product_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_id_idx ON history.product USING btree (id);


--
-- Name: product_shop_collection_product_id_collection_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_shop_collection_product_id_collection_id_idx ON history.product_shop_collection USING btree (product_id, collection_id);


--
-- Name: product_source_category_external_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_source_category_external_id_idx ON history.product_source_category_external USING btree (id);


--
-- Name: product_source_category_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_source_category_id_idx ON history.product_source_category USING btree (id);


--
-- Name: product_source_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_source_id_idx ON history.product_source USING btree (id);


--
-- Name: product_source_internal_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX product_source_internal_id_idx ON history.product_source_internal USING btree (id);


--
-- Name: purchase_order_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX purchase_order_id_shop_id_idx ON history.purchase_order USING btree (id, shop_id);


--
-- Name: purchase_refund_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX purchase_refund_id_shop_id_idx ON history.purchase_refund USING btree (id, shop_id);


--
-- Name: receipt_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX receipt_id_shop_id_idx ON history.receipt USING btree (id, shop_id);


--
-- Name: refund_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX refund_id_shop_id_idx ON history.refund USING btree (id, shop_id);


--
-- Name: shipnow_fulfillment_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shipnow_fulfillment_id_shop_id_idx ON history.shipnow_fulfillment USING btree (id, shop_id);


--
-- Name: shop_brand_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_brand_id_shop_id_idx ON history.shop_brand USING btree (id, shop_id);


--
-- Name: shop_carrier_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_carrier_id_shop_id_idx ON history.shop_carrier USING btree (id, shop_id);


--
-- Name: shop_collection_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_collection_id_idx ON history.shop_collection USING btree (id);


--
-- Name: shop_connection_shop_id_connection_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_connection_shop_id_connection_id_idx ON history.shop_connection USING btree (shop_id, connection_id);


--
-- Name: shop_customer_group_customer_customer_id_group_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_customer_group_customer_customer_id_group_id_idx ON history.shop_customer_group_customer USING btree (customer_id, group_id);


--
-- Name: shop_customer_group_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_customer_group_id_shop_id_idx ON history.shop_customer_group USING btree (id, shop_id);


--
-- Name: shop_customer_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_customer_id_shop_id_idx ON history.shop_customer USING btree (id, shop_id);


--
-- Name: shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_id_idx ON history.shop USING btree (id);


--
-- Name: shop_ledger_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_ledger_id_shop_id_idx ON history.shop_ledger USING btree (id, shop_id);


--
-- Name: shop_product_collection_shop_id_product_id_collection_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_product_collection_shop_id_product_id_collection_id_idx ON history.shop_product_collection USING btree (shop_id, product_id, collection_id);


--
-- Name: shop_product_shop_id_product_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_product_shop_id_product_id_idx ON history.shop_product USING btree (shop_id, product_id);


--
-- Name: shop_stocktake_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_stocktake_id_shop_id_idx ON history.shop_stocktake USING btree (id, shop_id);


--
-- Name: shop_supplier_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_supplier_id_shop_id_idx ON history.shop_supplier USING btree (id, shop_id);


--
-- Name: shop_trader_address_id_shop_id_trader_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_trader_address_id_shop_id_trader_id_idx ON history.shop_trader_address USING btree (id, shop_id, trader_id);


--
-- Name: shop_trader_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_trader_id_shop_id_idx ON history.shop_trader USING btree (id, shop_id);


--
-- Name: shop_variant_shop_id_variant_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_variant_shop_id_variant_id_idx ON history.shop_variant USING btree (shop_id, variant_id);


--
-- Name: shop_variant_supplier_variant_id_supplier_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_variant_supplier_variant_id_supplier_id_idx ON history.shop_variant_supplier USING btree (variant_id, supplier_id);


--
-- Name: shop_vendor_id_shop_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX shop_vendor_id_shop_id_idx ON history.shop_vendor USING btree (id, shop_id);


--
-- Name: supplier_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX supplier_id_idx ON history.supplier USING btree (id);


--
-- Name: transaction_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX transaction_id_idx ON history.transaction USING btree (id);


--
-- Name: user_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX user_id_idx ON history."user" USING btree (id);


--
-- Name: user_internal_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX user_internal_id_idx ON history.user_internal USING btree (id);


--
-- Name: variant_external_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX variant_external_id_idx ON history.variant_external USING btree (id);


--
-- Name: variant_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX variant_id_idx ON history.variant USING btree (id);


--
-- Name: webhook_id_idx; Type: INDEX; Schema: history; Owner: etop
--

CREATE INDEX webhook_id_idx ON history.webhook USING btree (id);


--
-- Name: account_auth_account_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX account_auth_account_id_idx ON public.account_auth USING btree (account_id);


--
-- Name: account_user_account_id_user_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX account_user_account_id_user_id_idx ON public.account_user USING btree (account_id, user_id) WHERE (deleted_at IS NULL);


--
-- Name: connection_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX connection_code_idx ON public.connection USING btree (code);


--
-- Name: connection_name_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX connection_name_idx ON public.connection USING btree (name);


--
-- Name: credit_sum_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX credit_sum_idx ON public.credit USING btree (shop_id, status, paid_at, amount);


--
-- Name: export_attempt_account_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX export_attempt_account_id_idx ON public.export_attempt USING btree (account_id);


--
-- Name: ffm_active_supplier_key; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX ffm_active_supplier_key ON public.fulfillment USING btree (order_id, public.ffm_active_supplier(supplier_id, status));


--
-- Name: fulfillment_address_to_district_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_address_to_district_code_idx ON public.fulfillment USING btree (address_to_district_code);


--
-- Name: fulfillment_address_to_province_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_address_to_province_code_idx ON public.fulfillment USING btree (address_to_province_code);


--
-- Name: fulfillment_address_to_ward_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_address_to_ward_code_idx ON public.fulfillment USING btree (address_to_ward_code);


--
-- Name: fulfillment_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_created_at_idx ON public.fulfillment USING btree (created_at);


--
-- Name: fulfillment_external_shipping_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_external_shipping_code_idx ON public.fulfillment USING btree (external_shipping_code);


--
-- Name: fulfillment_fulfillment_expected_pick_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_fulfillment_expected_pick_at_idx ON public.fulfillment USING btree (public.fulfillment_expected_pick_at(created_at));


--
-- Name: fulfillment_money_transaction_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_money_transaction_id_idx ON public.fulfillment USING btree (money_transaction_id);


--
-- Name: fulfillment_shipping_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_shipping_code_idx ON public.fulfillment USING btree (shipping_code);


--
-- Name: fulfillment_shipping_fee_shop_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_shipping_fee_shop_idx ON public.fulfillment USING btree (shipping_fee_shop);


--
-- Name: fulfillment_shipping_provider_external_shipping_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX fulfillment_shipping_provider_external_shipping_code_idx ON public.fulfillment USING btree (shipping_provider, external_shipping_code) WHERE (status <> ALL (ARRAY['-1'::integer, 1]));


--
-- Name: fulfillment_shipping_provider_shipping_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX fulfillment_shipping_provider_shipping_code_idx ON public.fulfillment USING btree (shipping_provider, shipping_code) WHERE (status <> ALL (ARRAY['-1'::integer, 1]));


--
-- Name: fulfillment_shipping_state_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_shipping_state_idx ON public.fulfillment USING btree (shipping_state);


--
-- Name: fulfillment_status_sum_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_status_sum_idx ON public.fulfillment USING btree (shop_id, status, shipping_status, etop_payment_status, shipping_fee_shop, total_cod_amount);


--
-- Name: fulfillment_total_cod_amount_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_total_cod_amount_idx ON public.fulfillment USING btree (total_cod_amount);


--
-- Name: fulfillment_updated_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX fulfillment_updated_at_idx ON public.fulfillment USING btree (updated_at);


--
-- Name: inventory_variant_shop_id_variant_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX inventory_variant_shop_id_variant_id_idx ON public.inventory_variant USING btree (shop_id, variant_id);


--
-- Name: inventory_voucher_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX inventory_voucher_created_at_idx ON public.inventory_voucher USING btree (created_at);


--
-- Name: inventory_voucher_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX inventory_voucher_id_idx ON public.inventory_voucher USING btree (id);


--
-- Name: inventory_voucher_ref_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX inventory_voucher_ref_code_idx ON public.inventory_voucher USING btree (ref_code);


--
-- Name: inventory_voucher_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX inventory_voucher_shop_id_code_idx ON public.inventory_voucher USING btree (shop_id, code);


--
-- Name: inventory_voucher_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX inventory_voucher_status_idx ON public.inventory_voucher USING btree (status);


--
-- Name: inventory_voucher_variant_ids_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX inventory_voucher_variant_ids_idx ON public.inventory_voucher USING gin (variant_ids);


--
-- Name: invitation_token_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX invitation_token_idx ON public.invitation USING btree (token);


--
-- Name: money_transaction_shipping_created_at_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX money_transaction_shipping_created_at_status_idx ON public.money_transaction_shipping USING btree (created_at, status);


--
-- Name: money_transaction_shipping_shop_id_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX money_transaction_shipping_shop_id_status_idx ON public.money_transaction_shipping USING btree (shop_id, status);


--
-- Name: order_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX order_code_idx ON public."order" USING btree (code);


--
-- Name: order_confirm_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_confirm_status_idx ON public."order" USING btree (confirm_status);


--
-- Name: order_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_created_at_idx ON public."order" USING btree (created_at);


--
-- Name: order_customer_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_customer_id_idx ON public."order" USING btree (customer_id);


--
-- Name: order_customer_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_customer_name_norm_idx ON public."order" USING gin (customer_name_norm);


--
-- Name: order_customer_phone_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_customer_phone_idx ON public."order" USING btree (customer_phone);


--
-- Name: order_etop_payment_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_etop_payment_status_idx ON public."order" USING btree (etop_payment_status);


--
-- Name: order_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_external_id_idx ON public."order" USING btree (external_order_id) WHERE (external_order_id IS NOT NULL);


--
-- Name: order_fulfillment_shipping_codes_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_fulfillment_shipping_codes_idx ON public."order" USING gin (fulfillment_shipping_codes);


--
-- Name: order_fulfillment_shipping_states_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_fulfillment_shipping_states_idx ON public."order" USING gin (fulfillment_shipping_states);


--
-- Name: order_fulfillment_shipping_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_fulfillment_shipping_status_idx ON public."order" USING btree (fulfillment_shipping_status);


--
-- Name: order_line_order_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_line_order_id_idx ON public.order_line USING btree (order_id);


--
-- Name: order_order_source_type_created_by_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_order_source_type_created_by_idx ON public."order" USING btree (order_source_type, created_by);


--
-- Name: order_partner_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX order_partner_external_id_idx ON public."order" USING btree (partner_id, external_order_id) WHERE ((external_order_id IS NOT NULL) AND (partner_id IS NOT NULL) AND (status <> '-1'::integer));


--
-- Name: order_partner_shop_id_external_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX order_partner_shop_id_external_code_idx ON public."order" USING btree (shop_id, ed_code, partner_id) WHERE ((partner_id IS NOT NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer));


--
-- Name: order_product_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_product_name_norm_idx ON public."order" USING gin (product_name_norm);


--
-- Name: order_shop_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX order_shop_external_id_idx ON public."order" USING btree (shop_id, external_order_id) WHERE ((external_order_id IS NOT NULL) AND (partner_id IS NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer) AND ((shop_id <> '1057792338951722956'::bigint) OR (created_at > '2019-02-14 02:00:00+00'::timestamp with time zone)));


--
-- Name: order_shop_id_ed_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX order_shop_id_ed_code_idx ON public."order" USING btree (shop_id, ed_code) WHERE ((partner_id IS NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer));


--
-- Name: order_shop_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_shop_id_idx ON public."order" USING btree (shop_id);


--
-- Name: order_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_status_idx ON public."order" USING btree (status);


--
-- Name: order_total_amount_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX order_total_amount_idx ON public."order" USING btree (total_amount);


--
-- Name: partner_relation_auth_key_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX partner_relation_auth_key_idx ON public.partner_relation USING btree (auth_key);


--
-- Name: partner_relation_partner_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX partner_relation_partner_id_idx ON public.partner_relation USING btree (partner_id);


--
-- Name: partner_relation_partner_id_subject_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX partner_relation_partner_id_subject_id_idx ON public.partner_relation USING btree (partner_id, subject_id);


--
-- Name: partner_relation_subject_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX partner_relation_subject_id_idx ON public.partner_relation USING btree (subject_id);


--
-- Name: payment_external_trans_id_payment_provider_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX payment_external_trans_id_payment_provider_idx ON public.payment USING btree (external_trans_id, payment_provider);


--
-- Name: purchase_order_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX purchase_order_shop_id_code_idx ON public.purchase_order USING btree (shop_id, code);


--
-- Name: purchase_order_variant_ids_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX purchase_order_variant_ids_idx ON public.purchase_order USING gin (variant_ids);


--
-- Name: purchase_refund_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX purchase_refund_shop_id_code_idx ON public.purchase_refund USING btree (shop_id, code);


--
-- Name: purchase_refund_shop_id_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX purchase_refund_shop_id_id_idx ON public.purchase_refund USING btree (shop_id, id);


--
-- Name: purchase_refund_shop_id_purchase_order_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX purchase_refund_shop_id_purchase_order_id_idx ON public.purchase_refund USING btree (shop_id, purchase_order_id);


--
-- Name: receipt_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_created_at_idx ON public.receipt USING btree (created_at);


--
-- Name: receipt_ledger_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_ledger_id_idx ON public.receipt USING btree (ledger_id);


--
-- Name: receipt_order_ids_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_order_ids_idx ON public.receipt USING gin (ref_ids);


--
-- Name: receipt_paid_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_paid_at_idx ON public.receipt USING btree (paid_at);


--
-- Name: receipt_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX receipt_shop_id_code_idx ON public.receipt USING btree (shop_id, code);


--
-- Name: receipt_trader_full_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_trader_full_name_norm_idx ON public.receipt USING gin (trader_full_name_norm);


--
-- Name: receipt_trader_full_name_norm_idx1; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_trader_full_name_norm_idx1 ON public.receipt USING gin (trader_full_name_norm);


--
-- Name: receipt_trader_phone_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_trader_phone_norm_idx ON public.receipt USING gin (trader_phone_norm);


--
-- Name: receipt_trader_type_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_trader_type_idx ON public.receipt USING btree (trader_type);


--
-- Name: receipt_type_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX receipt_type_idx ON public.receipt USING btree (type);


--
-- Name: refund_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX refund_shop_id_code_idx ON public.refund USING btree (shop_id, code);


--
-- Name: refund_shop_id_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX refund_shop_id_id_idx ON public.refund USING btree (shop_id, id);


--
-- Name: refund_shop_id_order_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX refund_shop_id_order_id_idx ON public.refund USING btree (shop_id, order_id);


--
-- Name: shipnow_fulfillment_order_ids_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shipnow_fulfillment_order_ids_idx ON public.shipnow_fulfillment USING gin (order_ids);


--
-- Name: shipping_source_name_type_username_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shipping_source_name_type_username_idx ON public.shipping_source USING btree (name, type, username);


--
-- Name: shop_brand_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_brand_partner_id_external_id_idx ON public.shop_brand USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_brand_shop_id_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_brand_shop_id_id_idx ON public.shop_brand USING btree (shop_id, id);


--
-- Name: shop_category_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_category_partner_id_external_id_idx ON public.shop_category USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_collection_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_collection_partner_id_external_id_idx ON public.shop_collection USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_connection_connection_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_connection_connection_id_idx ON public.shop_connection USING btree (connection_id) WHERE (is_global IS TRUE);


--
-- Name: shop_connection_shop_id_connection_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_connection_shop_id_connection_id_idx ON public.shop_connection USING btree (shop_id, connection_id) WHERE (deleted_at IS NULL);


--
-- Name: shop_customer_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_created_at_idx ON public.shop_customer USING btree (created_at);


--
-- Name: shop_customer_full_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_full_name_norm_idx ON public.shop_customer USING gin (full_name_norm);


--
-- Name: shop_customer_group_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_group_created_at_idx ON public.shop_customer_group USING btree (created_at);


--
-- Name: shop_customer_group_customer_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_group_customer_created_at_idx ON public.shop_customer_group_customer USING btree (created_at);


--
-- Name: shop_customer_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_customer_partner_id_external_id_idx ON public.shop_customer USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_customer_partner_id_shop_id_external_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_customer_partner_id_shop_id_external_code_idx ON public.shop_customer USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));


--
-- Name: shop_customer_phone_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_phone_idx ON public.shop_customer USING btree (phone);


--
-- Name: shop_customer_phone_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_phone_norm_idx ON public.shop_customer USING gin (phone_norm);


--
-- Name: shop_customer_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_customer_shop_id_code_idx ON public.shop_customer USING btree (shop_id, code);


--
-- Name: shop_customer_shop_id_email_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_shop_id_email_idx ON public.shop_customer USING btree (shop_id, email);


--
-- Name: shop_customer_shop_id_phone_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_shop_id_phone_idx ON public.shop_customer USING btree (shop_id, phone);


--
-- Name: shop_customer_type_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_type_idx ON public.shop_customer USING btree (type);


--
-- Name: shop_customer_updated_at_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_customer_updated_at_id_idx ON public.shop_customer USING btree (updated_at, id);


--
-- Name: shop_ledger_bank_account_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_ledger_bank_account_idx ON public.shop_ledger USING gin (bank_account);


--
-- Name: shop_ledger_bank_account_idx1; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_ledger_bank_account_idx1 ON public.shop_ledger USING gin (bank_account);


--
-- Name: shop_product_collection_partner_id_external_collection_id_e_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_collection_partner_id_external_collection_id_e_idx ON public.shop_product_collection USING btree (partner_id, external_collection_id, external_product_id) WHERE ((partner_id IS NOT NULL) AND (external_collection_id IS NOT NULL) AND (external_product_id IS NOT NULL));


--
-- Name: shop_product_list_price_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_list_price_idx ON public.shop_product USING btree (list_price);


--
-- Name: shop_product_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_name_norm_idx ON public.shop_product USING gin (name_norm);


--
-- Name: shop_product_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_partner_id_external_id_idx ON public.shop_product USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_product_partner_id_external_id_idx1; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_partner_id_external_id_idx1 ON public.shop_product USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_product_partner_id_shop_id_external_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_partner_id_shop_id_external_code_idx ON public.shop_product USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));


--
-- Name: shop_product_partner_id_shop_id_external_code_idx1; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_partner_id_shop_id_external_code_idx1 ON public.shop_product USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));


--
-- Name: shop_product_search_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_search_idx ON public.shop_product USING gin (name_norm);


--
-- Name: shop_product_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_shop_id_code_idx ON public.shop_product USING btree (shop_id, code) WHERE (deleted_at IS NULL);


--
-- Name: shop_product_shop_id_product_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_product_shop_id_product_id_idx ON public.shop_product USING btree (shop_id, product_id);


--
-- Name: shop_product_source_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_source_id_idx ON public.shop USING btree (product_source_id);


--
-- Name: shop_product_updated_at_product_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_updated_at_product_id_idx ON public.shop_product USING btree (updated_at, product_id);


--
-- Name: shop_product_updated_at_product_id_idx1; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_product_updated_at_product_id_idx1 ON public.shop_product USING btree (updated_at, product_id);


--
-- Name: shop_stocktake_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_stocktake_code_idx ON public.shop_stocktake USING btree (code);


--
-- Name: shop_stocktake_status_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_stocktake_status_idx ON public.shop_stocktake USING btree (status);


--
-- Name: shop_supplier_company_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_supplier_company_name_norm_idx ON public.shop_supplier USING gin (company_name_norm);


--
-- Name: shop_supplier_email_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_supplier_email_idx ON public.shop_supplier USING btree (email);


--
-- Name: shop_supplier_full_name_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_supplier_full_name_norm_idx ON public.shop_supplier USING gin (full_name_norm);


--
-- Name: shop_supplier_phone_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_supplier_phone_idx ON public.shop_supplier USING btree (phone);


--
-- Name: shop_supplier_phone_norm_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_supplier_phone_norm_idx ON public.shop_supplier USING gin (phone_norm);


--
-- Name: shop_supplier_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_supplier_shop_id_code_idx ON public.shop_supplier USING btree (shop_id, code);


--
-- Name: shop_variant_created_at_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_variant_created_at_idx ON public.shop_variant USING btree (created_at);


--
-- Name: shop_variant_partner_id_external_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_variant_partner_id_external_id_idx ON public.shop_variant USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));


--
-- Name: shop_variant_partner_id_shop_id_external_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_variant_partner_id_shop_id_external_code_idx ON public.shop_variant USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));


--
-- Name: shop_variant_product_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_variant_product_id_idx ON public.shop_variant USING btree (product_id);


--
-- Name: shop_variant_shop_id_code_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_variant_shop_id_code_idx ON public.shop_variant USING btree (shop_id, code) WHERE (deleted_at IS NULL);


--
-- Name: shop_variant_supplier_shop_id_supplier_id_variant_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX shop_variant_supplier_shop_id_supplier_id_variant_id_idx ON public.shop_variant_supplier USING btree (shop_id, supplier_id, variant_id);


--
-- Name: shop_variant_supplier_supplier_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_variant_supplier_supplier_id_idx ON public.shop_variant_supplier USING btree (supplier_id);


--
-- Name: shop_variant_supplier_variant_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_variant_supplier_variant_id_idx ON public.shop_variant_supplier USING btree (variant_id);


--
-- Name: shop_variant_updated_at_variant_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE INDEX shop_variant_updated_at_variant_id_idx ON public.shop_variant USING btree (updated_at, variant_id);


--
-- Name: user_email_key; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX user_email_key ON public."user" USING btree (email) WHERE (wl_partner_id IS NULL);


--
-- Name: user_email_wl_partner_id_idx; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX user_email_wl_partner_id_idx ON public."user" USING btree (email, wl_partner_id);


--
-- Name: user_phone_key; Type: INDEX; Schema: public; Owner: etop
--

CREATE UNIQUE INDEX user_phone_key ON public."user" USING btree (phone) WHERE (wl_partner_id IS NULL);


--
-- Name: account notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.account FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: address notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.address FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: etop_category notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.etop_category FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: fulfillment notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: money_transaction_shipping notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.money_transaction_shipping FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: order notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history."order" FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: order_external notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.order_external FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: order_source notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.order_source FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product_brand notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_brand FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product_external notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_external FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product_shop_collection notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_shop_collection FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_product_collection();


--
-- Name: product_source notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product_source_category notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source_category FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: product_source_category_external notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source_category_external FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop_collection notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_collection FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop_customer notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop_customer_group notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop_customer_group_customer notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group_customer FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_customer_group_customer();


--
-- Name: shop_product notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_product FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_product();


--
-- Name: shop_trader_address notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_trader_address FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: shop_variant notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_variant FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_variant();


--
-- Name: supplier notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.supplier FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: user notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history."user" FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: variant notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.variant FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: variant_external notify_pgrid; Type: TRIGGER; Schema: history; Owner: etop
--

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.variant_external FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();


--
-- Name: partner account_update; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON public.partner FOR EACH ROW EXECUTE PROCEDURE public.update_to_account();


--
-- Name: shop account_update; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON public.shop FOR EACH ROW EXECUTE PROCEDURE public.update_to_account();


--
-- Name: fulfillment fulfillment_update_order_status; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER fulfillment_update_order_status AFTER UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_order_status();


--
-- Name: fulfillment fulfillment_update_shipping_fees; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER fulfillment_update_shipping_fees BEFORE INSERT OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_shipping_fees();


--
-- Name: fulfillment fulfillment_update_status; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER fulfillment_update_status BEFORE INSERT OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_status();


--
-- Name: order order_update_status; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER order_update_status BEFORE INSERT OR UPDATE ON public."order" FOR EACH ROW EXECUTE PROCEDURE public.order_update_status();


--
-- Name: account save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.account FOR EACH ROW EXECUTE PROCEDURE public.save_history('account', '{id}');


--
-- Name: account_auth save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.account_auth FOR EACH ROW EXECUTE PROCEDURE public.save_history('account_auth', '{key,account_id}');


--
-- Name: account_user save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.account_user FOR EACH ROW EXECUTE PROCEDURE public.save_history('account_user', '{account_id,user_id}');


--
-- Name: address save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.address FOR EACH ROW EXECUTE PROCEDURE public.save_history('address', '{id,account_id}');


--
-- Name: connection save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.connection FOR EACH ROW EXECUTE PROCEDURE public.save_history('connection', '{id}');


--
-- Name: credit save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.credit FOR EACH ROW EXECUTE PROCEDURE public.save_history('credit', '{id,supplier_id,shop_id}');


--
-- Name: fulfillment save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.save_history('fulfillment', '{id,order_id,variant_ids,shop_id,supplier_id}', '{id}');


--
-- Name: import_attempt save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.import_attempt FOR EACH ROW EXECUTE PROCEDURE public.save_history('import_attempt', '{id,user_id,account_id}');


--
-- Name: inventory_variant save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.inventory_variant FOR EACH ROW EXECUTE PROCEDURE public.save_history('inventory_variant', '{variant_id,shop_id}');


--
-- Name: inventory_voucher save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.inventory_voucher FOR EACH ROW EXECUTE PROCEDURE public.save_history('inventory_voucher', '{id,shop_id}');


--
-- Name: invitation save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.invitation FOR EACH ROW EXECUTE PROCEDURE public.save_history('invitation', '{id,account_id}');


--
-- Name: money_transaction_shipping save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.money_transaction_shipping FOR EACH ROW EXECUTE PROCEDURE public.save_history('money_transaction_shipping', '{id}');


--
-- Name: order save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public."order" FOR EACH ROW EXECUTE PROCEDURE public.save_history('order', '{id,shop_id,partner_id}');


--
-- Name: order_line save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.order_line FOR EACH ROW EXECUTE PROCEDURE public.save_history('order_line', '{order_id,variant_id,shop_id,supplier_id}');


--
-- Name: order_source_internal save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.order_source_internal FOR EACH ROW EXECUTE PROCEDURE public.save_history('order_source_internal', '{id}');


--
-- Name: partner save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.partner FOR EACH ROW EXECUTE PROCEDURE public.save_history('partner', '{id}');


--
-- Name: partner_relation save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.partner_relation FOR EACH ROW EXECUTE PROCEDURE public.save_history('partner_relation', '{partner_id,subject_id}');


--
-- Name: payment save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.payment FOR EACH ROW EXECUTE PROCEDURE public.save_history('payment', '{id}');


--
-- Name: product_shop_collection save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.product_shop_collection FOR EACH ROW EXECUTE PROCEDURE public.save_history('product_shop_collection', '{product_id,collection_id}');


--
-- Name: purchase_order save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.purchase_order FOR EACH ROW EXECUTE PROCEDURE public.save_history('purchase_order', '{id,shop_id}');


--
-- Name: purchase_refund save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.purchase_refund FOR EACH ROW EXECUTE PROCEDURE public.save_history('purchase_refund', '{id,shop_id}');


--
-- Name: receipt save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.receipt FOR EACH ROW EXECUTE PROCEDURE public.save_history('receipt', '{id,shop_id}');


--
-- Name: refund save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.refund FOR EACH ROW EXECUTE PROCEDURE public.save_history('refund', '{id,shop_id}');


--
-- Name: shipnow_fulfillment save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shipnow_fulfillment FOR EACH ROW EXECUTE PROCEDURE public.save_history('shipnow_fulfillment', '{id,shop_id}');


--
-- Name: shop save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop', '{id}');


--
-- Name: shop_brand save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_brand FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_brand', '{id,shop_id}');


--
-- Name: shop_carrier save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_carrier FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_carrier', '{id,shop_id}');


--
-- Name: shop_category save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_category FOR EACH ROW EXECUTE PROCEDURE public.save_history('product_source_category', '{id}');


--
-- Name: shop_collection save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_collection FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_collection', '{id,shop_id}');


--
-- Name: shop_connection save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_connection FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_connection', '{shop_id,connection_id}');


--
-- Name: shop_customer save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_customer FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_customer', '{id,shop_id}');


--
-- Name: shop_customer_group save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_customer_group FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_customer_group', '{id,shop_id}');


--
-- Name: shop_customer_group_customer save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_customer_group_customer FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_customer_group_customer', '{customer_id,group_id}');


--
-- Name: shop_ledger save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_ledger FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_ledger', '{id,shop_id}');


--
-- Name: shop_product save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_product FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_product', '{shop_id,product_id}');


--
-- Name: shop_product_collection save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_product_collection FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_product_collection', '{shop_id,product_id,collection_id}');


--
-- Name: shop_stocktake save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_stocktake FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_stocktake', '{id,shop_id}');


--
-- Name: shop_supplier save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_supplier FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_supplier', '{id,shop_id}');


--
-- Name: shop_trader save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_trader FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_trader', '{id,shop_id}');


--
-- Name: shop_trader_address save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_trader_address FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_trader_address', '{id,shop_id,trader_id}');


--
-- Name: shop_variant save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_variant FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_variant', '{shop_id,variant_id}');


--
-- Name: shop_variant_supplier save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.shop_variant_supplier FOR EACH ROW EXECUTE PROCEDURE public.save_history('shop_variant_supplier', '{variant_id,supplier_id}');


--
-- Name: transaction save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.transaction FOR EACH ROW EXECUTE PROCEDURE public.save_history('transaction', '{id}');


--
-- Name: user save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public."user" FOR EACH ROW EXECUTE PROCEDURE public.save_history('user', '{id}');


--
-- Name: user_internal save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.user_internal FOR EACH ROW EXECUTE PROCEDURE public.save_history('user_internal', '{id}');


--
-- Name: webhook save_history; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER save_history AFTER INSERT OR DELETE OR UPDATE ON public.webhook FOR EACH ROW EXECUTE PROCEDURE public.save_history('webhook', '{id,account_id}');


--
-- Name: account update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.account FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_account_seq');


--
-- Name: account_auth update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.account_auth FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_account_auth_seq');


--
-- Name: account_user update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.account_user FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_account_user_seq');


--
-- Name: address update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.address FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_address_seq');


--
-- Name: connection update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.connection FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_connection_seq');


--
-- Name: credit update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.credit FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_credit_seq');


--
-- Name: fulfillment update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_fulfillment_seq');


--
-- Name: import_attempt update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.import_attempt FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_import_attempt_seq');


--
-- Name: inventory_variant update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.inventory_variant FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_inventory_variant_seq');


--
-- Name: inventory_voucher update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.inventory_voucher FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_inventory_voucher_seq');


--
-- Name: invitation update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.invitation FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_invitation_seq');


--
-- Name: money_transaction_shipping update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.money_transaction_shipping FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_money_transaction_seq');


--
-- Name: order update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public."order" FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_order_seq');


--
-- Name: order_line update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.order_line FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_order_line_seq');


--
-- Name: order_source_internal update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.order_source_internal FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_order_source_internal_seq');


--
-- Name: partner update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.partner FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_partner_seq');


--
-- Name: partner_relation update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.partner_relation FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_partner_relation_seq');


--
-- Name: payment update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.payment FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_payment_seq');


--
-- Name: product_shop_collection update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.product_shop_collection FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_product_shop_collection_seq');


--
-- Name: purchase_order update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.purchase_order FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_purchase_order_seq');


--
-- Name: purchase_refund update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.purchase_refund FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_purchase_refund_seq');


--
-- Name: receipt update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.receipt FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_receipt_seq');


--
-- Name: refund update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.refund FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_refund_seq');


--
-- Name: shipnow_fulfillment update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shipnow_fulfillment FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shipnow_fulfillment_seq');


--
-- Name: shop update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_seq');


--
-- Name: shop_brand update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_brand FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_brand_seq');


--
-- Name: shop_carrier update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_carrier FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_carrier_seq');


--
-- Name: shop_category update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_category FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_product_source_category_seq');


--
-- Name: shop_collection update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_collection FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_collection_seq');


--
-- Name: shop_connection update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_connection FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_connection_seq');


--
-- Name: shop_customer update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_customer FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_customer_seq');


--
-- Name: shop_customer_group update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_customer_group FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_customer_group_seq');


--
-- Name: shop_customer_group_customer update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_customer_group_customer FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_customer_group_customer_seq');


--
-- Name: shop_ledger update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_ledger FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_ledger_seq');


--
-- Name: shop_product update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_product FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_product_seq');


--
-- Name: shop_product_collection update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_product_collection FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_product_collection_seq');


--
-- Name: shop_stocktake update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_stocktake FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_stocktake_seq');


--
-- Name: shop_supplier update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_supplier FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_supplier_seq');


--
-- Name: shop_trader update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_trader FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_trader_seq');


--
-- Name: shop_trader_address update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_trader_address FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_trader_address_seq');


--
-- Name: shop_variant update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_variant FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_variant_seq');


--
-- Name: shop_variant_supplier update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.shop_variant_supplier FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_shop_variant_supplier_seq');


--
-- Name: transaction update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.transaction FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_transaction_seq');


--
-- Name: user update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public."user" FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_user_seq');


--
-- Name: user_internal update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.user_internal FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_user_internal_seq');


--
-- Name: webhook update_rid; Type: TRIGGER; Schema: public; Owner: etop
--

CREATE TRIGGER update_rid BEFORE INSERT OR UPDATE ON public.webhook FOR EACH ROW EXECUTE PROCEDURE public.update_rid('history_webhook_seq');


--
-- Name: order order_customer_id_fkey; Type: FK CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history."order"
  ADD CONSTRAINT order_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.shop_customer(id);


--
-- Name: order order_trading_shop_id_fkey; Type: FK CONSTRAINT; Schema: history; Owner: etop
--

ALTER TABLE ONLY history."order"
  ADD CONSTRAINT order_trading_shop_id_fkey FOREIGN KEY (trading_shop_id) REFERENCES public.shop(id);


--
-- Name: account_auth account_auth_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_auth
  ADD CONSTRAINT account_auth_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: account account_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account
  ADD CONSTRAINT account_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);


--
-- Name: account_user account_user_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: account_user account_user_disabled_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_disabled_by_fkey FOREIGN KEY (disabled_by) REFERENCES public."user"(id);


--
-- Name: account_user account_user_invitation_sent_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_invitation_sent_by_fkey FOREIGN KEY (invitation_sent_by) REFERENCES public."user"(id);


--
-- Name: account_user account_user_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: address address_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.address
  ADD CONSTRAINT address_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: affiliate affiliate_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);


--
-- Name: affiliate affiliate_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);


--
-- Name: connection connection_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.connection
  ADD CONSTRAINT connection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: credit credit_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.credit
  ADD CONSTRAINT credit_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: district district_province_code_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.district
  ADD CONSTRAINT district_province_code_fkey FOREIGN KEY (province_code) REFERENCES public.province(code);


--
-- Name: export_attempt export_attempt_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: export_attempt export_attempt_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: external_account_ahamove external_account_ahamove_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);


--
-- Name: external_account_haravan external_account_haravan_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.external_account_haravan
  ADD CONSTRAINT external_account_haravan_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: fulfillment fulfillment_connection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_connection_id_fkey FOREIGN KEY (connection_id) REFERENCES public.connection(id);


--
-- Name: fulfillment fulfillment_money_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_money_transaction_id_fkey FOREIGN KEY (money_transaction_id) REFERENCES public.money_transaction_shipping(id);


--
-- Name: fulfillment fulfillment_money_transaction_shipping_external_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_money_transaction_shipping_external_id_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);


--
-- Name: fulfillment fulfillment_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_order_id_fkey FOREIGN KEY (order_id) REFERENCES public."order"(id);


--
-- Name: fulfillment fulfillment_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: fulfillment fulfillment_shop_carrier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_shop_carrier_id_fkey FOREIGN KEY (shop_carrier_id) REFERENCES public.shop_carrier(id);


--
-- Name: fulfillment fulfillment_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: import_attempt import_attempt_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: import_attempt import_attempt_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: inventory_voucher inventory_voucher_trader_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.inventory_voucher
  ADD CONSTRAINT inventory_voucher_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);


--
-- Name: invitation invitation_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.shop(id);


--
-- Name: invitation invitation_invited_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_invited_by_fkey FOREIGN KEY (invited_by) REFERENCES public."user"(id);


--
-- Name: money_transaction_shipping money_transaction_money_transaction_shipping_external_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_money_transaction_shipping_external_id_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);


--
-- Name: money_transaction_shipping_external_line money_transaction_shipping_ex_money_transaction_shipping_e_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping_external_line
  ADD CONSTRAINT money_transaction_shipping_ex_money_transaction_shipping_e_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);


--
-- Name: money_transaction_shipping money_transaction_shipping_money_transaction_shipping_etop_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_shipping_money_transaction_shipping_etop_fkey FOREIGN KEY (money_transaction_shipping_etop_id) REFERENCES public.money_transaction_shipping_etop(id);


--
-- Name: money_transaction_shipping money_transaction_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: order order_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.shop_customer(id);


--
-- Name: order_line order_line_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.order_line
  ADD CONSTRAINT order_line_order_id_fkey FOREIGN KEY (order_id) REFERENCES public."order"(id);


--
-- Name: order order_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: order order_payment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payment(id);


--
-- Name: order order_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: order order_trading_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_trading_shop_id_fkey FOREIGN KEY (trading_shop_id) REFERENCES public.shop(id);


--
-- Name: partner partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);


--
-- Name: partner partner_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);


--
-- Name: partner_relation partner_relation_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.partner_relation
  ADD CONSTRAINT partner_relation_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: product_shop_collection product_shop_collection_collection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);


--
-- Name: product_shop_collection product_shop_collection_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_category product_source_category_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT product_source_category_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: purchase_order purchase_order_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);


--
-- Name: purchase_order purchase_order_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: purchase_order purchase_order_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.shop_supplier(id);


--
-- Name: receipt receipt_ledger_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_ledger_id_fkey FOREIGN KEY (ledger_id) REFERENCES public.shop_ledger(id);


--
-- Name: receipt receipt_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: receipt receipt_shop_ledger_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_shop_ledger_id_fkey FOREIGN KEY (shop_ledger_id) REFERENCES public.shop_ledger(id);


--
-- Name: receipt receipt_trader_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);


--
-- Name: receipt receipt_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_user_id_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);


--
-- Name: shipnow_fulfillment shipnow_fulfillment_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shipnow_fulfillment shipnow_fulfillment_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.account(id);


--
-- Name: shop shop_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(id);


--
-- Name: shop_brand shop_brand_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_brand
  ADD CONSTRAINT shop_brand_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_carrier shop_carrier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);


--
-- Name: shop_carrier shop_carrier_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_category shop_category_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT shop_category_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_collection shop_collection_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_collection shop_collection_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_connection shop_connection_connection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_connection
  ADD CONSTRAINT shop_connection_connection_id_fkey FOREIGN KEY (connection_id) REFERENCES public.connection(id);


--
-- Name: shop_connection shop_connection_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_connection
  ADD CONSTRAINT shop_connection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_customer_group_customer shop_customer_group_customer_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.shop_customer(id);


--
-- Name: shop_customer_group_customer shop_customer_group_customer_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.shop_customer_group(id);


--
-- Name: shop_customer_group shop_customer_group_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer_group
  ADD CONSTRAINT shop_customer_group_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_customer shop_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);


--
-- Name: shop_customer shop_customer_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_customer shop_customer_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);


--
-- Name: shop_ledger shop_ledger_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);


--
-- Name: shop_ledger shop_ledger_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop shop_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);


--
-- Name: shop_product_collection shop_product_collection_collection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);


--
-- Name: shop_product shop_product_collection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);


--
-- Name: shop_product_collection shop_product_collection_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_product_collection shop_product_collection_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.shop_product(product_id);


--
-- Name: shop_product_collection shop_product_collection_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_product shop_product_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_product shop_product_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop shop_ship_from_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_ship_from_address_id_fkey FOREIGN KEY (ship_from_address_id) REFERENCES public.address(id);


--
-- Name: shop shop_ship_to_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_ship_to_address_id_fkey FOREIGN KEY (ship_to_address_id) REFERENCES public.address(id);


--
-- Name: shop_stocktake shop_stocktake_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_stocktake
  ADD CONSTRAINT shop_stocktake_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_trader_address shop_trader_address_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_trader_address shop_trader_address_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_trader_address shop_trader_address_trader_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);


--
-- Name: shop_variant shop_variant_collection_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);


--
-- Name: shop_variant shop_variant_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);


--
-- Name: shop_variant shop_variant_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_variant_supplier shop_variant_supplier_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop_variant_supplier shop_variant_supplier_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.shop_supplier(id);


--
-- Name: shop_variant_supplier shop_variant_supplier_variant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.shop_variant(variant_id);


--
-- Name: shop_supplier shop_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);


--
-- Name: shop_supplier shop_vendor_shop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);


--
-- Name: shop shop_wl_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_wl_partner_id_fkey FOREIGN KEY (wl_partner_id) REFERENCES public.partner(id);


--
-- Name: user_auth user_auth_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.user_auth
  ADD CONSTRAINT user_auth_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);


--
-- Name: user_internal user_internal_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.user_internal
  ADD CONSTRAINT user_internal_id_fkey FOREIGN KEY (id) REFERENCES public."user"(id);


--
-- Name: user user_ref_sale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_ref_sale_id_fkey FOREIGN KEY (ref_sale_id) REFERENCES public."user"(id);


--
-- Name: user user_ref_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_ref_user_id_fkey FOREIGN KEY (ref_user_id) REFERENCES public."user"(id);


--
-- Name: user user_wl_partner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_wl_partner_id_fkey FOREIGN KEY (wl_partner_id) REFERENCES public.partner(id);


--
-- Name: shop_variant variant_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT variant_product_id_fkey FOREIGN KEY (shop_id, product_id) REFERENCES public.shop_product(shop_id, product_id);


--
-- Name: webhook webhook_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.webhook
  ADD CONSTRAINT webhook_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: webhook_changes webhook_changes_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);


--
-- Name: webhook_changes webhook_changes_webhook_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: etop
--

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_webhook_id_fkey FOREIGN KEY (webhook_id) REFERENCES public.webhook(id);

