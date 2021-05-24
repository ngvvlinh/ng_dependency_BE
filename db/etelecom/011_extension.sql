ALTER TABLE extension
    ADD COLUMN tenant_id INT8 REFERENCES tenant(id);

CREATE UNIQUE INDEX extension_extension_number_tenant_id_idx ON extension(extension_number,tenant_id) WHERE deleted_at is NULL;
