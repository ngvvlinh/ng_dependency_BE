alter table fulfillment add column if not exists last_sync_at timestamp with time zone;
alter table history.fulfillment add column if not exists last_sync_at timestamp with time zone;

create table shipping_source (
  id bigint primary key,
  name text,
  type shipping_provider not null,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  rid bigint
);

create table shipping_source_state (
  id bigint primary key,
  rid int8,
  last_sync_at timestamp with time zone,
  created_at timestamp with time zone,
  updated_at timestamp with time zone
);
insert into shipping_source (id, type) values (1038819005228305682, 'ghn');
insert into shipping_source_state (id) values (1038819005228305682);

CREATE UNIQUE INDEX ON fulfillment (shipping_provider, external_shipping_code) where status not in (-1, 1);
CREATE UNIQUE INDEX ON fulfillment (shipping_provider, shipping_code) where status not in (-1, 1);
