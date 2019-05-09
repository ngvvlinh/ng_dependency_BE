ALTER TABLE shipping_source ADD COLUMN username TEXT;

DELETE FROM shipping_source_internal where id in (
	SELECT id from shipping_source where type = 'vtpost'
);
DELETE FROM shipping_source WHERE type = 'vtpost';

DROP INDEX shipping_source_name_type_idx;
CREATE UNIQUE INDEX ON shipping_source (name, type, username);
