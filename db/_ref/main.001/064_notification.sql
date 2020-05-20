CREATE TRIGGER notify_pgrid AFTER INSERT ON history.money_transaction_shipping
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
