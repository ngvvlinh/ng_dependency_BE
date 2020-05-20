create table purchase_refund (
	id int8 not null primary key,
	shop_id int8 not null,
	purchase_order_id int8 not null,
	note text,
	discount int,
	code_norm int,
	code text,
	supplier_id int8,
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

create INDEX on "purchase_refund" (shop_id, purchase_order_id);
create INDEX on "purchase_refund" (shop_id, id);
create INDEX on "purchase_refund" (shop_id, code);

alter type receipt_ref_type add value 'purchase_refund';

alter table inventory_voucher add column "rollback" bool;
