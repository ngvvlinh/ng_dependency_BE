alter table history.credit
  alter column amount type int8
;

alter table history.fulfillment
  alter column total_items type int8
, alter column total_weight type int8
, alter column basket_value type int8
, alter column total_cod_amount type int8
, alter column shipping_fee_customer type int8
, alter column shipping_fee_shop type int8
, alter column external_shipping_fee type int8
, alter column etop_discount type int8
, alter column etop_fee_adjustment type int8
, alter column shipping_fee_main type int8
, alter column shipping_fee_return type int8
, alter column shipping_fee_insurance type int8
, alter column shipping_fee_adjustment type int8
, alter column shipping_fee_cods type int8
, alter column shipping_fee_info_change type int8
, alter column shipping_fee_other type int8
, alter column total_discount type int8
, alter column total_amount type int8
, alter column shipping_service_fee type int8
, alter column original_cod_amount type int8
, alter column shipping_fee_discount type int8
, alter column etop_adjusted_shipping_fee_main type int8
, alter column actual_compensation_amount type int8
, alter column gross_weight type int8
, alter column chargeable_weight type int8
, alter column length type int8
, alter column width type int8
, alter column height type int8
;

alter table history.import_attempt
  alter column duration_ms type int8
;

alter table history.inventory_variant
  alter column quantity_on_hand type int8
, alter column quantity_picked type int8
, alter column cost_price type int8
;

alter table history.inventory_voucher
  alter column total_amount type int8
, alter column code_norm type int8
;

alter table history.money_transaction_shipping
  alter column total_cod type int8
, alter column total_orders type int8
, alter column total_amount type int8
;

alter table history."order"
  alter column total_items type int8
, alter column basket_value type int8
, alter column total_weight type int8
, alter column total_tax type int8
, alter column total_discount type int8
, alter column total_amount type int8
, alter column shop_shipping_fee type int8
, alter column shop_cod type int8
, alter column order_discount type int8
, alter column total_fee type int8
;

alter table history.order_line
  alter column weight type int8
, alter column quantity type int8
, alter column wholesale_price_0 type int8
, alter column wholesale_price type int8
, alter column list_price type int8
, alter column retail_price type int8
, alter column payment_price type int8
, alter column line_amount type int8
, alter column total_discount type int8
, alter column total_line_amount type int8
;

alter table history.payment
  alter column amount type int8
;

alter table history.purchase_order
  alter column code_norm type int8
, alter column total_fee type int8
;

alter table history.purchase_refund
  alter column code_norm type int8
, alter column total_amount type int8
, alter column basket_value type int8
, alter column total_adjustment type int8
;

alter table history.receipt
  alter column amount type int8
, alter column code_norm type int8
;

alter table history.refund
  alter column code_norm type int8
, alter column total_amount type int8
, alter column basket_value type int8
, alter column total_adjustment type int8
;

alter table history.shipnow_fulfillment
  alter column shipping_service_fee type int8
, alter column chargeable_weight type int8
, alter column gross_weight type int8
, alter column basket_value type int8
, alter column cod_amount type int8
, alter column total_fee type int8
;

alter table history.shop_customer
  alter column code_norm type int8
;

alter table history.shop_product
  alter column retail_price type int8
, alter column cost_price type int8
, alter column list_price type int8
, alter column code_norm type int8
;

alter table history.shop_stocktake
  alter column total_quantity type int8
, alter column code_norm type int8
;

alter table history.shop_supplier
  alter column code_norm type int8
;

alter table history.shop_variant
  alter column retail_price type int8
, alter column cost_price type int8
, alter column list_price type int8
, alter column code_norm type int8
;

alter table history.transaction
  alter column amount type int8
;

alter table public.credit
  alter column amount type int8
;

alter table public.district
  alter column ghn_id type int8
, alter column vtpost_id type int8
;

alter table public.export_attempt
  alter column n_total type int8
, alter column n_exported type int8
, alter column n_error type int8
;

alter table public.fulfillment
  alter column total_items type int8
