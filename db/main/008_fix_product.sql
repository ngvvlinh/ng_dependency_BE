alter table product add column ed_code text;
alter table history.product add column ed_code text;
alter table product drop column sku;
alter table history.product drop column sku;
CREATE UNIQUE INDEX ON product (ed_code, product_source_id);
CREATE UNIQUE INDEX ON product (code, product_source_id);

alter table variant add column ed_code text;
alter table history.variant add column ed_code text;
CREATE UNIQUE INDEX ON variant (ed_code, product_source_id);
CREATE UNIQUE INDEX ON variant (code, product_source_id);

