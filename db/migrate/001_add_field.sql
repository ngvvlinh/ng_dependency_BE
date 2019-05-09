ALTER TABLE variant ADD COLUMN ed_status int2;
ALTER TABLE product ADD COLUMN ed_tags text[];
ALTER TABLE product_external ADD COLUMN external_units text[];
ALTER TABLE product_external DROP COLUMN IF EXISTS external_attributes;
ALTER TABLE variant_external
  ADD COLUMN external_units text[],
  ADD COLUMN product_id bigint;

ALTER TABLE product_source
  ADD COLUMN IF NOT EXISTS last_sync_at timestamp with time zone,
  ADD COLUMN IF NOT EXISTS sync_state_products jsonb,
  ADD COLUMN IF NOT EXISTS sync_state_categories jsonb,
  ADD COLUMN IF NOT EXISTS external_status int2,
  ADD COLUMN IF NOT EXISTS supplier_id INT8 UNIQUE REFERENCES supplier(id);

ALTER TABLE order_source
  ADD COLUMN IF NOT EXISTS external_status int2,
  ADD COLUMN IF NOT EXISTS shop_id int8 NOT NULL UNIQUE REFERENCES shop(id);

ALTER TABLE "order"
  ADD COLUMN IF NOT EXISTS external_order_id text,
  ADD COLUMN IF NOT EXISTS order_source_id int8,
  ADD COLUMN IF NOT EXISTS order_source_type order_source_type,
  ADD COLUMN IF NOT EXISTS external_order_source text,
  ADD COLUMN IF NOT EXISTS shop_address jsonb,
  ADD COLUMN IF NOT EXISTS fulfillment_status int2,
  ADD COLUMN IF NOT EXISTS payment_status int2,
  DROP COLUMN IF EXISTS processing_status;

ALTER TABLE order_external
  DROP COLUMN IF EXISTS external_order_source_type,
  DROP COLUMN IF EXISTS external_shop_id;

CREATE TABLE shop_product (
    shop_id bigint NOT NULL REFERENCES shop(id),
    product_id bigint NOT NULL REFERENCES product(id),
    rid bigint,
    collection_id bigint REFERENCES shop_collection(id),
    name text,
    description text,
    desc_html text,
    short_desc text,
    retail_price integer,
    tags text[],
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    note text,
    image_urls text[],
    status smallint,
    haravan_id text,
    name_norm tsvector,
    PRIMARY KEY (shop_id, product_id)
);

ALTER TABLE shop_variant
  ADD COLUMN product_id bigint REFERENCES product(id);

alter table product_source alter supplier_id drop not null;

alter table variant add column attributes jsonb,
  add column cost_price integer;

-- supplier -> shop     : self
-- supplier -> customer : haravan
-- shop     -> customer : etop_pos, etop_pxs

CREATE TYPE fulfillment_endpoint AS ENUM (
    'supplier',
    'shop',
    'customer'
);

CREATE TYPE shipping_provider AS ENUM ('ghn');

CREATE TABLE fulfillment (
    id bigint NOT NULL,
    order_id bigint references "order"(id) not null,
    lines jsonb not null,
    variant_ids bigint[] not null,
    type_from fulfillment_endpoint not null,
    type_to   fulfillment_endpoint not null,
    address_from jsonb not null,
    address_to   jsonb not null,
    shop_id bigint references shop(id) not null,
    supplier_id bigint references supplier(id),
    supplier_confirm smallint,
    shop_confirm smallint,
    total_items integer,
    total_weight integer,
    total_amount integer,
    total_cod_amount integer,
    shipping_fee_customer integer,
    shipping_fee_shop integer,
    external_shipping_fee integer,
    created_at timestamp with time zone,
--     created_by bigint,
    updated_at timestamp with time zone,
    delivered_at timestamp with time zone,
    returned_at timestamp with time zone,
--     updated_by bigint,
    cancelled_at timestamp with time zone,
    cancel_reason text,
--     cancelled_by bigint,
    closed_at timestamp with time zone,
    shipping_provider shipping_provider,
    shipping_code text,
    shipping_note text,
    external_shipping_id text not null,
    external_shipping_code text not null,
    external_shipping_service_id text not null,
    external_shipping_created_at timestamp with time zone,
    external_shipping_updated_at timestamp with time zone,
    external_shipping_cancelled_at timestamp with time zone,
    external_shipping_delivered_at timestamp with time zone,
    external_shipping_returned_at timestamp with time zone,
    external_shipping_closed_at timestamp with time zone,
    external_shipping_state text,
    external_shipping_status smallint,
    external_shipping_data jsonb,
    shipping_state shipping_state,
    status smallint,
    sync_status smallint,
    sync_states jsonb,
    rid bigint
);

