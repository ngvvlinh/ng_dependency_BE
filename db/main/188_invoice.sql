ALTER TABLE invoice
    ADD COLUMN classify TEXT
    , ADD COLUMN type TEXT;

ALTER TABLE history.invoice
    ADD COLUMN classify TEXT
    , ADD COLUMN type TEXT;
