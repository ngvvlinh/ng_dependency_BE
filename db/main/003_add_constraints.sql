ALTER TABLE order_line ADD COLUMN shop_id INT8;
UPDATE order_line SET shop_id = "order".shop_id FROM "order" WHERE order_line.order_id = "order".id;
ALTER TABLE order_line ALTER COLUMN shop_id SET NOT NULL;
