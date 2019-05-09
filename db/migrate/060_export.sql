CREATE TABLE export_attempt (
  id TEXT PRIMARY KEY
, user_id INT8 REFERENCES "user"(id) NOT NULL
, account_id INT8 REFERENCES account(id) NOT NULL
, export_type TEXT NOT NULL
, filename TEXT
, stored_file TEXT
, download_url TEXT
, request_query TEXT
, mime_type TEXT
, status INT2 NOT NULL
, errors JSONB
, error JSONB
, n_total INT4
, n_exported INT4
, n_error INT4
, created_at TIMESTAMPTZ
, deleted_at TIMESTAMPTZ
, started_at TIMESTAMPTZ
, done_at TIMESTAMPTZ
, expires_at TIMESTAMPTZ
);

CREATE INDEX ON export_attempt(account_id);
