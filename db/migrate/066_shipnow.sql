CREATE TABLE shipnow_fulfillment (
    id INT8 PRIMARY KEY
    , shop_id INT8 REFERENCES account(id)
    , partner_id INT8 REFERENCES partner(id)
    , order_ids INT8[]
    , pickup_address JSONB
    , carrier TEXT
    , shipping_service_code TEXT
    , shipping_service_fee INT
    , chargeable_weight INT
    , gross_weight INT
    , basket_value INT
    , cod_amount INT
    , shipping_note TEXT
    , request_pickup_at TIMESTAMP WITH TIME ZONE
    , delivery_points JSONB
    , status INT
    , shipping_state TEXT
    , sync_status INT
    , sync_states JSONB
    , last_sync_at TIMESTAMP WITH TIME ZONE
    , created_at TIMESTAMP WITH TIME ZONE
    , updated_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE adress ADD COLUMN coordinates JSONB;

ALTER TABLE shipnow_fulfillment
    ADD COLUMN confirm_status INT,
    ADD COLUMN shipping_status INT,
    ADD COLUMN shipping_code TEXT;

ALTER TABLE "order"
	ADD COLUMN fulfill int2,
	ADD COLUMN fulfill_ids int8[];

ALTER TABLE history."order"
	ADD COLUMN fulfill int2,
	ADD COLUMN fulfill_ids int8[];
