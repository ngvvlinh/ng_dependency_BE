CREATE INDEX account_auth_account_id_idx ON public.account_auth USING btree (account_id);

CREATE UNIQUE INDEX account_user_account_id_user_id_idx ON public.account_user USING btree (account_id, user_id) WHERE (deleted_at IS NULL);

CREATE UNIQUE INDEX connection_code_idx ON public.connection USING btree (code);

CREATE UNIQUE INDEX connection_name_idx ON public.connection USING btree (name);

CREATE INDEX credit_sum_idx ON public.credit USING btree (shop_id, status, paid_at, amount);

CREATE INDEX export_attempt_account_id_idx ON public.export_attempt USING btree (account_id);

CREATE UNIQUE INDEX ffm_active_supplier_key ON public.fulfillment USING btree (order_id, public.ffm_active_supplier(supplier_id, status));

CREATE INDEX fulfillment_address_to_district_code_idx ON public.fulfillment USING btree (address_to_district_code);

CREATE INDEX fulfillment_address_to_province_code_idx ON public.fulfillment USING btree (address_to_province_code);

CREATE INDEX fulfillment_address_to_ward_code_idx ON public.fulfillment USING btree (address_to_ward_code);

CREATE INDEX fulfillment_created_at_idx ON public.fulfillment USING btree (created_at);

CREATE INDEX fulfillment_external_shipping_code_idx ON public.fulfillment USING btree (external_shipping_code);

CREATE INDEX fulfillment_fulfillment_expected_pick_at_idx ON public.fulfillment USING btree (public.fulfillment_expected_pick_at(created_at));

CREATE INDEX fulfillment_money_transaction_id_idx ON public.fulfillment USING btree (money_transaction_id);

CREATE INDEX fulfillment_shipping_code_idx ON public.fulfillment USING btree (shipping_code);

CREATE INDEX fulfillment_shipping_fee_shop_idx ON public.fulfillment USING btree (shipping_fee_shop);

CREATE UNIQUE INDEX fulfillment_shipping_provider_external_shipping_code_idx ON public.fulfillment USING btree (shipping_provider, external_shipping_code) WHERE (status <> ALL (ARRAY['-1'::integer, 1]));

CREATE UNIQUE INDEX fulfillment_shipping_provider_shipping_code_idx ON public.fulfillment USING btree (shipping_provider, shipping_code) WHERE (status <> ALL (ARRAY['-1'::integer, 1]));

CREATE INDEX fulfillment_shipping_state_idx ON public.fulfillment USING btree (shipping_state);

CREATE INDEX fulfillment_status_sum_idx ON public.fulfillment USING btree (shop_id, status, shipping_status, etop_payment_status, shipping_fee_shop, total_cod_amount);

CREATE INDEX fulfillment_total_cod_amount_idx ON public.fulfillment USING btree (total_cod_amount);

CREATE INDEX fulfillment_updated_at_idx ON public.fulfillment USING btree (updated_at);

CREATE UNIQUE INDEX inventory_variant_shop_id_variant_id_idx ON public.inventory_variant USING btree (shop_id, variant_id);

CREATE INDEX inventory_voucher_created_at_idx ON public.inventory_voucher USING btree (created_at);

CREATE UNIQUE INDEX inventory_voucher_id_idx ON public.inventory_voucher USING btree (id);

CREATE INDEX inventory_voucher_ref_code_idx ON public.inventory_voucher USING btree (ref_code);

CREATE UNIQUE INDEX inventory_voucher_shop_id_code_idx ON public.inventory_voucher USING btree (shop_id, code);

CREATE INDEX inventory_voucher_status_idx ON public.inventory_voucher USING btree (status);

CREATE INDEX inventory_voucher_variant_ids_idx ON public.inventory_voucher USING gin (variant_ids);

CREATE INDEX invitation_token_idx ON public.invitation USING btree (token);

CREATE INDEX money_transaction_shipping_created_at_status_idx ON public.money_transaction_shipping USING btree (created_at, status);

