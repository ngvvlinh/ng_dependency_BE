CREATE UNIQUE INDEX ON extension(tenant_domain, extension_number);

ALTER TABLE call_log
    ADD COLUMN external_session_id TEXT UNIQUE;
