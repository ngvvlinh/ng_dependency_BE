CREATE FUNCTION notify_pgrid_id() RETURNS trigger
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
END$$ LANGUAGE 'plpgsql';

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.account
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
-- account_user
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.address
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.etop_category
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fulfillment
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.order_external
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
-- order_line
-- order_source_internal
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.order_source
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history."order"
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_brand
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_external
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
-- product_shop_collection
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source_category_external
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source_category
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_source
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_collection
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
-- shop_product
-- shop_variant
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.supplier
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
-- user_internal
CREATE TRIGGER notify_pgrid AFTER INSERT ON history."user"
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.variant_external
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
CREATE TRIGGER notify_pgrid AFTER INSERT ON history.variant
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
