ALTER TABLE shop_trader_address
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE history.shop_trader_address
    ADD COLUMN partner_id INT8;

SELECT init_history('shop_customer_group', '{id,shop_id}');

ALTER TABLE shop_customer_group
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE history.shop_customer_group
    ADD COLUMN partner_id INT8;

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

SELECT init_history('inventory_variant', '{variant_id,shop_id}');

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.inventory_variant
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_trader_address
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

SELECT init_history('shop_customer_group_customer', '{customer_id, group_id}');

CREATE FUNCTION notify_pgrid_shop_customer_group_customer() RETURNS trigger
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
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group_customer
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_shop_customer_group_customer();

CREATE FUNCTION notify_pgrid_shop_variant() RETURNS trigger
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
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_variant
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_shop_variant();

CREATE FUNCTION notify_pgrid_shop_product_collection() RETURNS trigger
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
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_shop_collection
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_shop_product_collection();