ALTER TABLE product ADD COLUMN name_norm_ua text;     -- unaccent normalization
ALTER TABLE variant ADD COLUMN attr_norm_kv tsvector; -- key-value normalization

ALTER TABLE history.product ADD COLUMN name_norm_ua text;
ALTER TABLE history.variant ADD COLUMN attr_norm_kv tsvector;

-- execute scripts/update_norm and clean up duplicated values before adding
-- following constraints

-- this column may be empty but not null, for simplify while querying
ALTER TABLE variant ALTER COLUMN attr_norm_kv SET NOT NULL;

CREATE UNIQUE INDEX product_product_source_id_name_norm_ua_idx ON product(product_source_id, name_norm_ua)
  WHERE deleted_at IS NULL AND ed_code IS NULL;

CREATE UNIQUE INDEX variant_product_id_attr_norm_kv_idx ON variant(product_id, attr_norm_kv) WHERE deleted_at IS NULL;
