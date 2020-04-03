CREATE TABLE custom_region (
    id INT8 PRIMARY KEY
    , name TEXT
    , description TEXT
    , province_codes TEXT[]
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE INDEX ON custom_region USING GIN(province_codes);

CREATE TABLE shipment_service (
    id INT8 PRIMARY KEY
    , connection_id INT8 REFERENCES connection(id)
    , name TEXT
    , ed_code TEXT
    , service_ids TEXT[]
    , description TEXT
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE INDEX ON shipment_service USING GIN(service_ids);

CREATE TABLE shipment_price_list (
    id INT8 PRIMARY KEY
    , name TEXT
    , description TEXT
    , is_active bool
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

ALTER TABLE connection ADD COLUMN services JSONB;
ALTER TABLE history.connection ADD COLUMN services JSONB;

CREATE TABLE shipment_price (
    id INT8 PRIMARY KEY
    , name TEXT
    , shipment_price_list_id INT8 REFERENCES shipment_price_list(id)
    , shipment_service_id INT8 REFERENCES shipment_service(id)
    , custom_region_types TEXT[]
    , custom_region_ids INT8[]
    , region_types TEXT[]
    , province_types TEXT[]
    , urban_types TEXT[]
    , priority_point INT2
    , details JSONB
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);

ALTER TABLE shipment_price
    ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE shipment_service
    ADD COLUMN status INT2
    , ADD COLUMN image_url TEXT;

ALTER TABLE connection
    ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE history.connection
    ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);

ALTER TABLE shipment_price
    ADD COLUMN status INT2;
