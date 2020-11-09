create table ticket_provider_webhook(
    id bigint primary key,
    ticket_provider text,
    connection_id bigint,
    external_status text,
    external_type text,
    client_id text,
    error jsonb,
    data jsonb,
    created_at timestamptz
);
