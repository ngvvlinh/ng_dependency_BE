
-- include partner_id in history of order and fulfillment
DROP TRIGGER save_history ON "order";
DROP TRIGGER save_history ON fulfillment;

CREATE TRIGGER save_history AFTER INSERT OR UPDATE OR DELETE ON "order"
  FOR EACH ROW EXECUTE PROCEDURE save_history('order', '{id,shop_id,partner_id}');
CREATE TRIGGER save_history AFTER INSERT OR UPDATE OR DELETE ON fulfillment
  FOR EACH ROW EXECUTE PROCEDURE save_history('fulfillment', '{id,order_id,variant_ids,shop_id,supplier_id,partner_id}');

-- webhook

CREATE TABLE webhook (
  id INT8 PRIMARY KEY,
  account_id INT8 NOT NULL REFERENCES account(id),
  entities TEXT[],
  fields TEXT[],
  url TEXT NOT NULL,
  metadata TEXT,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

SELECT init_history('webhook', '{id,account_id}', '{id}');
