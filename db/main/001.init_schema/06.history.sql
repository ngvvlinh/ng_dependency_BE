SELECT init_history('account', '{id}');
SELECT init_history('account_user', '{account_id, user_id}');
SELECT init_history('address', '{id,account_id}', '{id}');
SELECT init_history('fulfillment', '{id,order_id,variant_ids,shop_id,supplier_id}', '{id}');
SELECT init_history('order_line', '{order_id,variant_id,shop_id,supplier_id}', '{order_id,variant_id}');
SELECT init_history('order', '{id,shop_id}', '{id}');
SELECT init_history('product_shop_collection', '{product_id,collection_id}');
SELECT init_history('shop_collection', '{id,shop_id}', '{id}');
SELECT init_history('shop_product', '{shop_id,product_id}');
SELECT init_history('shop_variant', '{shop_id,variant_id}');
SELECT init_history('shop', '{id}');
SELECT init_history('user_internal', '{id}');
SELECT init_history('user', '{id}');

CREATE INDEX ON history.fulfillment (order_id);
CREATE INDEX ON history.fulfillment (shop_id);
CREATE INDEX ON history.order_line (shop_id);
CREATE INDEX ON history.order (shop_id);
CREATE INDEX ON history.address (account_id);

SELECT init_history('credit', '{id,supplier_id,shop_id}', '{id}');
SELECT init_history('money_transaction_shipping', '{id,supplier_id,shop_id}', '{id}');
SELECT init_history('money_transaction_shipping_external', '{id}', '{id}');
CREATE INDEX ON history.credit (shop_id);
CREATE INDEX ON history.money_transaction_shipping (shop_id);

SELECT init_history('import_attempt', '{id,user_id,account_id}');
SELECT init_history('partner', '{id}');
SELECT init_history('account_auth', '{auth_key,account_id}');
SELECT init_history('partner_relation', '{partner_id,subject_id}');
SELECT init_history('webhook', '{id,account_id}', '{id}');

SELECT init_history('shop_trader', '{id,shop_id}');
SELECT init_history('shop_customer', '{id,shop_id}');
SELECT init_history('shop_supplier', '{id,shop_id}');
SELECT init_history('shop_trader_address', '{id,shop_id,trader_id}');
SELECT init_history('payment', '{id}');
SELECT init_history('transaction', '{id}');
SELECT init_history('shop_carrier', '{id,shop_id}');
SELECT init_history('receipt', '{id,shop_id}');
SELECT init_history('shop_ledger', '{id, shop_id}');
SELECT init_history('purchase_order', '{id,shop_id}');
SELECT init_history('invitation', '{id,account_id}');
SELECT init_history('connection', '{id}');
SELECT init_history('shop_connection', '{shop_id,connection_id}');
SELECT init_history('shop_brand', '{id, shop_id}');
SELECT init_history('shop_customer_group', '{id,shop_id}');
SELECT init_history('shop_customer_group_customer', '{customer_id, group_id}');

select init_history('inventory_voucher', '{id, shop_id}');
select init_history('purchase_refund', '{id, shop_id}');
select init_history('inventory_variant', '{variant_id, shop_id}');
select init_history('refund', '{id, shop_id}');
select init_history('shipnow_fulfillment', '{id, shop_id}');
select init_history('shop_stocktake', '{id, shop_id}');
select init_history('shop_variant_supplier', '{variant_id, supplier_id}');
select init_history('shop_product_collection', '{shop_id, product_id, collection_id}');
