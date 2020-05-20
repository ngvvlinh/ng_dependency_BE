
alter table shop_product rename column vendor_id to supplier_id;

alter table history.shop_product
    rename column vendor_id to supplier_id;

alter table shop_vendor rename to shop_supplier;
