CREATE TYPE import_type AS ENUM('other', 'shop_order', 'shop_product', 'etop_money_transaction');

CREATE TABLE import_attempt (
  id             INT8 PRIMARY KEY,
  user_id        INT8 NOT NULL REFERENCES "user"(id),
  account_id     INT8 NOT NULL REFERENCES account(id),
  original_file  TEXT NOT NULL,
  stored_file    TEXT,
  type           import_type NOT NULL,
  n_created      INT2 NOT NULL,
  n_updated      INT2 NOT NULL,
  n_error        INT2 NOT NULL,
  status         INT2 NOT NULL,
  error_type     TEXT,
  errors         JSONB,
  duration_ms    INT4 NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL
);

SELECT init_history('import_attempt', '{id,user_id,account_id}');