CREATE INDEX money_transaction_shipping_shop_id_status_idx ON public.money_transaction_shipping USING btree (shop_id, status);

CREATE UNIQUE INDEX order_code_idx ON public."order" USING btree (code);

CREATE INDEX order_confirm_status_idx ON public."order" USING btree (confirm_status);

CREATE INDEX order_created_at_idx ON public."order" USING btree (created_at);

CREATE INDEX order_customer_id_idx ON public."order" USING btree (customer_id);

CREATE INDEX order_customer_name_norm_idx ON public."order" USING gin (customer_name_norm);

CREATE INDEX order_customer_phone_idx ON public."order" USING btree (customer_phone);

CREATE INDEX order_etop_payment_status_idx ON public."order" USING btree (etop_payment_status);

CREATE INDEX order_external_id_idx ON public."order" USING btree (external_order_id) WHERE (external_order_id IS NOT NULL);

CREATE INDEX order_fulfillment_shipping_codes_idx ON public."order" USING gin (fulfillment_shipping_codes);

CREATE INDEX order_fulfillment_shipping_states_idx ON public."order" USING gin (fulfillment_shipping_states);

CREATE INDEX order_fulfillment_shipping_status_idx ON public."order" USING btree (fulfillment_shipping_status);

CREATE INDEX order_line_order_id_idx ON public.order_line USING btree (order_id);

CREATE INDEX order_order_source_type_created_by_idx ON public."order" USING btree (order_source_type, created_by);

CREATE UNIQUE INDEX order_partner_external_id_idx ON public."order" USING btree (partner_id, external_order_id) WHERE ((external_order_id IS NOT NULL) AND (partner_id IS NOT NULL) AND (status <> '-1'::integer));

CREATE UNIQUE INDEX order_partner_shop_id_external_code_idx ON public."order" USING btree (shop_id, ed_code, partner_id) WHERE ((partner_id IS NOT NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer));

CREATE INDEX order_product_name_norm_idx ON public."order" USING gin (product_name_norm);

CREATE UNIQUE INDEX order_shop_external_id_idx ON public."order" USING btree (shop_id, external_order_id) WHERE ((external_order_id IS NOT NULL) AND (partner_id IS NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer));

CREATE UNIQUE INDEX order_shop_id_ed_code_idx ON public."order" USING btree (shop_id, ed_code) WHERE ((partner_id IS NULL) AND (status <> '-1'::integer) AND (fulfillment_shipping_status <> '-2'::integer));

CREATE INDEX order_shop_id_idx ON public."order" USING btree (shop_id);

CREATE INDEX order_status_idx ON public."order" USING btree (status);

CREATE INDEX order_total_amount_idx ON public."order" USING btree (total_amount);

CREATE UNIQUE INDEX partner_relation_auth_key_idx ON public.partner_relation USING btree (auth_key);

CREATE INDEX partner_relation_partner_id_idx ON public.partner_relation USING btree (partner_id);

CREATE UNIQUE INDEX partner_relation_partner_id_subject_id_idx ON public.partner_relation USING btree (partner_id, subject_id);

CREATE INDEX partner_relation_subject_id_idx ON public.partner_relation USING btree (subject_id);

CREATE UNIQUE INDEX payment_external_trans_id_payment_provider_idx ON public.payment USING btree (external_trans_id, payment_provider);

CREATE UNIQUE INDEX purchase_order_shop_id_code_idx ON public.purchase_order USING btree (shop_id, code);

CREATE INDEX purchase_order_variant_ids_idx ON public.purchase_order USING gin (variant_ids);

CREATE INDEX purchase_refund_shop_id_code_idx ON public.purchase_refund USING btree (shop_id, code);

CREATE INDEX purchase_refund_shop_id_id_idx ON public.purchase_refund USING btree (shop_id, id);

CREATE INDEX purchase_refund_shop_id_purchase_order_id_idx ON public.purchase_refund USING btree (shop_id, purchase_order_id);

CREATE INDEX receipt_created_at_idx ON public.receipt USING btree (created_at);

