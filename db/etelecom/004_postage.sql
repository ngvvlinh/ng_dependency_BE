ALTER TABLE hotline
    ADD COLUMN is_free_charge bool;

ALTER TABLE call_log
    ADD COLUMN postage int
    , ADD COLUMN duration_postage int;

DROP TABLE summary;

CREATE INDEX ON call_log(account_id, call_status);
