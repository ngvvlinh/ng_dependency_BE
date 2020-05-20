alter table shop_product drop column supplier_id;

alter type trader_type add VALUE 'supplier' after 'vendor';
