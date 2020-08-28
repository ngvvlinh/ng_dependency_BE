ALTER TABLE shipnow_fulfillment
    ADD COLUMN connection_id INT8
    , ADD COLUMN connection_method TEXT
    , ADD COLUMN external_id TEXT;

ALTER TABLE history.shipnow_fulfillment
    ADD COLUMN connection_id INT8
    , ADD COLUMN connection_method TEXT
    , ADD COLUMN external_id TEXT;

CREATE INDEX shipnow_fulfillment_external_id_idx ON "shipnow_fulfillment" (external_id) WHERE external_id IS NOT NULL;

CREATE UNIQUE INDEX shipnow_fulfillment_partner_external_id_idx ON "shipnow_fulfillment" (partner_id, external_id)
  WHERE external_id IS NOT NULL AND partner_id IS NOT NULL AND status != -1;

CREATE UNIQUE INDEX shipnow_fulfillment_shop_external_id_idx ON "shipnow_fulfillment" (shop_id, external_id)
  WHERE external_id IS NOT NULL AND partner_id IS NULL AND status != -1;

CREATE UNIQUE INDEX ON shipnow_fulfillment (connection_id, shipping_code) where status not in (-1);

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.shipnow_fulfillment FOR EACH ROW EXECUTE PROCEDURE public.notify_pgrid_id();

-- Insert connection ahamove
INSERT INTO code ("code", "type", "created_at") VALUES
    ('50FE', 'connection', NOW());

INSERT INTO "connection" ("id", "name", "status", "created_at", "updated_at", "driver", "connection_type", "connection_subtype", "connection_method", "connection_provider", "code") VALUES
('1000343411864064400', 'TopShip - Ahamove', '1', NOW(), NOW(), 'shipping/shipnow/builtin/ahamove', 'shipping', 'shipnow', 'builtin', 'ahamove', '50FE');

INSERT INTO "shop_connection" ("connection_id", "status", "created_at", "updated_at", "is_global") VALUES
('1000343411864064400','1', NOW(), NOW(), 't');
