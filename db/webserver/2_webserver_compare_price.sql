ALTER TABLE ws_product drop column compare_price ;

ALTER TABLE ws_product add column compare_price jsonb;

ALTER TABLE ws_website drop column banenrs ;

ALTER TABLE ws_website add column banners jsonb;

CREATE index ON "shop_product"("name");

CREATE index ON "shop_category"("name");
