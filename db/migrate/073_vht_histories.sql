create table vht_call_history (
	cdr_id text,
	call_id text,
	sip_call_id text,
	sdk_call_id text,
	cause text,
	q850_cause text,
	from_extension text,
	to_extension text,
	from_number text,
	"to_number" text,
	duration BIGINT,
	direction BIGINT,
	time_started timestamptz,
	time_connected timestamptz,
	time_ended timestamptz,
	recording_path text,
	recording_url text,
	record_file_size BIGINT,
	etop_account_id text,
	vtiger_account_id text,
	sync_status text,
	created_at timestamptz,
	updated_at timestamptz,
	o_data text,
	search_norm tsvector
)

create index search_idx_vht ON vht_call_history USING GIN(search_norm);
