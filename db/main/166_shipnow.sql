ALTER TABLE shop_connection
    ADD COLUMN owner_id INT8 REFERENCES "user"(id);

ALTER TABLE history.shop_connection
    ADD COLUMN owner_id INT8 REFERENCES "user"(id);

-- Insert ahamove connection direct
INSERT INTO "connection" ("id", "name", "status", "created_at", "updated_at", "driver", "connection_type", "connection_subtype", "connection_method", "connection_provider", "etop_affiliate_account", "code") VALUES
('1000212023297494791', 'Ahamove', '1', NOW(), NOW(), 'shipping/shipnow/direct/ahamove', 'shipping', 'shipnow', 'direct', 'ahamove', '{ "token": "", "user_id": ""}', 'A51B');

ALTER TABLE shipnow_fulfillment
    ADD COLUMN coupon TEXT;

ALTER TABLE history.shipnow_fulfillment
    ADD COLUMN coupon TEXT;

CREATE UNIQUE INDEX ON shop_connection(owner_id, connection_id) where deleted_at is null;

-- add connection_id to account_shipnow (external_account_ahamove)
-- default value is ahamove direct connection (1000212023297494791)
ALTER TABLE external_account_ahamove
    ADD COLUMN connection_id INT8 REFERENCES connection(id)
    DEFAULT '1000212023297494791';

ALTER TABLE external_account_ahamove DROP CONSTRAINT external_account_ahamove_phone_key;
CREATE UNIQUE INDEX ON external_account_ahamove(owner_id, phone);
