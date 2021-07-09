ALTER TABLE public."account_user"
    ADD COLUMN full_name_norm tsvector,
	ADD COLUMN phone_norm tsvector,
    ADD COLUMN extension_number_norm tsvector,
    ADD COLUMN phone TEXT;

ALTER TABLE history."account_user"
    ADD COLUMN full_name_norm tsvector,
	ADD COLUMN phone_norm tsvector,
    ADD COLUMN extension_number_norm tsvector,
	ADD COLUMN phone TEXT;

CREATE INDEX ON  public."account_user" USING GIN(full_name_norm);
CREATE INDEX ON  public."account_user" USING GIN(phone_norm);
CREATE INDEX ON  public."account_user" USING GIN(extension_number_norm);
