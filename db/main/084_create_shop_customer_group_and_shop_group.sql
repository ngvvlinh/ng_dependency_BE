create table shop_customer_group(
	id bigint primary key,
	name text,

	created_at timestamp with time zone,
	updated_at timestamp with time zone,
	deleted_at timestamp with time zone
);

create table shop_customer_group_customer (
	customer_id bigint  REFERENCES shop_customer(id),
	group_id bigint REFERENCES shop_customer_group(id),

	created_at timestamp with time zone,
	updated_at timestamp with time zone
);

alter table shop_customer_group_customer add constraint shop_customer_group_customer_constraint primary key (group_id, customer_id);

CREATE INDEX ON shop_customer_group (created_at);
CREATE INDEX ON shop_customer_group_customer (created_at);
