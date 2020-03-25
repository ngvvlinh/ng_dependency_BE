CREATE TABLE if not exists purchase_order (
    id INT8 PRIMARY KEY,
    shop_id INT8,
    supplier_id INT8,
    supplier JSONB,
    basket_value INT8,
    total_discount INT8,
    total_amount INT8,
    paid_amount INT8,
    code TEXT,
    code_norm INT4,
    note TEXT,
    status INT2,
    variant_ids INT8[],
    lines JSONB,
    cancelled_reason TEXT,
    created_by INT8,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    total_fee int,
    rid int8
);

alter table purchase_order
    drop column if exists discount_lines,
    drop column if exists fee_lines;
