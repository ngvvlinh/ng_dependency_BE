-- Update order fulfillment_type to type `shipment (10)`
UPDATE "order" set fulfillment_type = 10 WHERE id in (
	SELECT o.id from "order" as o, fulfillment as f where o.id = f.order_id and fulfillment_type = 0
);
