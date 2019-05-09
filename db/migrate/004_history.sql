CREATE EXTENSION hstore;
CREATE SCHEMA history;

CREATE TYPE tg_op_type AS ENUM ('INSERT', 'UPDATE', 'DELETE');

CREATE OR REPLACE FUNCTION latest_rid(_tbl regclass, OUT result INT8)
  LANGUAGE plpgsql
  AS $$
BEGIN
  EXECUTE format('SELECT rid FROM %s WHERE rid IS NOT NULL ORDER BY rid DESC LIMIT 1::int', _tbl)
  INTO result;
  result := COALESCE(result, 0);
END
$$;

CREATE OR REPLACE FUNCTION update_rid() RETURNS trigger
  LANGUAGE plpgsql
  AS $$
BEGIN
  NEW.rid = nextval(TG_ARGV[0]);
  RETURN NEW;
END
$$;

-- save_history('code', '{code}')
CREATE OR REPLACE FUNCTION save_history() RETURNS trigger
  LANGUAGE plpgsql
  AS $$
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
    IF (rnew - rold - '{rid,updated_at,last_sync_at}'::TEXT[] = '') THEN
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
$$;

-- init_history('code', {code}', '{code}')
--
-- cols: columns which always store value (instead of NULL on unchanged)
-- idx : columns to automatically create index
CREATE OR REPLACE FUNCTION init_history(
  tbl TEXT,
  cols TEXT[],
  idx TEXT[] DEFAULT NULL
) RETURNS TEXT
  LANGUAGE plpgsql
  AS $$
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
$$;

-- config history
SELECT init_history('account', '{id}');
SELECT init_history('account_user', '{account_id, user_id}');
SELECT init_history('address', '{id,account_id}', '{id}');
SELECT init_history('etop_category', '{id}');
SELECT init_history('fulfillment', '{id,order_id,variant_ids,shop_id,supplier_id}', '{id}');
SELECT init_history('order_external', '{id}');
SELECT init_history('order_line', '{order_id,variant_id,shop_id,supplier_id}', '{order_id,variant_id}');
SELECT init_history('order_source_internal', '{id}');
SELECT init_history('order_source', '{id}');
SELECT init_history('order', '{id,shop_id}', '{id}');
SELECT init_history('product_brand', '{id,supplier_id}', '{id}');
SELECT init_history('product_external', '{id,supplier_id}', '{id}');
SELECT init_history('product_shop_collection', '{product_id,collection_id}');
SELECT init_history('product_source_category_external', '{id}');
SELECT init_history('product_source_category', '{id}');
SELECT init_history('product_source_internal', '{id}');
SELECT init_history('product_source', '{id}');
SELECT init_history('product', '{id}');
SELECT init_history('shop_collection', '{id,shop_id}', '{id}');
SELECT init_history('shop_product', '{shop_id,product_id}');
SELECT init_history('shop_variant', '{shop_id,variant_id}');
SELECT init_history('shop', '{id}');
SELECT init_history('supplier', '{id}');
SELECT init_history('user_internal', '{id}');
SELECT init_history('user', '{id}');
SELECT init_history('variant_external', '{id}');
SELECT init_history('variant', '{id}');

CREATE INDEX ON history.fulfillment (order_id);
CREATE INDEX ON history.fulfillment (shop_id);
CREATE INDEX ON history.order_line (shop_id);
CREATE INDEX ON history.order (shop_id);
CREATE INDEX ON history.address (account_id);