, alter column total_weight type int8
, alter column basket_value type int8
, alter column total_cod_amount type int8
, alter column shipping_fee_customer type int8
, alter column shipping_fee_shop type int8
, alter column external_shipping_fee type int8
, alter column etop_discount type int8
, alter column etop_fee_adjustment type int8
, alter column shipping_fee_main type int8
, alter column shipping_fee_return type int8
, alter column shipping_fee_insurance type int8
, alter column shipping_fee_adjustment type int8
, alter column shipping_fee_cods type int8
, alter column shipping_fee_info_change type int8
, alter column shipping_fee_other type int8
, alter column total_discount type int8
, alter column total_amount type int8
, alter column shipping_service_fee type int8
, alter column original_cod_amount type int8
, alter column shipping_fee_discount type int8
, alter column etop_adjusted_shipping_fee_main type int8
, alter column actual_compensation_amount type int8
, alter column gross_weight type int8
, alter column chargeable_weight type int8
, alter column length type int8
, alter column width type int8
, alter column height type int8
;

alter table public.import_attempt
  alter column duration_ms type int8
;

alter table public.inventory_variant
  alter column quantity_on_hand type int8
, alter column quantity_picked type int8
, alter column cost_price type int8
;

alter table public.inventory_voucher
  alter column total_amount type int8
, alter column code_norm type int8
;

alter table public.money_transaction_shipping
  alter column total_cod type int8
, alter column total_orders type int8
, alter column total_amount type int8
;

alter table public.money_transaction_shipping_etop
  alter column total_cod type int8
, alter column total_orders type int8
, alter column total_amount type int8
, alter column total_fee type int8
, alter column total_money_transaction type int8
;

alter table public.money_transaction_shipping_external
  alter column total_cod type int8
, alter column total_orders type int8
;

alter table public.money_transaction_shipping_external_line
  alter column external_total_cod type int8
, alter column external_total_shipping_fee type int8
;

alter table public."order"
  alter column total_items type int8
, alter column basket_value type int8
, alter column total_weight type int8
, alter column total_tax type int8
, alter column total_discount type int8
, alter column total_amount type int8
, alter column shop_shipping_fee type int8
, alter column shop_cod type int8
, alter column order_discount type int8
, alter column total_fee type int8
;

alter table public.order_line
  alter column weight type int8
, alter column quantity type int8
, alter column wholesale_price_0 type int8
, alter column wholesale_price type int8
, alter column list_price type int8
, alter column retail_price type int8
, alter column payment_price type int8
, alter column line_amount type int8
, alter column total_discount type int8
, alter column total_line_amount type int8
;

alter table public.payment
  alter column amount type int8
;

alter table public.province
  alter column vtpost_id type int8
;

alter table public.purchase_order
  alter column code_norm type int8
, alter column total_fee type int8
;

alter table public.purchase_refund
  alter column code_norm type int8
, alter column total_amount type int8
, alter column basket_value type int8
, alter column total_adjustment type int8
;

alter table public.receipt
  alter column amount type int8
, alter column code_norm type int8
;

alter table public.refund
  alter column code_norm type int8
, alter column total_amount type int8
, alter column basket_value type int8
, alter column total_adjustment type int8
;

alter table public.shipnow_fulfillment
  alter column shipping_service_fee type int8
, alter column chargeable_weight type int8
, alter column gross_weight type int8
, alter column basket_value type int8
, alter column cod_amount type int8
, alter column total_fee type int8
;

alter table public.shop_customer
  alter column code_norm type int8
;

alter table public.shop_product
  alter column retail_price type int8
, alter column cost_price type int8
, alter column list_price type int8
, alter column code_norm type int8
;

alter table public.shop_stocktake
  alter column total_quantity type int8
, alter column code_norm type int8
;

alter table public.shop_supplier
  alter column code_norm type int8
;

alter table public.shop_variant
  alter column retail_price type int8
, alter column cost_price type int8
, alter column list_price type int8
, alter column code_norm type int8
;

alter table public.transaction
  alter column amount type int8
;
