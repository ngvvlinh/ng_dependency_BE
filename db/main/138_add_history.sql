select init_history('inventory_voucher', '{id, shop_id}');

select init_history('purchase_refund', '{id, shop_id}');

select init_history('inventory_variant', '{variant_id, shop_id}');

select init_history('refund', '{id, shop_id}');

select init_history('shipnow_fulfillment', '{id, shop_id}');

select init_history('shop_stocktake', '{id, shop_id}');

select init_history('shop_variant_supplier', '{variant_id, supplier_id}');