CREATE OR REPLACE FUNCTION ffm_active_supplier(supplier_id INT8, status INT2) RETURNS INT8
LANGUAGE plpgsql IMMUTABLE AS $$
BEGIN
  IF (status = -1) THEN RETURN NULL; END IF;
  RETURN supplier_id;
END;
$$;

CREATE UNIQUE INDEX ffm_active_supplier_key ON fulfillment (order_id, ffm_active_supplier(supplier_id, status));

CREATE OR REPLACE FUNCTION ids_not_empty(ids INT8[]) RETURNS BOOLEAN
LANGUAGE plpgsql IMMUTABLE AS $$
BEGIN
  RETURN ((array_length(ids, 1) is null) != (0 != ALL(ids)));
END;
$$;

-- ALTER TABLE "order" ADD CONSTRAINT ids_not_empty CHECK (
--   ids_not_empty(product_ids) AND
--   ids_not_empty(variant_ids) AND
--   ids_not_empty(supplier_ids)
-- );

-- ALTER TABLE fulfillment ADD CONSTRAINT ids_not_empty CHECK (
--   ids_not_empty(variant_ids)
-- );

ALTER TABLE fulfillment ADD CONSTRAINT type_from_supplier_id CHECK (
  (type_from = 'supplier') = (supplier_id IS NOT NULL)
);

ALTER TABLE "order"
	ADD COLUMN IF NOT EXISTS shipping_fee INTEGER,
	ADD COLUMN IF NOT EXISTS reference_url TEXT;

ALTER TABLE "order"
	ADD COLUMN IF NOT EXISTS fulfillment_status int2,
	ADD COLUMN IF NOT EXISTS payment_status int2,
  ADD COLUMN IF NOT EXISTS shop_cod integer;

CREATE OR REPLACE FUNCTION variant_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    supplier_status INT2;
    external_status INT2;
    etop_status     INT2;
		product_status 	INT2;
BEGIN
    SELECT ve.external_status INTO external_status FROM variant_external as ve WHERE id = NEW.id;
    IF TG_OP = 'UPDATE' THEN
        supplier_status := COALESCE(NEW.ed_status, OLD.ed_status, 0);
        etop_status     := COALESCE(NEW.etop_status, OLD.etop_status, 0);
    ELSE
        supplier_status := COALESCE(NEW.ed_status, 0);
        etop_status     := COALESCE(NEW.etop_status,     0);
    END IF;

    NEW.status := LEAST(supplier_status, external_status, etop_status);

		SELECT MAX(v.status) INTO product_status FROM variant as v WHERE product_id = NEW.product_id AND id != NEW.id;
		product_status := GREATEST(product_status, NEW.status);
		UPDATE product SET status = product_status WHERE id = NEW.product_id;
    RETURN NEW;
END
$$;

CREATE TRIGGER variant_update BEFORE INSERT OR UPDATE ON variant FOR EACH ROW EXECUTE PROCEDURE variant_update();


CREATE OR REPLACE FUNCTION variant_external_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    supplier_status INT2;
    x_status INT2;
    etop_status     INT2;
    final_status    INT2;
BEGIN
    SELECT v.ed_status, v.etop_status INTO supplier_status, etop_status FROM variant as v WHERE id = NEW.id;
    IF TG_OP = 'UPDATE' THEN
        x_status := COALESCE(NEW.external_status, OLD.external_status, 0);
    ELSE
        x_status := COALESCE(NEW.external_status, 0);
    END IF;

    final_status := LEAST(supplier_status, x_status, etop_status);
    UPDATE variant SET status = final_status where id = NEW.id;
    RETURN NEW;
END
$$;

CREATE TRIGGER variant_external_update AFTER INSERT OR UPDATE ON variant_external FOR EACH ROW EXECUTE PROCEDURE variant_external_update();

ALTER TABLE address RENAME COLUMN address TO address1;
ALTER TABLE address
    ADD COLUMN full_name text,
    ADD COLUMN first_name text,
    ADD COLUMN last_name text,
    ADD COLUMN email text,
    ADD COLUMN position text,
    ADD COLUMN city text,
    ADD COLUMN zip text,
    ADD COLUMN address2 text,
    ADD COLUMN phone text,
    ADD COLUMN company text;

