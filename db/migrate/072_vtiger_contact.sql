
CREATE	TABLE vtiger_contact (
	id text,
	firstname text,
	contact_no text,
	phone text,
	lastname text,
	mobile text,
	email text,
	email2 text,
	leadsource text,
	secondaryemail text,
	assigned_user_id text,
	created_at timestamptz,
	etop_user_id bigint,
	updated_at timestamptz,
	description text,
	"source" text,
	used_shipping_provider text,
	orders_per_day text,
	company text,
	city text,
	state text,
	website text,
	lane text,
	country text,
	search_norm tsvector,
	vtiger_created_at timestamptz,
	vtiger_updated_at timestamptz
)

CREATE INDEX search_idx ON vtiger_contact USING GIN(search_norm);
