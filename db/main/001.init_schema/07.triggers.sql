CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON public.partner FOR EACH ROW EXECUTE PROCEDURE public.update_to_account();

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON public.shop FOR EACH ROW EXECUTE PROCEDURE public.update_to_account();

CREATE TRIGGER fulfillment_update_order_status AFTER UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_order_status();

CREATE TRIGGER fulfillment_update_shipping_fees BEFORE INSERT OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_shipping_fees();

CREATE TRIGGER fulfillment_update_status BEFORE INSERT OR UPDATE ON public.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.fulfillment_update_status();

CREATE TRIGGER order_update_status BEFORE INSERT OR UPDATE ON public."order" FOR EACH ROW EXECUTE PROCEDURE public.order_update_status();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.account FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.address FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fulfillment FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.money_transaction_shipping FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history."order" FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.product_shop_collection FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_product_collection();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_collection FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_customer_group_customer FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_customer_group_customer();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_product FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_product();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_trader_address FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shop_variant FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_shop_variant();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history."user" FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();
