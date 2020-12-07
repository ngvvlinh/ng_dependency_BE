CREATE TABLE call_log(
    id INT8 PRIMARY KEY,
    external_id TEXT,
    account_id INT8,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    duration INT,
    caller TEXT,
    callee TEXT,
    audio_urls TEXT[],
    external_direction TEXT,
    direction INT,
    external_call_status TEXT,
    extension_id INT8 REFERENCES extension(id),
    hotline_id INT8 REFERENCES hotline(id),
    contact_id INT8,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    call_state TEXT,
    call_status INT2
);

CREATE UNIQUE INDEX "call_log_extension_id_external_id_direction_key" ON call_log(extension_id, external_id, direction);

CREATE UNIQUE INDEX "call_log_extension_id_external_id_key" ON call_log(extension_id, external_id) WHERE direction IS NULL;

ALTER TABLE extension
    DROP COLUMN connection_id,
    DROP COLUMN connection_method;

CREATE INDEX ON call_log(started_at, id);

ALTER TABLE hotline
    ADD COLUMN name TEXT
    , ADD COLUMN description TEXT
    , ADD COLUMN status INT2;
