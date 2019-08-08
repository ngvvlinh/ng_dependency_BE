ALTER TYPE account_type ADD VALUE 'partner' AFTER 'shop';

CREATE TABLE partner (
  id INT8 PRIMARY KEY REFERENCES account(id)
, rid INT8
, name TEXT NOT NULL
, public_name TEXT NOT NULL
, owner_id INT8 NOT NULL REFERENCES "user"(id)
, status INT2 NOT NULL
, is_test INT2 NOT NULL
, phone TEXT
, email TEXT
, website_url TEXT
, image_url TEXT
, created_at TIMESTAMPTZ
, updated_at TIMESTAMPTZ
, deleted_at TIMESTAMPTZ
, contact_persons JSONB
);

CREATE TABLE account_auth (
  key TEXT PRIMARY KEY
, account_id INT8 NOT NULL REFERENCES account(id)
, status INT2
, roles TEXT[]
, permissions TEXT[]
, created_at TIMESTAMPTZ
, updated_at TIMESTAMPTZ
, deleted_at TIMESTAMPTZ
);

CREATE TYPE subject_type AS ENUM(
  'account'
, 'user'
);

CREATE TABLE partner_relation (
  partner_id INT8 NOT NULL REFERENCES partner(id)
, subject_id INT8 NOT NULL REFERENCES account(id)
, subject_type subject_type
, external_subject_id TEXT
, nonce INT8
, status INT2 NOT NULL
, roles TEXT[]
, permissions TEXT[]
, created_at TIMESTAMPTZ
, updated_at TIMESTAMPTZ
, deleted_at TIMESTAMPTZ
, PRIMARY KEY(partner_id, subject_id, subject_type)
);

CREATE INDEX ON account_auth(account_id);
CREATE INDEX ON partner_relation(subject_id);
CREATE INDEX ON partner_relation(partner_id);
CREATE UNIQUE INDEX ON partner_relation(partner_id, subject_id);

ALTER FUNCTION account_shop_supplier_update() RENAME TO update_to_account;

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON partner
  FOR EACH ROW EXECUTE PROCEDURE update_to_account();

SELECT init_history('partner', '{id}');
SELECT init_history('account_auth', '{key, account_id}');
SELECT init_history('partner_relation', '{partner_id, subject_id}');

ALTER TABLE "order" ADD COLUMN partner_id INT8 REFERENCES partner(id);
ALTER TABLE fulfillment ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE history."order" ADD COLUMN partner_id INT8;
ALTER TABLE history.fulfillment ADD COLUMN partner_id INT8;
