ALTER TABLE shop_collection
    ADD COLUMN  deleted_at timestamp with time zone;
ALTER TABLE history.shop_collection
    ADD COLUMN  deleted_at timestamp with time zone;