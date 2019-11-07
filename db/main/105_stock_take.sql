CREATE TABLE shop_stocktake(
	id BIGINT primary key,
	shop_id BIGINT references shop(id),
	created_by BIGINT,
	updated_by BIGINT,
	created_at timestamptz,
	updated_at timestamptz,
	confirmed_at timestamptz,
	cancelled_at timestamptz,
	variant_ids INT8[],
	total_quantity INT4,
	status int4,
	lines jsonb,
	code text,
	code_norm INT4,
	note text
);
