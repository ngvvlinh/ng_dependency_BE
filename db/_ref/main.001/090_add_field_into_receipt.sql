CREATE TYPE receipt_created_type AS enum('manual', 'auto');

ALTER TABLE receipt ADD COLUMN created_type receipt_created_type;
ALTER TABLE "history".receipt ADD COLUMN created_type receipt_created_type;

ALTER TABLE receipt ADD COLUMN ledger_id INT8 references shop_ledger(id);
ALTER TABLE "history".receipt ADD COLUMN ledger_id INT8;

CREATE INDEX ON receipt(ledger_id);
