CREATE TYPE account_type AS ENUM (
    'etop',
    'supplier',
    'shop'
);

CREATE TYPE address_type AS ENUM (
    'billing',
    'shipping',
    'general',
    'warehouse'
);

create type product_source_type as enum (
    'custom',
    'kiotviet'
);

CREATE TYPE order_source_type AS ENUM (
    'unknown',
    'self',
    'etop_pos',
    'etop_pxs',
    'haravan'
);

CREATE TYPE partial_status AS ENUM (
    'default',
    'partial',
    'done',
    'cancelled'
);

CREATE TYPE processing_status AS ENUM (
    'default',
    'processing',
    'done',
    'cancelled'
);

CREATE TYPE shipping_state AS ENUM (
    'default',       -- have not created shipping order on GHN yet
    'unknown',
    'created',       -- created shipping order on GHN
    'confirmed',     -- unused
    'picking',
    'processing',    -- unused
    'holding',
    'returning',
    'returned',
    'delivering',
    'delivered',
    'undeliverable',
    'cancelled',
    'closed'
);

CREATE TABLE account (
    id bigint PRIMARY KEY,
    name text NOT NULL,
    type account_type NOT NULL,
    deleted_at timestamp with time zone,
    image_url text
);

CREATE TABLE "user" (
    id bigint PRIMARY KEY,
    rid bigint,
    status int2 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    full_name text NOT NULL,
    short_name text NOT NULL,
    email text NOT NULL,
    phone text NOT NULL
);

CREATE TABLE product_source (
    id int8 PRIMARY KEY,
    rid int8,
    type product_source_type,
    name text,
    supplier_id int8 UNIQUE,
    status int2,
    external_status int2,
    external_key text,
    external_info jsonb,
    extra_info jsonb,
    created_at timestamptz,
    updated_at timestamptz,
    last_sync_at timestamp with time zone,
    sync_state_products jsonb,
    sync_state_categories jsonb
);

CREATE UNIQUE INDEX ON product_source (type, external_key) WHERE type in ('kiotviet');

