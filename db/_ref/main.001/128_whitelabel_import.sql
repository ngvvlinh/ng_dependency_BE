SELECT init_history('shop_brand', '{id, shop_id}');

ALTER TABLE shop_brand
    ADD COLUMN external_id TEXT,
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE history."shop_brand"
    ADD COLUMN external_id TEXT,
    ADD COLUMN partner_id INT8;

CREATE UNIQUE INDEX ON shop_brand(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;

ALTER TABLE shop_category
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_parent_id TEXT,
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE history.product_source_category
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_parent_id TEXT,
    ADD COLUMN partner_id INT8;

CREATE UNIQUE INDEX ON shop_category(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;

ALTER TABLE shop_product
    ADD COLUMN external_brand_id TEXT,
    ADD COLUMN external_category_id TEXT;

ALTER TABLE history.shop_product
    ADD COLUMN external_brand_id TEXT,
    ADD COLUMN external_category_id TEXT;

ALTER TABLE shop_variant
    ADD COLUMN external_product_id TEXT;

ALTER TABLE history.shop_variant
    ADD COLUMN external_product_id TEXT;

ALTER TABLE shop_collection
    ADD COLUMN partner_id INT8 REFERENCES partner(id),
    ADD COLUMN external_id TEXT;

ALTER TABLE history.shop_collection
    ADD COLUMN partner_id INT8,
    ADD COLUMN external_id TEXT;

CREATE UNIQUE INDEX ON shop_collection(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;

ALTER TABLE shop_product_collection
    ADD COLUMN partner_id INT8 REFERENCES partner(id),
    ADD COLUMN external_collection_id TEXT,
    ADD COLUMN external_product_id TEXT,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE history.product_shop_collection
    ADD COLUMN partner_id INT8,
    ADD COLUMN external_collection_id TEXT,
    ADD COLUMN external_product_id TEXT,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

CREATE UNIQUE INDEX ON shop_product_collection(partner_id, external_collection_id, external_product_id) WHERE partner_id IS NOT NULL AND external_collection_id IS NOT NULL AND external_product_id IS NOT NULL;