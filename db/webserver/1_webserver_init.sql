create table ws_category (
  id bigint not null primary key,
  shop_id bigint,
  slug text,
  seo_config jsonb,
  image text,
  appear bool,
  created_at timestamptz,
  updated_at timestamptz
);

create table ws_product (
  id bigint not null primary key,
  shop_id bigint,
  slug text,
  seo_config jsonb,
  desc_html text,
  compare_price jsonb,
  appear bool,
  created_at timestamptz,
  updated_at timestamptz
);

create table ws_page (
  shop_id bigint ,
  id bigint not null primary key,
  seo_config jsonb,
  "name" text,
  slug text,
  desc_html text,
  image text,
  appear bool,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz
);

create table ws_website (
  id bigint not null primary key,
  shop_id bigint,
  main_color text,
  banenrs jsonb,
  outstanding_product jsonb,
  new_product jsonb,
  seo_config jsonb,
  facebook jsonb,
  google_analytics_id text,
  domain_name text,
  over_stock int,
  shop_info jsonb,
  description text,
  logo_image text,
  favicon_image text,
  created_at timestamptz,
  updated_at timestamptz
);
