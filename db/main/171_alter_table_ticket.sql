ALTER TABLE ticket_comment
    DROP COLUMN image_url,
    ADD COLUMN image_urls TEXT[],
    ADD COLUMN created_name TEXT,
    ADD COLUMN created_source account_type;

ALTER TABLE ticket
    ADD COLUMN created_name TEXT,
    ADD COLUMN created_source account_type;

SELECT init_history('ticket_comment', '{id}');
SELECT init_history('ticket', '{id}');
