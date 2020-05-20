CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;

CREATE SCHEMA history;

CREATE TYPE public.account_type AS ENUM (
  'etop',
  'supplier',
  'shop',
  'partner',
  'affiliate'
  );

CREATE TYPE public.address_type AS ENUM (
  'billing',
  'shipping',
  'general',
  'warehouse',
  'shipfrom',
  'shipto'
  );

CREATE TYPE public.code_type AS ENUM (
  'order',
  'money_transaction',
  'shop',
  'money_transaction_external',
  'money_transaction_shipping',
  'money_transaction_shipping_etop',
  'connection'
  );

CREATE TYPE public.contact_type AS ENUM (
  'phone',
  'email'
  );

CREATE TYPE public.customer_type AS ENUM (
  'individual',
  'organization',
  'independent',
  'anonymous'
  );

CREATE TYPE public.fulfillment_endpoint AS ENUM (
  'supplier',
  'shop',
  'customer'
  );

CREATE TYPE public.gender_type AS ENUM (
  'male',
  'female',
  'other'
  );

CREATE TYPE public.ghn_note_code AS ENUM (
  'CHOTHUHANG',
  'CHOXEMHANGKHONGTHU',
  'KHONGCHOXEMHANG'
  );

CREATE TYPE public.import_type AS ENUM (
  'other',
  'shop_order',
  'shop_product',
  'etop_money_transaction'
  );

CREATE TYPE public.inventory_voucher_type AS ENUM (
  'in',
  'out'
  );

CREATE TYPE public.order_source_type AS ENUM (
  'unknown',
  'self',
  'import',
  'api',
  'etop_pos',
  'etop_pxs',
  'etop_cmx',
  'ts_app',
  'etop_app',
  'haravan'
  );

CREATE TYPE public.partial_status AS ENUM (
  'default',
  'partial',
  'done',
  'cancelled'
  );

CREATE TYPE public.payment_method_type AS ENUM (
  'other',
  'bank',
  'cod'
  );

CREATE TYPE public.processing_status AS ENUM (
  'default',
  'processing',
  'done',
  'cancelled'
  );

CREATE TYPE public.product_source_type AS ENUM (
  'custom',
  'kiotviet'
  );

CREATE TYPE public.province_region AS ENUM (
  'north',
  'south',
  'middle'
  );

CREATE TYPE public.receipt_created_type AS ENUM (
  'manual',
  'auto'
  );

CREATE TYPE public.receipt_ref_type AS ENUM (
  'order',
  'fulfillment',
  'inventory_voucher',
  'purchase_order',
  'refund',
  'purchase_refund'
  );

CREATE TYPE public.receipt_type AS ENUM (
  'receipt',
  'payment'
  );

CREATE TYPE public.shipping_provider AS ENUM (
  'ghn',
  'manual',
  'ghtk',
  'vtpost',
  'partner'
  );

CREATE TYPE public.shipping_state AS ENUM (
  'default',
  'unknown',
  'created',
  'confirmed',
  'picking',
  'processing',
  'holding',
  'returning',
  'returned',
  'delivering',
  'delivered',
  'undeliverable',
  'cancelled',
  'closed'
  );

CREATE TYPE public.shop_ledger_type AS ENUM (
  'cash',
  'bank'
  );

CREATE TYPE public.subject_type AS ENUM (
  'account',
  'user'
  );

CREATE TYPE public.tg_op_type AS ENUM (
  'INSERT',
  'UPDATE',
  'DELETE'
  );

CREATE TYPE public.trader_type AS ENUM (
  'customer',
  'vendor',
  'supplier',
  'carrier'
  );

CREATE TYPE public.try_on AS ENUM (
  'none',
  'open',
  'try'
  );

CREATE TYPE public.user_identifying_type AS ENUM (
  'full',
  'half',
  'stub'
  );
