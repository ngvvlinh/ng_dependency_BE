DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'shop_ledger_type') THEN
        create type shop_ledger_type as enum (
            'cash',
            'bank'
        );
    END IF;
END
$$;

CREATE TABLE if not exists shop_ledger (
    id INT8 PRIMARY KEY
    , shop_id INT8
    , name TEXT
    , bank_account JSONB
    , status INT2
    , note TEXT
    , type shop_ledger_type
    , created_by INT8
    , created_at TIMESTAMP WITH TIME ZONE
    , updated_at TIMESTAMP WITH TIME ZONE
    , rid INT8
);