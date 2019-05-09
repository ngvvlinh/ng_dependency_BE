ALTER TABLE "order" DROP COLUMN shop_name;
ALTER TABLE order_line DROP COLUMN product_exists;

ALTER TABLE history."order" DROP COLUMN shop_name;
ALTER TABLE history.order_line DROP COLUMN product_exists;

ALTER TABLE "order" ADD COLUMN ed_code TEXT;
ALTER TABLE history."order" ADD COLUMN ed_code TEXT;

ALTER TABLE "order" ALTER COLUMN order_source_type SET NOT NULL;

ALTER TYPE order_source_type ADD VALUE 'import' AFTER 'self';

CREATE TYPE payment_method_type AS ENUM ('other', 'bank', 'cod');
