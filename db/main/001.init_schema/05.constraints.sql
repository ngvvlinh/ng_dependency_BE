ALTER TABLE ONLY public.account_auth
  ADD CONSTRAINT account_auth_pkey PRIMARY KEY (auth_key);

ALTER TABLE ONLY public.account
  ADD CONSTRAINT account_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.address
  ADD CONSTRAINT address_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.code
  ADD CONSTRAINT code_pkey PRIMARY KEY (code, type);

ALTER TABLE ONLY public.connection
  ADD CONSTRAINT connection_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.credit
  ADD CONSTRAINT credit_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.district
  ADD CONSTRAINT district_pkey PRIMARY KEY (code);

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_phone_key UNIQUE (phone);

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.money_transaction_shipping_etop
  ADD CONSTRAINT money_transaction_shipping_etop_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.money_transaction_shipping_external_line
  ADD CONSTRAINT money_transaction_shipping_external_line_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.money_transaction_shipping_external
  ADD CONSTRAINT money_transaction_shipping_external_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.partner_relation
  ADD CONSTRAINT partner_relation_pkey PRIMARY KEY (partner_id, subject_id, subject_type);

ALTER TABLE ONLY public.payment
  ADD CONSTRAINT payment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_pkey PRIMARY KEY (product_id, collection_id);

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT product_source_category_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.province
  ADD CONSTRAINT province_pkey PRIMARY KEY (code);

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.purchase_refund
  ADD CONSTRAINT purchase_refund_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.refund
  ADD CONSTRAINT refund_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shipping_source
  ADD CONSTRAINT shipping_source_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shipping_source_internal
  ADD CONSTRAINT shipping_source_state_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_code_key UNIQUE (code);

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_constraint PRIMARY KEY (group_id, customer_id);

ALTER TABLE ONLY public.shop_customer_group
  ADD CONSTRAINT shop_customer_group_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_constraint PRIMARY KEY (product_id, collection_id);

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_pkey PRIMARY KEY (product_id);

ALTER TABLE ONLY public.shop_stocktake
  ADD CONSTRAINT shop_stocktake_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_trader
  ADD CONSTRAINT shop_trader_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_pkey PRIMARY KEY (variant_id);

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.transaction
  ADD CONSTRAINT transaction_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.user_auth
  ADD CONSTRAINT user_auth_pkey PRIMARY KEY (auth_type, auth_key);

