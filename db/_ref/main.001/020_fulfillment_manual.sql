ALTER TYPE shipping_provider ADD VALUE 'manual';
ALTER TABLE fulfillment ADD column shipping_status smallint;
ALTER TABLE history.fulfillment ADD column shipping_status smallint;

ALTER TABLE "order" RENAME payment_status TO customer_payment_status;
ALTER TABLE history."order" RENAME payment_status TO customer_payment_status;
ALTER TABLE "order" ADD COLUMN etop_payment_status SMALLINT;
ALTER TABLE history."order" ADD COLUMN etop_payment_status SMALLINT;

CREATE OR REPLACE FUNCTION ffm_update_status() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
	ffm_status smallint;
	payment_status smallint;
BEGIN
	IF NEW.shipping_provider = 'manual' THEN
		RETURN NEW;
	END IF;
	ffm_status := NEW.shipping_status;
	IF NEW.cod_etop_paid_at IS NULL THEN
		payment_status := 0;
	ELSE
		payment_status := 1;
	END IF;
	UPDATE "order" SET fulfillment_status = ffm_status, etop_payment_status = payment_status WHERE id = NEW.order_id;
   RETURN NEW;
END;
$$;
CREATE TRIGGER ffm_update_status AFTER UPDATE ON fulfillment FOR EACH ROW EXECUTE PROCEDURE ffm_update_status();

CREATE OR REPLACE FUNCTION order_update_status() RETURNS trigger
	LANGUAGE plpgsql
	AS $$
BEGIN
	IF OLD.status = 1 OR OLD.status = -1 THEN
		RETURN NEW;
	END IF;
	IF NEW.shop_confirm != 0 THEN
		NEW.confirm_status = NEW.shop_confirm;
	ELSE
		NEW.confirm_status = LEAST(NEW.shop_confirm, NEW.customer_confirm, NEW.external_confirm);
	END IF;
	RETURN NEW;
END;
$$;
CREATE TRIGGER order_update_status BEFORE UPDATE ON "order" FOR EACH ROW EXECUTE PROCEDURE order_update_status();
