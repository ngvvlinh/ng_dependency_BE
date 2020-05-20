create table credit (
  id bigint primary key,
  rid bigint,
  amount integer,
  shop_id bigint REFERENCES shop(id),
  supplier_id bigint REFERENCES supplier(id),
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  paid_at timestamp with time zone,
  type account_type,
  status smallint
);

alter table fulfillment add column if not exists external_shipping_returned_fee integer;
alter table history.fulfillment add column if not exists external_shipping_returned_fee integer;

alter table money_transaction
  add column total_amount integer;

SELECT init_history('credit', '{id,supplier_id,shop_id}', '{id}');
SELECT init_history('money_transaction', '{id,supplier_id,shop_id}', '{id}');
CREATE INDEX ON history.credit (shop_id);
CREATE INDEX ON history.money_transaction (shop_id);
