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
	ADD COLUMN fulfillment_type int2,
	ADD COLUMN fulfillment_ids int8[];

ALTER TABLE history."order"
	ADD COLUMN fulfillment_type int2,
	ADD COLUMN fulfillment_ids int8[];

ALTER table  shipnow_fulfillment
    ADD COLUMN fee_lines JSONB,
    ADD COLUMN carrier_fee_lines JSONB,
    ADD COLUMN total_fee INT,
    ADD COLUMN shipping_created_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN etop_payment_status INT,
    ADD COLUMN cod_etop_transfered_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN shipping_picking_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN shipping_delivering_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN shipping_delivered_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN shipping_cancelled_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX ON shipnow_fulfillment USING gin(order_ids);

CREATE TABLE  external_account_ahamove (
    id INT8 PRIMARY KEY
    , owner_id INT8 NOT NULL REFERENCES "user"(id)
    , phone TEXT NOT NULL UNIQUE
    , name TEXT NOT NULL
    , external_token TEXT
    , external_verified BOOLEAN
    , created_at TIMESTAMP WITH TIME ZONE
    , updated_at TIMESTAMP WITH TIME ZONE
    , external_created_at TIMESTAMP WITH TIME ZONE
    , last_send_verified_at TIMESTAMP WITH TIME ZONE
    , external_ticket_id TEXT
    , external_id TEXT
);

ALTER TABLE "external_account_ahamove"
    ADD COLUMN id_card_front_img TEXT
    , ADD COLUMN id_card_back_img TEXT
    , ADD COLUMN portrait_img TEXT
    , ADD COLUMN uploaded_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE shipnow_fulfillment
    ADD COLUMN shipping_service_name TEXT
    , ADD COLUMN shipping_service_description TEXT
    , ADD COLUMN cancel_reason TEXT
    , ADD COLUMN shipping_shared_link TEXT;
