CREATE TYPE user_identifying_type AS ENUM('full', 'half', 'stub');

ALTER TABLE "user"
  ALTER COLUMN full_name DROP NOT NULL,
  ALTER COLUMN short_name DROP NOT NULL,
  ALTER COLUMN email DROP NOT NULL,
  ALTER COLUMN phone DROP NOT NULL,
  ADD COLUMN identifying user_identifying_type;
ALTER TABLE history."user" ADD COLUMN identifying user_identifying_type;

UPDATE "user" SET identifying = 'full' WHERE identifying IS NULL;

ALTER TABLE "user" ALTER COLUMN identifying SET NOT NULL;

ALTER TABLE "user"
  ADD CONSTRAINT user_identifying CHECK
    ((identifying = 'full' AND status != 0 AND  phone IS NOT NULL AND email IS NOT NULL) OR
     (identifying = 'half' AND status != 0 AND  phone IS NOT NULL) OR
     (identifying = 'stub' AND status  = 0 AND (phone IS NOT NULL  OR email IS NOT NULL)));

ALTER TABLE "user"
  ADD CONSTRAINT user_not_full CHECK
     (identifying = 'stub' OR (full_name IS NOT NULL AND short_name IS NOT NULL));

ALTER TABLE "account_user"
  ADD COLUMN full_name TEXT,
  ADD COLUMN short_name TEXT,
  ADD COLUMN position TEXT,
  ADD COLUMN response_status INT2,
  ADD COLUMN invitation_sent_at TIMESTAMPTZ,
  ADD COLUMN invitation_sent_by INT8 REFERENCES "user" (id),
  ADD COLUMN invitation_accepted_at TIMESTAMPTZ,
  ADD COLUMN invitation_rejected_at TIMESTAMPTZ,
  ADD COLUMN disabled_at TIMESTAMPTZ,
  ADD COLUMN disabled_by INT8 REFERENCES "user" (id),
  ADD COLUMN disable_reason TEXT;

ALTER TABLE history."account_user"
  ADD COLUMN full_name TEXT,
  ADD COLUMN short_name TEXT,
  ADD COLUMN position TEXT,
  ADD COLUMN response_status INT2,
  ADD COLUMN invitation_sent_at TIMESTAMPTZ,
  ADD COLUMN invitation_sent_by INT8,
  ADD COLUMN invitation_accepted_at TIMESTAMPTZ,
  ADD COLUMN invitation_rejected_at TIMESTAMPTZ,
  ADD COLUMN disabled_at TIMESTAMPTZ,
  ADD COLUMN disabled_by INT8,
  ADD COLUMN disable_reason TEXT;

ALTER TABLE account ADD COLUMN owner_id INT8 REFERENCES "user"(id);
ALTER TABLE history.account ADD COLUMN owner_id INT8;

CREATE OR REPLACE FUNCTION account_shop_supplier_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE account SET
      name = NEW.name,
      owner_id = NEW.owner_id,
      image_url = NEW.image_url,
      deleted_at = NEW.deleted_at
    WHERE id = NEW.id;
    RETURN NEW;
END;
$$;

-- owner_id of supplier and etop may be null
UPDATE account SET owner_id = shop.owner_id FROM shop WHERE account.id = shop.id AND account.owner_id IS NULL;
