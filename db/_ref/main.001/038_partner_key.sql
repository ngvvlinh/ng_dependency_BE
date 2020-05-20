ALTER TABLE account_auth RENAME COLUMN key to auth_key;

ALTER TABLE history.account_auth RENAME COLUMN key to auth_key;

ALTER TABLE partner_relation ADD COLUMN auth_key TEXT NOT NULL;

ALTER TABLE history.partner_relation ADD COLUMN auth_key TEXT;

CREATE UNIQUE INDEX ON partner_relation(auth_key);
