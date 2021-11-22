CREATE TABLE public. "department" (
    id bigint NOT NULL PRIMARY KEY,
    account_id BIGINT REFERENCES public. "account" (id),
    name text,
    description text,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE public. "account_user"
    ADD COLUMN department_id bigint REFERENCES public. "department" (id);

ALTER TABLE history. "account_user"
    ADD COLUMN department_id bigint;
