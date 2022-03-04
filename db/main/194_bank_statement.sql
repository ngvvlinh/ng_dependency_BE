CREATE TABLE public."bank_statement" (
    id int8 NOT NULL PRIMARY KEY,
    account_id INT8 REFERENCES public."account" (id),
    amount INT4,
    description TEXT,
    transfered_at TIMESTAMPTZ,
    external_transaction_id TEXT UNIQUE ,
    sender_name TEXT,
    sender_bank_account TEXT,
    other_info JSONB,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON public."bank_statement"(external_transaction_id);

CREATE INDEX ON public."bank_statement"(account_id);

ALTER TABLE public."partner"
    ADD COLUMN whitelist_ips TEXT[];

ALTER TABLE history."partner"
    ADD COLUMN whitelist_ips TEXT[];

ALTER TABLE public."credit"
    ADD COLUMN bank_statement_id INT8 REFERENCES public."bank_statement" (id);

ALTER TABLE history."credit"
    ADD COLUMN bank_statement_id INT8;
