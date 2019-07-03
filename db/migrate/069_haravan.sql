CREATE TABLE external_account_haravan (
    id INT8 PRIMARY KEY
    , shop_id INT8 NOT NULL REFERENCES shop(id)
    , subdomain text
    , external_shop_id INT
    , external_carrier_service_id INT
    , external_connected_carrier_service_at TIMESTAMP  WITH TIME ZONE
    , access_token text
    , expires_at TIMESTAMP WITH TIME ZONE
    , created_at TIMESTAMP WITH TIME ZONE
    , updated_at TIMESTAMP WITH TIME ZONE
);
