ALTER TABLE shop_customer
    ADD COLUMN created_by INT8 references "user"(id);
ALTER TABLE "history".shop_customer
    ADD COLUMN created_by INT8;
