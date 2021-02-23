ALTER TABLE extension
    ADD COLUMN subscription_id INT8
    , ADD COLUMN expires_at TIMESTAMPTZ;

DROP INDEX extension_user_id_account_id_idx;
DROP INDEX extension_tenant_domain_extension_number_idx;

CREATE UNIQUE INDEX extension_user_id_account_id_idx ON extension(user_id, account_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX extension_tenant_domain_extension_number_idx ON extension(tenant_domain, extension_number) WHERE deleted_at IS NULL;
