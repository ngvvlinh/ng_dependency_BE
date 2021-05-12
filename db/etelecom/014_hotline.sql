CREATE UNIQUE INDEX hotline_hotline_idx ON hotline(hotline) WHERE deleted_at is NULL;
