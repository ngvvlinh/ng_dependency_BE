ALTER TABLE shop_product
  ADD COLUMN deleted_at TIMESTAMPTZ
, ADD COLUMN code TEXT
, ADD COLUMN name_norm_ua TEXT
, ADD COLUMN category_id INT8
, ADD COLUMN cost_price INT4
, ADD COLUMN list_price INT4
, ADD COLUMN unit TEXT
;

ALTER TABLE history.shop_product
  ADD COLUMN deleted_at TIMESTAMPTZ
, ADD COLUMN code TEXT
, ADD COLUMN name_norm_ua TEXT
, ADD COLUMN category_id INT8
, ADD COLUMN cost_price INT4
, ADD COLUMN list_price INT4
, ADD COLUMN unit TEXT
;

ALTER TABLE shop_variant
  ADD COLUMN deleted_at TIMESTAMPTZ
, ADD COLUMN attr_norm_kv TSVECTOR
, ADD COLUMN code TEXT
, ADD COLUMN cost_price INT4
, ADD COLUMN list_price INT4
, ADD COLUMN attributes JSONB
;

ALTER TABLE history.shop_variant
  ADD COLUMN deleted_at TIMESTAMPTZ
, ADD COLUMN attr_norm_kv TSVECTOR
, ADD COLUMN code TEXT
, ADD COLUMN cost_price INT4
, ADD COLUMN list_price INT4
, ADD COLUMN attributes JSONB
;

DROP TRIGGER save_history ON shop_product;
DROP TRIGGER save_history ON shop_variant;
DROP TRIGGER save_history ON product;
DROP TRIGGER save_history ON variant;

UPDATE shop_product AS sp
SET
  deleted_at   = COALESCE (sp.deleted_at, p.deleted_at)
, code         = COALESCE (sp.code, p.code)
, name         = COALESCE (sp.name, p.name)
, description  = COALESCE (sp.description, p.description)
, desc_html    = COALESCE (sp.desc_html, p.desc_html)
, short_desc   = COALESCE (sp.short_desc, p.short_desc)
, image_urls   = COALESCE (sp.image_urls, p.image_urls)
, unit         = COALESCE (sp.unit, p.unit)
, name_norm    = COALESCE (sp.name_norm, p.name_norm)
, name_norm_ua = COALESCE (sp.name_norm_ua, p.name_norm_ua)
FROM product AS p
WHERE sp.product_id = p.id
;

UPDATE shop_variant AS sv
SET
  deleted_at   = COALESCE (sv.deleted_at, v.deleted_at)
, code         = COALESCE (sv.code, v.code)
, name         = COALESCE (sv.name, v.name)
, description  = COALESCE (sv.description, v.description)
, desc_html    = COALESCE (sv.desc_html, v.desc_html)
, short_desc   = COALESCE (sv.short_desc, v.short_desc)
, image_urls   = COALESCE (sv.image_urls, v.image_urls)
, cost_price   = COALESCE (sv.cost_price, v.cost_price)
, list_price   = COALESCE (sv.list_price, v.list_price)
, attributes   = COALESCE (sv.attributes, v.attributes)
, name_norm    = COALESCE (sv.name_norm, v.name_norm)
, attr_norm_kv = COALESCE (sv.attr_norm_kv, v.attr_norm_kv)
FROM variant as v
WHERE sv.variant_id = v.id
;

ALTER TABLE product DROP CONSTRAINT product_product_source_id_fkey;
ALTER TABLE variant DROP CONSTRAINT variant_product_source_id_fkey;

UPDATE product
SET product_source_id = shop.id
FROM shop
WHERE shop.product_source_id = product.product_source_id;

INSERT INTO shop_product (
  shop_id
, product_id
, code
, name
, description
, desc_html
, short_desc
, image_urls
, unit
, cost_price
, list_price
, retail_price
, status
, created_at
, updated_at
, deleted_at
, name_norm
, name_norm_ua
)
SELECT
  p.product_source_id
, p.id
, p.code
, p.name
, p.description
, p.desc_html
, p.short_desc
, p.image_urls
, p.unit
, 0
, 0
, 0
, 1 -- status
, p.created_at
, p.updated_at
, p.deleted_at
, p.name_norm
, p.name_norm_ua
FROM product AS p
LEFT JOIN shop_product AS sp
ON p.id = sp.product_id
WHERE sp.product_id IS NULL
;

INSERT INTO shop_variant (
  shop_id
, variant_id
, product_id
, code
, name
, description
, desc_html
, short_desc
, image_urls
, cost_price
, list_price
, retail_price
, status
, attributes
, created_at
, updated_at
, deleted_at
, name_norm
, attr_norm_kv
)
SELECT
  p.product_source_id
, v.id
, v.product_id
, v.code
, v.name
, v.description
, v.desc_html
, v.short_desc
, v.image_urls
, v.cost_price
, v.list_price
, 0
, 1 -- status
, v.attributes
, v.created_at
, v.updated_at
, v.deleted_at
, v.name_norm
, v.attr_norm_kv
FROM variant AS v
LEFT JOIN product AS p ON v.product_id = p.id
LEFT JOIN shop_variant AS sv ON v.id = sv.variant_id
WHERE sv.variant_id IS NULL
;

CREATE TRIGGER save_history
  AFTER INSERT OR UPDATE OR DELETE ON shop_product
  FOR EACH ROW EXECUTE PROCEDURE save_history('shop_product', '{shop_id,product_id}');

CREATE TRIGGER save_history
  AFTER INSERT OR UPDATE OR DELETE ON shop_variant
  FOR EACH ROW EXECUTE PROCEDURE save_history('shop_variant', '{shop_id,variant_id}');

ALTER TABLE shop_product DROP CONSTRAINT shop_product_product_id_fkey;
ALTER TABLE shop_variant DROP CONSTRAINT shop_variant_variant_id_fkey;
