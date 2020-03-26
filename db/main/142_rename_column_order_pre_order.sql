ALTER TABLE "order"
RENAME COLUMN "preorder" TO "pre_order";

ALTER TABLE history."order"
RENAME COLUMN "preorder" TO "pre_order";
