create table if not exists refund (
	id int8 primary key,
	shop_id int8,
	order_id int8,
	note text,
	code_norm int,
	code text,
	customer_id int8,
	lines jsonb,
	created_at timestamptz,
	updated_at timestamptz,
	cancelled_at timestamptz,
	confirmed_at timestamptz,
	created_by int8,
	updated_by int8,
	total_amount int,
	basket_value int,
	cancel_reason text,
	adjustment_lines jsonb,
	total_adjustment int,
	status int,
	rid int8
);

alter table purchase_refund
    drop column if exists adjustment_lines,
    drop column if exists total_adjustment;