CREATE INDEX receipt_ledger_id_idx ON public.receipt USING btree (ledger_id);

CREATE INDEX receipt_order_ids_idx ON public.receipt USING gin (ref_ids);

CREATE INDEX receipt_paid_at_idx ON public.receipt USING btree (paid_at);

CREATE UNIQUE INDEX receipt_shop_id_code_idx ON public.receipt USING btree (shop_id, code);

CREATE INDEX receipt_trader_full_name_norm_idx ON public.receipt USING gin (trader_full_name_norm);

CREATE INDEX receipt_trader_full_name_norm_idx1 ON public.receipt USING gin (trader_full_name_norm);

CREATE INDEX receipt_trader_phone_norm_idx ON public.receipt USING gin (trader_phone_norm);

CREATE INDEX receipt_trader_type_idx ON public.receipt USING btree (trader_type);

CREATE INDEX receipt_type_idx ON public.receipt USING btree (type);

CREATE INDEX refund_shop_id_code_idx ON public.refund USING btree (shop_id, code);

CREATE INDEX refund_shop_id_id_idx ON public.refund USING btree (shop_id, id);

CREATE INDEX refund_shop_id_order_id_idx ON public.refund USING btree (shop_id, order_id);

CREATE INDEX shipnow_fulfillment_order_ids_idx ON public.shipnow_fulfillment USING gin (order_ids);

CREATE UNIQUE INDEX shipping_source_name_type_username_idx ON public.shipping_source USING btree (name, type, username);

CREATE UNIQUE INDEX shop_brand_partner_id_external_id_idx ON public.shop_brand USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_brand_shop_id_id_idx ON public.shop_brand USING btree (shop_id, id);

CREATE UNIQUE INDEX shop_category_partner_id_external_id_idx ON public.shop_category USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_collection_partner_id_external_id_idx ON public.shop_collection USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_connection_connection_id_idx ON public.shop_connection USING btree (connection_id) WHERE (is_global IS TRUE);

CREATE UNIQUE INDEX shop_connection_shop_id_connection_id_idx ON public.shop_connection USING btree (shop_id, connection_id) WHERE (deleted_at IS NULL);

CREATE INDEX shop_customer_created_at_idx ON public.shop_customer USING btree (created_at);

CREATE INDEX shop_customer_full_name_norm_idx ON public.shop_customer USING gin (full_name_norm);

CREATE INDEX shop_customer_group_created_at_idx ON public.shop_customer_group USING btree (created_at);

CREATE INDEX shop_customer_group_customer_created_at_idx ON public.shop_customer_group_customer USING btree (created_at);

CREATE UNIQUE INDEX shop_customer_partner_id_external_id_idx ON public.shop_customer USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_customer_partner_id_shop_id_external_code_idx ON public.shop_customer USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));

CREATE INDEX shop_customer_phone_idx ON public.shop_customer USING btree (phone);

CREATE INDEX shop_customer_phone_norm_idx ON public.shop_customer USING gin (phone_norm);

CREATE UNIQUE INDEX shop_customer_shop_id_code_idx ON public.shop_customer USING btree (shop_id, code);

CREATE INDEX shop_customer_shop_id_email_idx ON public.shop_customer USING btree (shop_id, email);

CREATE INDEX shop_customer_shop_id_phone_idx ON public.shop_customer USING btree (shop_id, phone);

CREATE INDEX shop_customer_type_idx ON public.shop_customer USING btree (type);

CREATE INDEX shop_customer_updated_at_id_idx ON public.shop_customer USING btree (updated_at, id);

CREATE INDEX shop_ledger_bank_account_idx ON public.shop_ledger USING gin (bank_account);

CREATE INDEX shop_ledger_bank_account_idx1 ON public.shop_ledger USING gin (bank_account);

CREATE UNIQUE INDEX shop_product_collection_partner_id_external_collection_id_e_idx ON public.shop_product_collection USING btree (partner_id, external_collection_id, external_product_id) WHERE ((partner_id IS NOT NULL) AND (external_collection_id IS NOT NULL) AND (external_product_id IS NOT NULL));

