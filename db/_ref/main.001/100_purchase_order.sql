CREATE TABLE purchase_order (
    id INT8 PRIMARY KEY,
    shop_id INT8 NOT NULL REFERENCES shop(id),
    supplier_id INT8 NOT NULL REFERENCES shop_supplier(id),
    supplier JSONB,
    basket_value INT8,
    total_discount INT8,
    total_amount INT8,
    paid_amount INT8,
    code TEXT NOT NULL,
    code_norm INT4,
    note TEXT,
    status INT2,
    variant_ids INT8[],
    lines JSONB,
    cancelled_reason TEXT,
    created_by INT8 NOT NULL REFERENCES "user"(id),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX ON purchase_order USING GIN(variant_ids);
CREATE UNIQUE INDEX ON purchase_order (shop_id, code);

select init_history('purchase_order', '{id,shop_id}');
