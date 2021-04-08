CREATE TABLE tenant (
    id INT8 PRIMARY KEY
    , owner_id INT8
    , name TEXT
    , domain TEXT
    , password TEXT
    , connection_id INT8
    , connection_method TEXT
    , external_data JSONB
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , status INT2
);

DROP INDEX tenant_owner_id_idx;

CREATE UNIQUE INDEX tenant_owner_id_connection_id_idx ON tenant(owner_id,connection_id) WHERE deleted_at IS NULL;

ALTER TABLE hotline
    ADD COLUMN tenant_id INT8 REFERENCES tenant(id);