CREATE INDEX shop_product_list_price_idx ON public.shop_product USING btree (list_price);

CREATE INDEX shop_product_name_norm_idx ON public.shop_product USING gin (name_norm);

CREATE UNIQUE INDEX shop_product_partner_id_external_id_idx ON public.shop_product USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_product_partner_id_external_id_idx1 ON public.shop_product USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_product_partner_id_shop_id_external_code_idx ON public.shop_product USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));

CREATE UNIQUE INDEX shop_product_partner_id_shop_id_external_code_idx1 ON public.shop_product USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));

CREATE INDEX shop_product_search_idx ON public.shop_product USING gin (name_norm);

CREATE UNIQUE INDEX shop_product_shop_id_code_idx ON public.shop_product USING btree (shop_id, code) WHERE (deleted_at IS NULL);

CREATE UNIQUE INDEX shop_product_shop_id_product_id_idx ON public.shop_product USING btree (shop_id, product_id);

CREATE INDEX shop_product_source_id_idx ON public.shop USING btree (product_source_id);

CREATE INDEX shop_product_updated_at_product_id_idx ON public.shop_product USING btree (updated_at, product_id);

CREATE INDEX shop_product_updated_at_product_id_idx1 ON public.shop_product USING btree (updated_at, product_id);

CREATE INDEX shop_stocktake_code_idx ON public.shop_stocktake USING btree (code);

CREATE INDEX shop_stocktake_status_idx ON public.shop_stocktake USING btree (status);

CREATE INDEX shop_supplier_company_name_norm_idx ON public.shop_supplier USING gin (company_name_norm);

CREATE INDEX shop_supplier_email_idx ON public.shop_supplier USING btree (email);

CREATE INDEX shop_supplier_full_name_norm_idx ON public.shop_supplier USING gin (full_name_norm);

CREATE INDEX shop_supplier_phone_idx ON public.shop_supplier USING btree (phone);

CREATE INDEX shop_supplier_phone_norm_idx ON public.shop_supplier USING gin (phone_norm);

CREATE UNIQUE INDEX shop_supplier_shop_id_code_idx ON public.shop_supplier USING btree (shop_id, code);

CREATE INDEX shop_variant_created_at_idx ON public.shop_variant USING btree (created_at);

CREATE UNIQUE INDEX shop_variant_partner_id_external_id_idx ON public.shop_variant USING btree (partner_id, external_id) WHERE ((partner_id IS NOT NULL) AND (external_id IS NOT NULL));

CREATE UNIQUE INDEX shop_variant_partner_id_shop_id_external_code_idx ON public.shop_variant USING btree (partner_id, shop_id, external_code) WHERE ((partner_id IS NOT NULL) AND (external_code IS NOT NULL));

CREATE INDEX shop_variant_product_id_idx ON public.shop_variant USING btree (product_id);

CREATE UNIQUE INDEX shop_variant_shop_id_code_idx ON public.shop_variant USING btree (shop_id, code) WHERE (deleted_at IS NULL);

CREATE UNIQUE INDEX shop_variant_supplier_shop_id_supplier_id_variant_id_idx ON public.shop_variant_supplier USING btree (shop_id, supplier_id, variant_id);

CREATE INDEX shop_variant_supplier_supplier_id_idx ON public.shop_variant_supplier USING btree (supplier_id);

CREATE INDEX shop_variant_supplier_variant_id_idx ON public.shop_variant_supplier USING btree (variant_id);

CREATE INDEX shop_variant_updated_at_variant_id_idx ON public.shop_variant USING btree (updated_at, variant_id);

CREATE UNIQUE INDEX user_email_key ON public."user" USING btree (email) WHERE (wl_partner_id IS NULL);

CREATE UNIQUE INDEX user_email_wl_partner_id_idx ON public."user" USING btree (email, wl_partner_id);

CREATE UNIQUE INDEX user_phone_key ON public."user" USING btree (phone) WHERE (wl_partner_id IS NULL);
