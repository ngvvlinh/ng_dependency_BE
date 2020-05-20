ALTER TABLE shop_customer
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE "history".shop_customer
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8;

CREATE UNIQUE INDEX ON shop_customer(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;
CREATE UNIQUE INDEX ON shop_customer(partner_id, shop_id, external_code) WHERE partner_id IS NOT NULL AND external_code IS NOT NULL;
CREATE INDEX ON shop_customer(updated_at, id);

ALTER TABLE shop_product
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE "history".shop_product
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8;

CREATE UNIQUE INDEX ON shop_product(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;
CREATE UNIQUE INDEX ON shop_product(partner_id, shop_id, external_code) WHERE partner_id IS NOT NULL AND external_code IS NOT NULL;
CREATE INDEX ON shop_product(updated_at, product_id);

ALTER TABLE shop_variant
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8 REFERENCES partner(id);

ALTER TABLE "history".shop_variant
    ADD COLUMN external_id TEXT,
    ADD COLUMN external_code TEXT,
    ADD COLUMN partner_id INT8;

CREATE UNIQUE INDEX ON shop_variant(partner_id, external_id) WHERE partner_id IS NOT NULL AND external_id IS NOT NULL;
CREATE UNIQUE INDEX ON shop_variant(partner_id, shop_id, external_code) WHERE partner_id IS NOT NULL AND external_code IS NOT NULL;
CREATE INDEX ON shop_variant(updated_at, variant_id);