create table refund (
	id int8 not null primary key,
	shop_id int8 not null,
	order_id int8 not null,
	note text,
	discount int,
	code_norm int,
	code text,
	customer_id int8,
	lines jsonb,
	created_at timestamptz,
	updated_at timestamptz,
	cancelled_at timestamptz,
	confirmed_at timestamptz,
	created_by int8 not null,
	updated_by int8 not null,
	total_amount int,
	basket_value int,
	cancel_reason  text,
	status int
);

CREATE INDEX ON "refund" (shop_id, order_id);
CREATE INDEX ON "refund" (shop_id, id);
create INDEX ON "refund" (shop_id, code);

ALTER TYPE receipt_ref_type ADD VALUE 'refund';
