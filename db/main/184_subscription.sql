ALTER TABLE transaction
    ADD COLUMN classify TEXT
    , ADD COLUMN name TEXT
    , ADD COLUMN referral_type TEXT
    , ADD COLUMN referral_ids INT8[];

ALTER TABLE history.transaction
    ADD COLUMN classify TEXT
    , ADD COLUMN name TEXT
    , ADD COLUMN referral_type TEXT
    , ADD COLUMN referral_ids INT8[];

ALTER TABLE transaction
    DROP COLUMN metadata;

ALTER TABLE history.transaction
    DROP COLUMN metadata;

CREATE TABLE invoice (
    id INT8 PRIMARY KEY
    , account_id INT8 REFERENCES account(id)
    , subscription_id INT8 REFERENCES subscription(id)
    , total_amount INT8
    , description TEXT
    , payment_id INT8 REFERENCES payment(id)
    , payment_status INT2
    , status INT2
    , customer JSONB
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
    , referral_type TEXT
    , referral_ids INT8[]
);

CREATE TABLE invoice_line (
    id INT8 PRIMARY KEY
    , line_amount INT8
    , price INT8
    , quantity INT8
    , description TEXT
    , invoice_id INT8 REFERENCES invoice(id)
    , referral_type TEXT
    , referral_id INT8
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
);

SELECT init_history('invoice', '{id,account_id}');
SELECT init_history('invoice_line', '{id,invoice_id}');

-- delete later
ALTER TABLE subscription_bill
    RENAME TO _subscription_bill;
ALTER TABLE subscription_bill_line
    RENAME TO _subscription_bill_line;