ALTER TABLE ONLY public.user_internal
  ADD CONSTRAINT user_internal_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.webhook
  ADD CONSTRAINT webhook_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.account_auth
  ADD CONSTRAINT account_auth_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.account
  ADD CONSTRAINT account_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_disabled_by_fkey FOREIGN KEY (disabled_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_invitation_sent_by_fkey FOREIGN KEY (invitation_sent_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.account_user
  ADD CONSTRAINT account_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.address
  ADD CONSTRAINT address_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);

ALTER TABLE ONLY public.affiliate
  ADD CONSTRAINT affiliate_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.connection
  ADD CONSTRAINT connection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.credit
  ADD CONSTRAINT credit_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.district
  ADD CONSTRAINT district_province_code_fkey FOREIGN KEY (province_code) REFERENCES public.province(code);

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.export_attempt
  ADD CONSTRAINT export_attempt_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.external_account_ahamove
  ADD CONSTRAINT external_account_ahamove_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_connection_id_fkey FOREIGN KEY (connection_id) REFERENCES public.connection(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_money_transaction_id_fkey FOREIGN KEY (money_transaction_id) REFERENCES public.money_transaction_shipping(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_money_transaction_shipping_external_id_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_order_id_fkey FOREIGN KEY (order_id) REFERENCES public."order"(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_shop_carrier_id_fkey FOREIGN KEY (shop_carrier_id) REFERENCES public.shop_carrier(id);

ALTER TABLE ONLY public.fulfillment
  ADD CONSTRAINT fulfillment_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.import_attempt
  ADD CONSTRAINT import_attempt_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.inventory_voucher
  ADD CONSTRAINT inventory_voucher_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.invitation
  ADD CONSTRAINT invitation_invited_by_fkey FOREIGN KEY (invited_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_money_transaction_shipping_external_id_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);

ALTER TABLE ONLY public.money_transaction_shipping_external_line
  ADD CONSTRAINT money_transaction_shipping_ex_money_transaction_shipping_e_fkey FOREIGN KEY (money_transaction_shipping_external_id) REFERENCES public.money_transaction_shipping_external(id);

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_shipping_money_transaction_shipping_etop_fkey FOREIGN KEY (money_transaction_shipping_etop_id) REFERENCES public.money_transaction_shipping_etop(id);

ALTER TABLE ONLY public.money_transaction_shipping
  ADD CONSTRAINT money_transaction_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.shop_customer(id);

ALTER TABLE ONLY public.order_line
  ADD CONSTRAINT order_line_order_id_fkey FOREIGN KEY (order_id) REFERENCES public."order"(id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payment(id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public."order"
  ADD CONSTRAINT order_trading_shop_id_fkey FOREIGN KEY (trading_shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);

ALTER TABLE ONLY public.partner
  ADD CONSTRAINT partner_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.partner_relation
  ADD CONSTRAINT partner_relation_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);

ALTER TABLE ONLY public.product_shop_collection
  ADD CONSTRAINT product_shop_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT product_source_category_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.purchase_order
  ADD CONSTRAINT purchase_order_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.shop_supplier(id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_ledger_id_fkey FOREIGN KEY (ledger_id) REFERENCES public.shop_ledger(id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_shop_ledger_id_fkey FOREIGN KEY (shop_ledger_id) REFERENCES public.shop_ledger(id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.receipt
  ADD CONSTRAINT receipt_user_id_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shipnow_fulfillment
  ADD CONSTRAINT shipnow_fulfillment_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(id);

ALTER TABLE ONLY public.shop_brand
  ADD CONSTRAINT shop_brand_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.shop_carrier
  ADD CONSTRAINT shop_carrier_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_category
  ADD CONSTRAINT shop_category_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_collection
  ADD CONSTRAINT shop_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_connection
  ADD CONSTRAINT shop_connection_connection_id_fkey FOREIGN KEY (connection_id) REFERENCES public.connection(id);

ALTER TABLE ONLY public.shop_connection
  ADD CONSTRAINT shop_connection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.shop_customer(id);

ALTER TABLE ONLY public.shop_customer_group_customer
  ADD CONSTRAINT shop_customer_group_customer_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.shop_customer_group(id);

ALTER TABLE ONLY public.shop_customer_group
  ADD CONSTRAINT shop_customer_group_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_customer
  ADD CONSTRAINT shop_customer_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_id_fkey FOREIGN KEY (id) REFERENCES public.account(id);

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_created_by_fkey FOREIGN KEY (created_by) REFERENCES public."user"(id);

ALTER TABLE ONLY public.shop_ledger
  ADD CONSTRAINT shop_ledger_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.shop_product(product_id);

ALTER TABLE ONLY public.shop_product_collection
  ADD CONSTRAINT shop_product_collection_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_product
  ADD CONSTRAINT shop_product_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_ship_from_address_id_fkey FOREIGN KEY (ship_from_address_id) REFERENCES public.address(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_ship_to_address_id_fkey FOREIGN KEY (ship_to_address_id) REFERENCES public.address(id);

ALTER TABLE ONLY public.shop_stocktake
  ADD CONSTRAINT shop_stocktake_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_trader_address
  ADD CONSTRAINT shop_trader_address_trader_id_fkey FOREIGN KEY (trader_id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_collection_id_fkey FOREIGN KEY (collection_id) REFERENCES public.shop_collection(id);

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_partner_id_fkey FOREIGN KEY (partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT shop_variant_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.shop_supplier(id);

ALTER TABLE ONLY public.shop_variant_supplier
  ADD CONSTRAINT shop_variant_supplier_variant_id_fkey FOREIGN KEY (variant_id) REFERENCES public.shop_variant(variant_id);

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_id_fkey FOREIGN KEY (id) REFERENCES public.shop_trader(id);

ALTER TABLE ONLY public.shop_supplier
  ADD CONSTRAINT shop_vendor_shop_id_fkey FOREIGN KEY (shop_id) REFERENCES public.shop(id);

ALTER TABLE ONLY public.shop
  ADD CONSTRAINT shop_wl_partner_id_fkey FOREIGN KEY (wl_partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.user_auth
  ADD CONSTRAINT user_auth_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public.user_internal
  ADD CONSTRAINT user_internal_id_fkey FOREIGN KEY (id) REFERENCES public."user"(id);

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_ref_sale_id_fkey FOREIGN KEY (ref_sale_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_ref_user_id_fkey FOREIGN KEY (ref_user_id) REFERENCES public."user"(id);

ALTER TABLE ONLY public."user"
  ADD CONSTRAINT user_wl_partner_id_fkey FOREIGN KEY (wl_partner_id) REFERENCES public.partner(id);

ALTER TABLE ONLY public.shop_variant
  ADD CONSTRAINT variant_product_id_fkey FOREIGN KEY (shop_id, product_id) REFERENCES public.shop_product(shop_id, product_id);

ALTER TABLE ONLY public.webhook
  ADD CONSTRAINT webhook_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id);

ALTER TABLE ONLY public.webhook_changes
  ADD CONSTRAINT webhook_changes_webhook_id_fkey FOREIGN KEY (webhook_id) REFERENCES public.webhook(id);
