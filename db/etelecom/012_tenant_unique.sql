CREATE UNIQUE INDEX tenant_connection_id_connection_method_idx ON tenant(connection_id, connection_method) WHERE connection_method = 'builtin';
