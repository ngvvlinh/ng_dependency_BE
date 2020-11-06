CREATE TABLE contact (
    id int8,
    shop_id int8,
    full_name text,
    phone text,
    wl_partner_id int8,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

SELECT init_history('contact', '{id, shop_id}');
