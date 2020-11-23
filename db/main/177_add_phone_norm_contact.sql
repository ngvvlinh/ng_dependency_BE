ALTER TABLE contact
    ADD COLUMN phone_norm tsvector;

ALTER TABLE history.contact
    ADD COLUMN phone_norm tsvector;
