create table money_transaction_shipping_external (
  id bigint primary key,
  code text not null,
  total_cod integer not null,
  total_orders integer not null,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  -- khi confirm => admin xác nhận ghn đã chuyển tiền + tạo phiên chuyển tiền bên etop
  status smallint,
  external_paid_at timestamp with time zone,
  provider shipping_provider  NOT NULL
);

create table money_transaction_shipping_external_line (
  id bigint primary key,
  external_code text not null,
  external_total_cod integer not null,
  external_created_at timestamp with time zone,
  external_closed_at timestamp with time zone,
  external_customer text,
  external_address text,
  etop_fulfillment_id_raw text,   -- raw data
  etop_fulfillment_id     bigint,  -- after parsing successfully
  note text,
  money_transaction_shipping_external_id bigint NOT NULL REFERENCES money_transaction_shipping_external(id),
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  import_error jsonb
);

create table money_transaction (
  id bigint primary key,
  shop_id bigint REFERENCES shop(id),
  supplier_id bigint REFERENCES supplier(id),
  status smallint,
  total_cod integer not null,
  total_orders integer not null,
  code text not null,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  closed_at timestamp with time zone,
  money_transaction_shipping_external_id bigint REFERENCES money_transaction_shipping_external(id),
  provider shipping_provider  NOT NULL,
  etop_paid_at  timestamp with time zone
);

alter table fulfillment
  add column if not exists money_transaction_id bigint REFERENCES money_transaction(id);

ALTER TYPE code_type ADD VALUE 'money_transaction' AFTER 'order';

ALTER table fulfillment
  add column if not exists cod_etop_paid_at timestamp with time zone,
  add column if not exists shipping_fee_shop_paid_at timestamp with time zone;
