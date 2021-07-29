CREATE INDEX account_user_roles_idx ON account_user USING GIN("roles");
