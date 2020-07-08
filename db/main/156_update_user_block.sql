alter table "user"
  add column blocked_by int8
  , add column blocked_at timestamptz
  , add column block_reason text;

 alter table history."user"
  add column blocked_by int8
  , add column blocked_at timestamptz
  , add column block_reason text;
