ALTER TABLE ws_product drop column compare_price ;

ALTER TABLE ws_product add column compare_price jsonb;

CREATE index ON "shop_product"("name");

CREATE index ON "shop_category"("name");
