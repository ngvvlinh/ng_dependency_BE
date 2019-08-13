ALTER TABLE product RENAME TO _product;
ALTER TABLE variant RENAME TO _variant;
ALTER TABLE etop_category RENAME TO _etop_category;
ALTER TABLE order_external RENAME TO _order_external;
ALTER TABLE order_source RENAME TO _order_source;
ALTER TABLE product_brand RENAME TO _product_brand;
ALTER TABLE product_external RENAME TO _product_external;
ALTER TABLE product_source RENAME TO _product_source;
ALTER TABLE product_source_category_external RENAME TO _product_source_category_external;
ALTER TABLE supplier RENAME TO _supplier;
ALTER TABLE variant_external RENAME TO _variant_external;

ALTER TABLE product_source_category RENAME to shop_category;

ALTER TABLE order_line DROP CONSTRAINT order_line_variant_id_fkey;
ALTER TABLE order_line DROP CONSTRAINT order_line_product_id_fkey;
ALTER TABLE order_line DROP CONSTRAINT order_line_supplier_id_fkey;
ALTER TABLE shop_variant DROP CONSTRAINT shop_variant_product_id_fkey;
ALTER TABLE product_shop_collection DROP CONSTRAINT product_shop_collection_product_id_fkey;
ALTER TABLE shop DROP CONSTRAINT shop_product_source_id_fkey;
ALTER TABLE shop DROP CONSTRAINT shop_order_source_id_fkey;
ALTER TABLE shop_category DROP CONSTRAINT product_source_category_product_source_id_fkey;
ALTER TABLE money_transaction_shipping DROP CONSTRAINT money_transaction_supplier_id_fkey;
ALTER TABLE fulfillment DROP CONSTRAINT fulfillment_supplier_id_fkey;
ALTER TABLE credit DROP CONSTRAINT credit_supplier_id_fkey;

ALTER TABLE _supplier DROP CONSTRAINT supplier_product_source_id_fkey;
ALTER TABLE _product DROP CONSTRAINT product_product_brand_id_fkey;
ALTER TABLE _product DROP CONSTRAINT product_supplier_id_fkey;
ALTER TABLE _product DROP CONSTRAINT product_etop_category_id_fkey;
ALTER TABLE _variant DROP CONSTRAINT variant_product_brand_id_fkey;
ALTER TABLE _variant DROP CONSTRAINT variant_supplier_id_fkey;
ALTER TABLE _variant DROP CONSTRAINT variant_etop_category_id_fkey;

DROP TABLE product_source_internal;
DROP TABLE _product_external;
DROP TABLE _variant_external;
DROP TABLE _product_brand;
DROP TABLE _product_source_category_external;
DROP TABLE _product_source;
DROP TABLE _supplier;
DROP TABLE _order_external;
DROP TABLE _order_source;
DROP TABLE _etop_category;

-- be careful, we want to keep them for safety!
--
-- DROP TABLE _variant;
-- DROP TABLE _product;

ALTER TABLE shop_product DROP CONSTRAINT shop_product_pkey;
ALTER TABLE shop_product ADD CONSTRAINT shop_product_pkey PRIMARY KEY (product_id);

ALTER TABLE shop_variant DROP CONSTRAINT shop_variant_pkey;
ALTER TABLE shop_variant ADD CONSTRAINT shop_variant_pkey PRIMARY KEY (variant_id);

CREATE UNIQUE INDEX ON shop_product (shop_id, product_id);
ALTER TABLE shop_variant ADD CONSTRAINT variant_product_id_fkey
    FOREIGN KEY (shop_id, product_id) REFERENCES shop_product (shop_id, product_id);

CREATE UNIQUE INDEX ON shop_product (shop_id, code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX ON shop_variant (shop_id, code) WHERE deleted_at IS NULL;
