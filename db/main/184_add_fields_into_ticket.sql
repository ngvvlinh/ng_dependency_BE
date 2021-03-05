ALTER TABLE ticket
    ADD COLUMN "type" INT;
ALTER TABLE "history".ticket
    ADD COLUMN "type" INT;

UPDATE ticket
SET "type" = 74;

ALTER TABLE ticket_label
    ADD COLUMN shop_id INT8 REFERENCES account(id),
    ADD COLUMN "type" INT;
ALTER TABLE "history".ticket_label
    ADD COLUMN shop_id INT8
    ADD COLUMN "type" INT;

UPDATE ticket_label
SET "type" = 74;
