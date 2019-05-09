CREATE FUNCTION notify_pgrid_shop_product() RETURNS trigger
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
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_product
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_shop_product();