CREATE TABLE address (
    id bigint primary key,
    country text DEFAULT 'VN'::text,
    province_code text,
    province text,
    district_code text,
    district text,
    ward text,
    ward_code text,
    address text,
    is_default boolean DEFAULT false,
    type address_type,
    account_id bigint REFERENCES account(id),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE order_source (
    id bigint PRIMARY key,
    rid bigint,
    type order_source_type,
    name text not null,
    shop_id int8 not null unique,
    status int2 not null,
    external_status int2,
    external_key text,
    external_info jsonb,
    extra_info jsonb,
    created_at timestamptz,
    updated_at timestamptz,
    last_sync_at timestamptz,
    sync_state_orders jsonb
);

CREATE TABLE shop (
    id bigint PRIMARY KEY REFERENCES account(id),
    rid bigint,
    name text NOT NULL,
    owner_id bigint NOT NULL REFERENCES "user"(id),
    status int2 NOT NULL,
    product_source_id int8 REFERENCES product_source(id),
    order_source_id int8 REFERENCES order_source(id),
    created_at timestamp with time zone DEFAULT date_trunc('second', now()),
    updated_at timestamp with time zone,
    rules jsonb,
    is_test smallint DEFAULT '0'::smallint,
    image_url text,
    phone text,
    website_url text,
    email text,
    deleted_at timestamp with time zone,
    address_id bigint REFERENCES address(id),
    bank_account jsonb,
    contact_persons jsonb
);

CREATE TABLE supplier (
    id bigint PRIMARY KEY REFERENCES account(id),
    rid bigint,
    status int2 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    name text NOT NULL,
    owner_id bigint NOT NULL REFERENCES "user"(id),
    product_source_id int8 NOT NULL UNIQUE REFERENCES product_source(id) DEFERRABLE INITIALLY DEFERRED,
    rules jsonb,
    is_test smallint DEFAULT '0'::smallint,
    company_info jsonb,
    warehouse_address_id bigint REFERENCES address(id),
    bank_account jsonb,
    contact_persons jsonb,
    image_url text,
    deleted_at timestamp with time zone
);

CREATE TABLE user_internal (
    id bigint PRIMARY KEY REFERENCES "user"(id),
    hashpwd text,
    updated_at timestamp with time zone
);

CREATE TABLE user_auth (
    user_id bigint REFERENCES "user"(id),
    auth_type text NOT NULL,
    auth_key text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (auth_type, auth_key)
);

CREATE TABLE account_user (
    account_id bigint NOT NULL REFERENCES account(id),
    user_id bigint NOT NULL REFERENCES "user"(id),
    status int2 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    roles text[],
    permissions text[],
    deleted_at text,
    PRIMARY KEY (account_id, user_id)
);

CREATE TABLE etop_category (
    id bigint primary key,
    name text NOT NULL,
    status int2 not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    parent_id bigint REFERENCES etop_category(id)
);

CREATE UNIQUE INDEX ON order_source(type, external_key) WHERE TYPE IN ('haravan');
CREATE UNIQUE INDEX ON order_source(shop_id) WHERE TYPE IN ('haravan');

CREATE TABLE order_source_internal (
    id int8 PRIMARY KEY,
    rid int8,
    secret jsonb,
    access_token text,
    expires_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE product_brand (
    id bigint PRIMARY KEY,
    name text,
    description text,
    policy text,
    supplier_id bigint NOT NULL REFERENCES supplier(id),
    image_urls text[],
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE product_source_category (
    id bigint PRIMARY KEY,
    rid bigint,
    product_source_id bigint REFERENCES product_source(id),
    product_source_type product_source_type,
    supplier_id int8,
    parent_id int8,
    name text,
    status int2 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE product_source_category_external (
    id bigint PRIMARY KEY REFERENCES product_source_category(id),
    rid bigint,
    product_source_id bigint REFERENCES product_source(id),
    product_source_type product_source_type,
    external_id text,
    external_parent_id text,
    external_code text,
    external_name text,
    external_status int2 NOT NULL,
    external_updated_at timestamp with time zone,
    external_created_at timestamp with time zone,
    external_deleted_at timestamp with time zone,
    last_sync_at timestamp with time zone
);

CREATE TABLE product (
    id bigint primary key,
    rid bigint,
    product_source_id bigint REFERENCES product_source(id),
    -- external_id text,
    supplier_id bigint REFERENCES supplier(id),
    product_source_category_id bigint REFERENCES product_source_category(id),
    etop_category_id bigint REFERENCES etop_category(id),
    product_brand_id bigint REFERENCES product_brand(id),

    name text,
    short_desc text,
    description text,
    desc_html text,

    ed_name text,
    ed_short_desc text,
    ed_description text,
    ed_desc_html text,

    status int2 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamptz,

    sku text,
    code text,

    quantity_available integer,
    quantity_on_hand integer,
    quantity_reserved integer,

    image_urls text[]
);

CREATE TABLE product_external (
    id bigint PRIMARY key REFERENCES product(id),
    rid bigint,
    product_source_id bigint REFERENCES product_source(id),
    product_source_type product_source_type,

    external_id text not NULL,
    external_name text not null,
    external_code text,
    external_category_id text,
    external_description text,
    -- external_attributes text,
    external_image_urls text[],
    external_unit text,
    external_data jsonb,
    external_status int2,
    external_created_at timestamptz,
    external_updated_at timestamptz,
    external_deleted_at timestamptz,
    last_sync_at timestamptz
);

CREATE UNIQUE INDEX ON product_external(product_source_type, external_id) WHERE product_source_type in ('kiotviet');

CREATE TABLE variant (
    id bigint PRIMARY KEY,
    rid bigint,
    product_id bigint REFERENCES product(id),
    product_source_id bigint REFERENCES product_source(id),
    supplier_id bigint REFERENCES supplier(id),

    name text,
    short_desc text,
    description text,
    desc_html text,

    ed_name text,
    ed_short_desc text,
    ed_description text,
    ed_desc_html text,

    desc_norm tsvector,
    name_norm tsvector,

    status int2 not null,
    etop_status int2 not null,

    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    sku text,
    code text,

    wholesale_price integer,
    wholesale_price_0 integer,
    list_price integer,
    retail_price_min integer,
    retail_price_max integer,

    ed_wholesale_price integer,
    ed_wholesale_price_0 integer,
    ed_list_price integer,
    ed_retail_price_min integer,
    ed_retail_price_max integer,

    quantity_available integer,
    quantity_on_hand integer,
    quantity_reserved integer,

    image_urls text[],
    supplier_meta jsonb,
    product_source_category_id bigint REFERENCES product_source_category(id),
    etop_category_id bigint REFERENCES etop_category(id),
    product_brand_id bigint REFERENCES product_brand(id)
);

CREATE TABLE variant_external (
	id int8 primary key references variant(id),
    rid bigint,
    product_source_id bigint NOT NULL REFERENCES product_source(id),
    product_source_type product_source_type NOT NULL,

    external_id text NOT NULL,
    external_product_id text,
    external_category_id text,
    external_code text,
    external_name text,
    -- external_full_name text,
    external_description text,
    external_attributes jsonb,
    external_image_urls text[],
    external_unit text,
    external_conversion_value real,
    external_price integer,
    external_unit_conv real,
    external_base_unit_id text,
    external_data jsonb NOT NULL,
    external_name_norm tsvector,
    external_status int2,
    external_created_at timestamp with time zone,
    external_updated_at timestamp with time zone,
    external_deleted_at timestamp with time zone,
    last_sync_at timestamp with time zone
);

CREATE UNIQUE INDEX ON variant_external(product_source_type, external_id);

CREATE TABLE shop_collection (
    id bigint PRIMARY KEY,
    rid bigint,
    shop_id bigint NOT NULL REFERENCES shop(id),
    name text,
    description text,
    desc_html text,
    short_desc text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE shop_variant (
    shop_id bigint NOT NULL REFERENCES shop(id),
    variant_id bigint NOT NULL REFERENCES variant(id),
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
    PRIMARY KEY (shop_id, variant_id)
);

CREATE TABLE product_source_internal (
    id int8 PRIMARY key REFERENCES product_source(id),
    rid int8,
    secret jsonb,
    access_token text,
    expires_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

CREATE TABLE "order" (
    id bigint primary key,
    rid bigint,
    shop_id bigint NOT NULL REFERENCES shop(id),
    shop_name text,
    code text NOT NULL,
--     product_ids bigint[],
    variant_ids bigint[],
    supplier_ids bigint[],
    currency text,
    payment_method text,
    customer jsonb,
    customer_address jsonb,
    billing_address jsonb,
    shipping_address jsonb,
    shop_address jsonb,
    customer_phone text,
    customer_email text,
    created_at timestamp with time zone,
    processed_at timestamp with time zone,
    updated_at timestamp with time zone,
    closed_at timestamp with time zone,
    confirmed_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    cancel_reason text,
    customer_confirm smallint,
    external_confirm smallint,
    shop_confirm smallint,
    confirm_status smallint,
    fulfillment_status smallint,
    status smallint,
    lines jsonb,
    discounts jsonb,
    total_items integer,
    basket_value integer,
    total_weight integer,
    total_tax integer,
    total_discount integer,
    total_amount integer,

    order_note text,
    shop_note text,
    shipping_note text,
    order_source_id int8,
    order_source_type order_source_type,
    external_order_id text
);

create table order_external (
    id bigint primary key REFERENCES "order"(id),
    rid bigint,
    source_id bigint REFERENCES order_source(id),
    external_order_source text,
    external_provider text,
--     external_shop_id text,
    external_order_id text,
    external_order_code text,
    external_user_id text,
    external_customer_id text,
    external_created_at timestamp with time zone,
    external_updated_at timestamp with time zone,
    external_processed_at timestamp with time zone,
    external_closed_at timestamp with time zone,
    external_cancelled_at timestamp with time zone,
    external_cancel_reason text,
    external_data jsonb,
    external_lines jsonb
);

CREATE TABLE order_line (
    rid bigint,
    order_id bigint NOT NULL REFERENCES "order"(id),
    product_id bigint not null references product(id),
    variant_id bigint NOT NULL REFERENCES variant(id),
    supplier_id bigint REFERENCES supplier(id),
    external_variant_id text,
    external_supplier_order_id text,
    product_name text,
    supplier_name text,
    image_url text,
    product_exists boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    closed_at timestamp with time zone,
    confirmed_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    cancel_reason text,
    supplier_confirm smallint,
    status smallint,
    weight integer,
    quantity integer NOT NULL,
    wholesale_price_0 integer NOT NULL,
    wholesale_price integer NOT NULL,
    list_price integer NOT NULL,
    retail_price integer NOT NULL,
    payment_price integer NOT NULL,
    line_amount integer NOT NULL,
    total_discount integer NOT NULL,
    total_line_amount integer NOT NULL,
    requires_shipping boolean,
    primary key (order_id, variant_id)
);

CREATE FUNCTION account_shop_supplier_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE account SET name = NEW.name, image_url = NEW.image_url, deleted_at = NEW.deleted_at WHERE id = NEW.id;
    RETURN NEW;
END;
$$;

CREATE FUNCTION product_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    supplier_status INT2;
    external_status INT2;
    etop_status     INT2;
    price_status    INT2;
BEGIN
    IF TG_OP = 'UPDATE' THEN
        supplier_status := COALESCE(NEW.supplier_status, OLD.supplier_status, 0);
        external_status := COALESCE(NEW.external_status, OLD.external_status, 0);
        etop_status     := COALESCE(NEW.etop_status,     OLD.etop_status,     0);
    ELSE
        supplier_status := COALESCE(NEW.supplier_status, 0);
        external_status := COALESCE(NEW.external_status, 0);
        etop_status     := COALESCE(NEW.etop_status,     0);
    END IF;

    IF NEW.wholesale_price > 0 AND NEW.list_price > 0 THEN
        price_status := 1;
    ELSE
        price_status := 0;
    END IF;
    NEW.status := LEAST(supplier_status, external_status, etop_status, price_status);
    RETURN NEW;
END
$$;

CREATE FUNCTION supplier_rules_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (NEW.rules != OLD.rules) THEN
        PERFORM pg_notify('supplier_rules_update', NEW.id::TEXT);
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON shop FOR EACH ROW EXECUTE PROCEDURE account_shop_supplier_update();

CREATE TRIGGER account_update AFTER INSERT OR UPDATE ON supplier FOR EACH ROW EXECUTE PROCEDURE account_shop_supplier_update();

CREATE TRIGGER product_update BEFORE INSERT OR UPDATE ON product FOR EACH ROW EXECUTE PROCEDURE product_update();

CREATE TRIGGER supplier_rules_update AFTER UPDATE ON supplier FOR EACH ROW EXECUTE PROCEDURE supplier_rules_update();

ALTER TABLE product_source ADD CONSTRAINT product_source_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES supplier(id);
ALTER TABLE order_source ADD CONSTRAINT order_source_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES shop(id);
