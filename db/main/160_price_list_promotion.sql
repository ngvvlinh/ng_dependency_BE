CREATE TABLE shipment_price_list_promotion (
    id INT8 PRIMARY KEY
    , price_list_id INT8 REFERENCES shipment_price_list(id)
    , name TEXT
    , description TEXT
    , status INT2
    , date_from TIMESTAMPTZ
    , date_to TIMESTAMPTZ
    , applied_rules JSONB
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
    , connection_id INT8 REFERENCES connection(id)
    , priority_point INT2
);

SELECT init_history('shipment_price_list_promotion', '{id,connection_id}');

ALTER TABLE fulfillment
    ADD COLUMN shipment_price_info JSONB;

ALTER TABLE history.fulfillment
    ADD COLUMN shipment_price_info JSONB;
