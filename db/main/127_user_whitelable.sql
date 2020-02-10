ALTER TABLE "user" ADD COLUMN wl_partner_id INT8 REFERENCES partner(id);
ALTER TABLE history."user" ADD COLUMN wl_partner_id INT8;

CREATE UNIQUE INDEX user_email_wl_partner_id_idx ON "user" (email, wl_partner_id);

DROP INDEX IF EXISTS user_phone_key;
DROP INDEX IF EXISTS user_email_key;

CREATE UNIQUE INDEX user_phone_key ON "user" (phone) where wl_partner_id IS NULL;
CREATE UNIQUE INDEX user_email_key ON "user" (email) where wl_partner_id IS NULL;