ALTER TYPE address_type ADD VALUE 'shipto' AFTER 'warehouse';
ALTER TYPE address_type ADD VALUE 'shipfrom' AFTER 'warehouse';

ALTER TABLE shop
    ADD COLUMN ship_to_address_id bigint REFERENCES address(id),
    ADD COLUMN ship_from_address_id bigint REFERENCES address(id);

ALTER TABLE supplier
    ADD COLUMN ship_from_address_id bigint REFERENCES address(id);

ALTER TABLE address
    ADD COLUMN notes jsonb;

ALTER TABLE "order" RENAME COLUMN shipping_fee TO shop_shipping_fee;
ALTER TABLE "order" ADD COLUMN shop_shipping jsonb;

ALTER TABLE product ADD COLUMN subcode text;

alter table product_source_category ADD COLUMN shop_id BIGINT REFERENCES shop(id);

CREATE TABLE product_shop_collection (
    product_id bigint references product(id),
    shop_id bigint references shop(id),
    collection_id bigint references shop_collection(id),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    status smallint,
    primary key (product_id, collection_id)
);

-- Product Name Normalize -- use for searching --
-- CREATE EXTENSION unaccent;

ALTER TABLE product ADD COLUMN name_norm tsvector;
-- UPDATE product SET name_norm = to_tsvector('simple', unaccent(product.name));
CREATE INDEX product_search_idx ON product USING gin(name_norm);
-- CREATE OR REPLACE FUNCTION update_name_norm_product() RETURNS trigger
--     LANGUAGE plpgsql
--     AS $$
-- BEGIN
-- 		if NEW.name != '' then
--             NEW.name_norm = to_tsvector('simple', unaccent(NEW.name));
-- 		END IF;
--     RETURN NEW;
-- END
-- $$;
-- CREATE TRIGGER shop_product_update BEFORE INSERT OR UPDATE ON product FOR EACH ROW EXECUTE PROCEDURE update_name_norm_product();

-- Shop Product Name Normalize -- use for searching --
-- UPDATE shop_product SET name_norm = to_tsvector('simple', unaccent(shop_product.name));
CREATE INDEX shop_product_search_idx ON shop_product USING gin(name_norm);
-- CREATE OR REPLACE FUNCTION update_name_norm_shop_product() RETURNS trigger
--     LANGUAGE plpgsql
--     AS $$
-- BEGIN
-- 		if NEW.name != '' then
-- 			NEW.name_norm = to_tsvector('simple', unaccent(NEW.name));
-- 		END IF;
--     RETURN NEW;
-- END
-- $$;
-- CREATE TRIGGER shop_product_update BEFORE INSERT OR UPDATE ON shop_product FOR EACH ROW EXECUTE PROCEDURE update_name_norm_shop_product();

ALTER TABLE order_external
  RENAME COLUMN source_id TO order_source_id;

-- ALTER TABLE "order" DROP CONSTRAINT ids_not_empty;
-- ALTER TABLE "fulfillment" DROP CONSTRAINT ids_not_empty;

alter table order_line ADD COLUMN IF NOT EXISTS is_outside_etop BOOLEAN DEFAULT FALSE;
alter table "order" ADD COLUMN IF NOT EXISTS is_outside_etop BOOLEAN DEFAULT FALSE;

-- Order Code --
ALTER TABLE shop ADD COLUMN IF NOT EXISTS code TEXT UNIQUE;
ALTER TABLE order_line ADD COLUMN IF NOT EXISTS code TEXT UNIQUE;

CREATE TYPE code_type AS ENUM ('order', 'shop');

CREATE TABLE code (
  code text NOT NULL,
  type code_type,
  created_at timestamp with time zone,
  PRIMARY KEY (code, type)
);

ALTER TABLE shop add column if not EXISTS auto_create_ffm BOOLEAN DEFAULT FALSE;
ALTER TABLE fulfillment
  ADD COLUMN last_sync_at timestamptz;

CREATE TYPE ghn_note_code AS ENUM (
    'CHOTHUHANG',
    'CHOXEMHANGKHONGTHU',
    'KHONGCHOXEMHANG'
);

ALTER TABLE shop
    add column if not EXISTS ghn_note_code ghn_note_code;
ALTER TABLE "order" ADD COLUMN IF NOT EXISTS ghn_note_code ghn_note_code;
ALTER TABLE fulfillment ADD COLUMN IF NOT EXISTS expected_delivery_at TIMESTAMP WITH TIME ZONE;
