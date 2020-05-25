CREATE TABLE shipment_sub_price_list (
    id INT8 PRIMARY KEY
    , name TEXT
    , description TEXT
    , status INT2
    , connection_id INT8 REFERENCES connection(id)
    , wl_partner_id INT8 REFERENCES partner(id)
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);

ALTER TABLE shipment_price
    ADD COLUMN shipment_sub_price_list_id INT8 REFERENCES shipment_sub_price_list(id),
    DROP COLUMN shipment_price_list_id;

CREATE TABLE shop_shipment_price_list (
    shop_id INT8 REFERENCES shop(id)
    , shipment_price_list_id INT8 REFERENCES shipment_price_list(id)
    , note TEXT
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , updated_by INT8 REFERENCES "user"(id)
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE UNIQUE INDEX ON shop_shipment_price_list(shop_id, shipment_price_list_id) where deleted_at is null;

ALTER TABLE shipment_price_list
    ADD COLUMN shipment_sub_price_list_ids INT8[];

CREATE INDEX ON shipment_price_list USING GIN(shipment_sub_price_list_ids);

select init_history('shipment_service', '{id}');
select init_history('shipment_sub_price_list', '{id}');
select init_history('shipment_price_list', '{id}');
select init_history('shipment_price', '{id}');
select init_history('shop_shipment_price_list', '{shop_id}');
