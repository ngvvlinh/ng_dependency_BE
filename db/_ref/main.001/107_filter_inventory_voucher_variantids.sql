alter table inventory_voucher add column variant_ids int8[];
CREATE INDEX ON inventory_voucher USING gin(variant_ids);
