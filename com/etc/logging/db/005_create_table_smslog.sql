create table sms_log(
	id bigint primary key,
	external_id text,
	content text,
	phone text,
	provider text,
	status int2,
	error text,
    created_at timestamptz
);
