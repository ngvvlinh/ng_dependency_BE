CREATE TABLE subscription_product (
    id INT8 PRIMARY KEY
    , name TEXT
    , type TEXT
    , description TEXT
    , image_url TEXT
    , status INT2
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE TABLE subscription_plan (
    id INT8 PRIMARY KEY
    , name TEXT
    , price INT8
    , status INT2
    , description TEXT
    , product_id INT8 REFERENCES subscription_product(id)
    , interval TEXT
    , interval_count INT8
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE TABLE subscription (
    id INT8 PRIMARY KEY
    , account_id INT8 REFERENCES account(id)
    , cancel_at_period_end BOOLEAN
    , current_period_end_at TIMESTAMPTZ
    , current_period_start_at TIMESTAMPTZ
    , status INT2
    , billing_cycle_anchor_at TIMESTAMPTZ
    , started_at TIMESTAMPTZ
    , customer JSONB
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
    , deleted_at TIMESTAMPTZ
    , wl_partner_id INT8 REFERENCES partner(id)
);

CREATE TABLE subscription_line (
    id INT8 PRIMARY KEY
    , plan_id INT8 REFERENCES subscription_plan(id)
    , subscription_id INT8 REFERENCES subscription(id)
    , quantity INT8
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
);

CREATE TABLE subscription_bill (
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
);

CREATE TABLE subscription_bill_line (
    id INT8 PRIMARY KEY
    , line_amount INT8
    , price INT8
    , quantity INT8
    , description TEXT
    , period_start_at TIMESTAMPTZ
    , period_end_at TIMESTAMPTZ
    , subscription_bill_id INT8 REFERENCES subscription_bill(id)
    , subscription_id INT8 REFERENCES subscription(id)
    , created_at TIMESTAMPTZ
    , updated_at TIMESTAMPTZ
);

SELECT init_history('subscription_product', '{id}');
SELECT init_history('subscription_plan', '{id}');
SELECT init_history('subscription', '{id,account_id}');
SELECT init_history('subscription_line', '{id,subscription_id}');
SELECT init_history('subscription_bill', '{id,account_id}');
SELECT init_history('subscription_bill_line', '{id,subscription_bill_id,subscription_id}');

ALTER TABLE subscription ADD COLUMN plan_ids INT8[];
ALTER TABLE history.subscription ADD COLUMN plan_ids INT8[];

CREATE UNIQUE INDEX ON subscription_product(type);
