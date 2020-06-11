ALTER TABLE shipment_price_list
    DROP COLUMN shipment_sub_price_list_ids
    , ADD COLUMN connection_id INT8 REFERENCES connection(id);

ALTER TABLE history.shipment_price_list
    DROP COLUMN shipment_sub_price_list_ids
    , ADD COLUMN connection_id INT8;

ALTER TABLE shipment_price
    DROP COLUMN shipment_sub_price_list_id
    , ADD COLUMN shipment_price_list_id INT8 REFERENCES shipment_price_list(id);

ALTER TABLE history.shipment_price
    DROP COLUMN shipment_sub_price_list_id
    , ADD COLUMN shipment_price_list_id INT8;

DROP TABLE shipment_sub_price_list, history.shipment_sub_price_list;

-- shop_shipment_price_list
DROP INDEX shop_shipment_price_list_shop_id_shipment_price_list_id_idx;

ALTER TABLE shop_shipment_price_list
    ADD COLUMN connection_id INT8 REFERENCES connection(id)
    , ALTER COLUMN shipment_price_list_id SET NOT NULL;

ALTER TABLE history.shop_shipment_price_list
    ADD COLUMN connection_id INT8;

CREATE UNIQUE INDEX ON shop_shipment_price_list(shop_id, connection_id) where deleted_at is null;
