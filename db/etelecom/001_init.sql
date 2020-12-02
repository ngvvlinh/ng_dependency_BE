CREATE TABLE hotline (
    id INT8 PRIMARY KEY
    , user_id INT8
    , hotline TEXT
    , network TEXT
    , connection_id INT8
    , connection_method TEXT
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);

CREATE TABLE extension (
    id INT8 PRIMARY KEY
    , user_id INT8
    , account_id INT8
    , extension_number TEXT
    , extension_password TEXT
    , external_data JSONB
    , connection_id INT8
    , connection_method TEXT
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON extension(user_id, account_id);
CREATE UNIQUE INDEX ON extension(account_id, connection_id,  extension_number);

CREATE TABLE summary (
    id INT8 PRIMARY KEY
    , extension_id INT8 REFERENCES extension(id)
    , date DATE
    , total_phone_call INT2
    , total_call_time INT2
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
);


