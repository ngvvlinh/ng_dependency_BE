alter table purchase_order
    add column supplier_full_name_norm Text,
    add column supplier_phone_norm Text;

alter table history."purchase_order"
    add column supplier_full_name_norm Text,
    add column supplier_phone_norm Text;
