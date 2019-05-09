CREATE TABLE callback (
  id INT8 PRIMARY KEY
, webhook_id INT8 NOT NULL
, account_id INT8 NOT NULL
, created_at TIMESTAMPTZ
, changes JSONB
, result JSONB
);
