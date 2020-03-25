CREATE TABLE if not exists shipnow_fulfillment (
    id INT8 PRIMARY KEY
    , shop_id INT8
    , partner_id INT8
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
    , address_to_province_code text
	, address_to_district_code text
	, confirm_status INT
    , shipping_status INT
    , shipping_code TEXT
    , shipping_service_name TEXT
    , shipping_service_description TEXT
    , cancel_reason TEXT
    , shipping_shared_link TEXT
    , fee_lines JSONB
    , carrier_fee_lines JSONB
    , total_fee INT
    , shipping_created_at TIMESTAMP WITH TIME ZONE
    , etop_payment_status INT
    , cod_etop_transfered_at TIMESTAMP WITH TIME ZONE
    , shipping_picking_at TIMESTAMP WITH TIME ZONE
    , shipping_delivering_at TIMESTAMP WITH TIME ZONE
    , shipping_delivered_at TIMESTAMP WITH TIME ZONE
    , shipping_cancelled_at TIMESTAMP WITH TIME ZONE
    , rid INT8
);