CREATE TABLE money_transaction_shipping_etop (
  id BIGINT PRIMARY KEY,
  code TEXT NOT NULL,
  total_cod INTEGER NOT NULL,
  total_orders INTEGER NOT NULL,
  total_amount INTEGER NOT NULL,
  total_fee INTEGER NOT NULL,
  total_money_transaction INTEGER NOT NULL, 
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  confirmed_at TIMESTAMP WITH TIME ZONE,
  status SMALLINT,
  bank_account JSONB,
  note TEXT,
  invoice_number TEXT
);

ALTER TABLE money_transaction_shipping 
  ADD COLUMN money_transaction_shipping_etop_id BIGINT REFERENCES money_transaction_shipping_etop(id),
  ADD COLUMN bank_account JSON,
  ADD COLUMN note TEXT,
  ADD COLUMN invoice_number TEXT;

ALTER TABLE history.money_transaction_shipping 
  ADD COLUMN money_transaction_shipping_etop_id BIGINT,
  ADD COLUMN bank_account JSON,
  ADD COLUMN note TEXT,
  ADD COLUMN invoice_number TEXT;

ALTER TYPE code_type ADD VALUE 'money_transaction_shipping_etop';

ALTER TABLE money_transaction_shipping_external
  ADD COLUMN bank_account JSON,
  ADD COLUMN note TEXT,
  ADD COLUMN invoice_number TEXT;
