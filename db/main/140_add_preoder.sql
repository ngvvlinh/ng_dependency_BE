alter table "order" add column "preorder" bool;
alter table history."order" add column "preorder" bool;

ALTER TYPE order_source_type ADD VALUE 'ecomify' AFTER 'self';
