package authcommon

const shopPolicy = `#connection
	p, shop/connection:create, owner, m_admin
	p, shop/connection:update, owner, m_admin
	p, shop/connection:delete, owner, m_admin
	p, shop/connection:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	# carrier
	p, shop/carrier:create, admin, owner, salesman, m_admin
	p, shop/carrier:update, admin, owner, salesman, m_admin
	p, shop/carrier:delete, admin, owner, salesman, m_admin
	p, shop/carrier:view, admin, owner, salesman, accountant, m_admin
	# category
	p, shop/category:create, admin, owner, m_admin
	p, shop/category:update, admin, owner, purchasing_management, m_admin
	p, shop/category:delete, admin, owner, m_admin
	p, shop/category:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, m_admin
	# collection
	p, shop/collection:create, admin, owner, purchasing_management, m_admin
	p, shop/collection:update, admin, owner, purchasing_management, m_admin
	p, shop/collection:view, admin, owner, accountant, salesman, purchasing_management, inventory_management, m_admin
	# customer
	p, shop/customer:create, admin, owner, salesman, m_admin
	p, shop/customer:update, admin, owner, salesman, m_admin
	p, shop/customer:delete, admin, owner, salesman, m_admin
	p, shop/customer:view, admin, owner, salesman, accountant, m_admin
	p, shop/customer:manage, admin, owner, salesman, m_admin
	p, shop/customer_group:manage, admin, owner, salesman, m_admin
	# external account
	p, shop/external_account:manage, admin, owner, m_admin
	# fulfillment
	p, shop/fulfillment:create, admin, owner, salesman, m_admin
	p, shop/fulfillment:update, admin, owner, salesman, m_admin
	p, shop/fulfillment:cancel, admin, owner, salesman, m_admin
	p, shop/fulfillment:view, admin, owner, salesman, accountant, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/fulfillment:export, admin, owner, salesman, m_admin
	# inventory
	p, shop/inventory:view, admin, owner, inventory_management, purchasing_management, m_admin
	p, shop/inventory:update, admin, owner, inventory_management, purchasing_management, m_admin
	p, shop/inventory:create, admin, owner, salesman, inventory_management, purchasing_management, m_admin
	p, shop/inventory:confirm, admin, owner, inventory_management, salesman, m_admin
	p, shop/inventory:cancel, admin, owner, inventory_management, m_admin
	# ledger
	p, shop/ledger:create, admin, owner, accountant, m_admin
	p, shop/ledger:view, admin, owner, accountant, salesman, purchasing_management, m_admin
	p, shop/ledger:update, admin, owner, accountant, m_admin
	p, shop/ledger:delete, admin, owner, accountant, m_admin
	# money_transaction
	p, shop/money_transaction:view, admin, owner, accountant, m_admin
	p, shop/money_transaction:export, admin, owner, accountant, m_admin
	# order
	p, shop/order:create, admin, owner, salesman, m_admin
	p, shop/order:confirm, admin, owner, salesman, m_admin
	p, shop/order:complete, admin, owner, salesman, m_admin
	p, shop/order:update, admin, owner, salesman, m_admin
	p, shop/order:cancel, admin, owner, salesman, m_admin
	p, shop/order:import, admin, owner, salesman, m_admin
	p, shop/order:export, admin, owner, salesman, inventory_management, m_admin
	# order: hotfix, all users can access order page
	p, shop/order:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	# purchase_order
	p, shop/purchase_order:view, admin, owner, accountant, purchasing_management, m_admin
	p, shop/purchase_order:create, admin, owner, purchasing_management, m_admin
	p, shop/purchase_order:update, admin, owner, purchasing_management, m_admin
	p, shop/purchase_order:confirm, admin, owner, purchasing_management, m_admin
	p, shop/purchase_order:cancel, admin, owner, purchasing_management, m_admin
	# receipt
	p, shop/receipt:create, admin, owner, salesman, accountant, purchasing_management, m_admin
	p, shop/receipt:update, admin, owner, accountant, m_admin
	p, shop/receipt:confirm, admin, owner, salesman, accountant, purchasing_management, m_admin
	p, shop/receipt:cancel, admin, owner, accountant, m_admin
	p, shop/receipt:view, admin, owner, accountant, m_admin
	# settings
	p, shop/settings/shop_info:update, admin, owner, m_admin
	p, shop/settings/shop_info:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/settings/company_info:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/settings/company_info:update, admin, owner, m_admin
	p, shop/settings/bank_info:view, admin, owner, m_admin
	p, shop/settings/bank_info:update, admin, owner, m_admin
	p, shop/settings/payment_schedule:view, admin, owner, m_admin
	p, shop/settings/payment_schedule:update, admin, owner, m_admin
	p, shop/settings/topship_balance:view, admin, owner, m_admin
	p, shop/settings/wallet_balance:view, admin, owner, m_admin
	p, shop/settings/shipping_setting:view, admin, owner, m_admin
	p, shop/settings/shipping_setting:update, admin, owner, m_admin
	# shipnow FFM
	p, shop/shipnow:view, admin, owner, salesman, accountant, m_admin
	p, shop/shipnow:create, admin, owner, salesman, m_admin
	p, shop/shipnow:update, admin, owner, salesman, m_admin
	p, shop/shipnow:cancel, admin, owner, salesman, m_admin
	p, shop/shipnow:confirm, admin, owner, salesman, m_admin
	# shop
	p, shop/balance:view, admin, owner, accountant, m_admin
	p, shop/account:delete, owner, m_admin
	p, shop/user/balance:view, owner, accountant, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	# staff
	p, relationship/invitation:create, admin, owner, staff_management, m_admin
	p, relationship/invitation:view, admin, owner, staff_management, m_admin
	p, relationship/invitation:delete, admin, owner, staff_management, m_admin
	p, relationship/permission:update, admin, owner, staff_management, m_admin
	p, relationship/relationship:update, admin, owner, staff_management, m_admin
	p, relationship/relationship:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	p, relationship/relationship:remove, admin, owner, staff_management, m_admin
	# stocktake
	p, shop/stocktake:create, admin, owner, purchasing_management, inventory_management, m_admin
	p, shop/stocktake:update, admin, owner, inventory_management, m_admin
	p, shop/stocktake:confirm, admin, owner, inventory_management, m_admin
	p, shop/stocktake:cancel, admin, owner, inventory_management, m_admin
	p, shop/stocktake:view, admin, owner, accountant, purchasing_management, inventory_management, m_admin
	p, shop/stocktake:self_update, purchasing_management, m_admin
	# supplier
	p, shop/supplier:create, admin, owner, purchasing_management, m_admin
	p, shop/supplier:update, admin, owner, purchasing_management, m_admin
	p, shop/supplier:delete, admin, owner, purchasing_management, m_admin
	p, shop/supplier:view, admin, owner, accountant, purchasing_management, m_admin
	# refund
	p, shop/refund:create, admin, owner, salesman, m_admin
	p, shop/refund:update, admin, owner, salesman, m_admin
	p, shop/refund:delete, admin, owner, salesman, m_admin
	p, shop/refund:view, admin, owner, salesman, m_admin
	p, shop/refund:confirm, admin, owner, salesman, m_admin
	p, shop/refund:cancel, admin, owner, salesman, m_admin
	# purchaserefund
	p, shop/purchaserefund:create, admin, owner, purchasing_management, m_admin
	p, shop/purchaserefund:update, admin, owner, purchasing_management, m_admin
	p, shop/purchaserefund:delete, admin, owner, purchasing_management, m_admin
	p, shop/purchaserefund:view, admin, owner, purchasing_management, m_admin
	p, shop/purchaserefund:confirm, admin, owner, purchasing_management, m_admin
	p, shop/purchaserefund:cancel, admin, owner, purchasing_management, m_admin
	# trading
	p, trading/order:view, admin, owner, m_admin
	p, trading/order:create, admin, owner, m_admin
	p, shop/trading/product:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	p, shop/trading/order:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	p, shop/trading/order:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, m_admin
	# payment
	p, shop/payment:create, admin, owner, m_admin
	p, shop/payment:view, admin, owner, m_admin
	# webserver
	p, shop/webserver/wswebsite:create, admin, owner, m_admin
	p, shop/webserver/wswebsite:update, admin, owner, m_admin
	p, shop/webserver/wswebsite:view, admin, owner, m_admin
	p, shop/webserver/wsproduct:create, admin, owner, m_admin
	p, shop/webserver/wsproduct:update, admin, owner, m_admin
	p, shop/webserver/wsproduct:view, admin, owner, m_admin
	p, shop/webserver/wscategory:create, admin, owner, m_admin
	p, shop/webserver/wscategory:update, admin, owner, m_admin
	p, shop/webserver/wscategory:view, admin, owner, m_admin
	p, shop/webserver/wspage:create, admin, owner, m_admin
	p, shop/webserver/wspage:update, admin, owner, m_admin
	p, shop/webserver/wspage:delete, admin, owner, m_admin
	p, shop/webserver/wspage:view, admin, owner, m_admin
	# ticket
	p, shop/shop_ticket:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket:update, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket:assign, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket:reopen, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_comment:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_comment:update, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_comment:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_comment:delete, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_label:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_label:update, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/shop_ticket_label:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	# subscription
	p, shop/subscription_product:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/subscription_plan:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/subscription:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/subscription:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	p, shop/subscription:update, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, telecom_customerservice, telecom_customerservice_management, m_admin
	# credit
	p, shop/credit:create, admin, owner, accountant, staff_management, telecom_customerservice, telecom_customerservice_management
	# transaction, m_admin
	p, shop/transaction:view, admin, owner, salesman, accountant, purchasing_management, m_admin
	# invoice
	p, shop/invoice:view, admin, owner, accountant, salesman, m_admin
	# jira
	p, shop/jira:view, admin, owner, staff_management, salesman, m_admin
	p, shop/jira:create, admin, owner, staff_management, salesman, m_admin
	# account_user
    p, shop/account_user:create, admin, owner, staff_management, m_admin
	p, shop/account_user:update, admin, owner, staff_management, m_admin
	p, shop/account_user:delete, admin, owner, staff_management, m_admin
	p, shop/account_user:view, admin, owner, staff_management, telecom_customerservice, inventory_management, purchasing_management, accountant, telecom_customerservice_management, m_admin
`

const adminPolicy = `
	# account
	p, admin/admin_user:create, admin
	p, admin/admin_user:update, admin
	p, admin/admin_user:view, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_user:delete, admin
	# ticket
	p, admin/admin_ticket:assign, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_ticket:create, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_ticket:update, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_ticket:delete, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_ticket:view, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/admin_lead_ticket:update, admin, ad_salelead, ad_customerservice_lead
	# ticket_label
	p, admin/admin_ticket_label:create, admin, ad_salelead, ad_customerservice_lead
	p, admin/admin_ticket_label:update, admin, ad_salelead, ad_customerservice_lead
	p, admin/admin_ticket_label:view, admin, ad_salelead, ad_customerservice_lead
	# credit
	p, admin/credit:view, admin, ad_accountant, ad_customerservice, ad_customerservice_lead, ad_voip
	p, admin/credit:update, admin, ad_accountant, ad_voip
	p, admin/credit:create, admin, ad_accountant, ad_voip
	p, admin/credit:confirm, admin, ad_accountant, ad_voip
	p, admin/credit:delete, admin, ad_accountant, ad_voip
	# fulfillment
	p, admin/fulfillment:view, admin, ad_accountant, ad_salelead, ad_sale, ad_customerservice, ad_customerservice_lead
	p, admin/fulfillment:create, admin, ad_accountant, ad_customerservice, ad_customerservice_lead
	p, admin/fulfillment:update, admin, ad_customerservice, ad_customerservice_lead
	p, admin/fulfillment_state:update, admin, ad_customerservice_lead, ad_accountant
	p, admin/fulfillment_shipping_fees:update, admin, ad_accountant
	# money transaction
	p, admin/money_transaction:create, admin, ad_accountant
	p, admin/money_transaction:view, admin, ad_accountant, ad_customerservice, ad_customerservice_lead
	p, admin/money_transaction:update, admin, ad_accountant
	p, admin/money_transaction:confirm, admin, ad_accountant
	p, admin/money_transaction_shipping_etop:view, admin, ad_accountant, ad_customerservice, ad_customerservice_lead
	p, admin/money_transaction_shipping_etop:create, admin, ad_accountant
	p, admin/money_transaction_shipping_etop:update, admin, ad_accountant
	p, admin/money_transaction_shipping_etop:confirm, admin, ad_accountant
	p, admin/money_transaction_shipping_etop:delete, admin, ad_accountant
	p, admin/money_transaction_shipping_external:view, admin, ad_accountant, ad_customerservice, ad_customerservice_lead
	p, admin/money_transaction_shipping_external:create, admin, ad_accountant
	p, admin/money_transaction_shipping_external:update, admin, ad_accountant
	p, admin/money_transaction_shipping_external:confirm, admin, ad_accountant
	p, admin/money_transaction_shipping_external:delete, admin, ad_accountant
	p, admin/money_transaction_shipping_external_lines:remove, admin, ad_accountant
	# order
	p, admin/order:view, admin, ad_salelead, ad_sale, ad_customerservice, ad_accountant, ad_customerservice_lead
	# shop shipment price
	p, admin/shop_shipment_price_list:update, admin, ad_salelead, ad_sale
	p, admin/shop_shipment_price_list:create, admin, ad_salelead, ad_sale
	p, admin/shop_shipment_price_list:delete, admin, ad_salelead, ad_sale
	p, admin/shop_shipment_price_list:view, admin, ad_salelead, ad_sale
	# shipment price services
	p, admin/shipment_service:view, admin, ad_salelead
	p, admin/shipment_service:create, admin, ad_salelead
	p, admin/shipment_service:delete, admin, ad_salelead
	p, admin/shipment_service:update, admin, ad_salelead
	# shipment price list
	p, admin/shipment_price_list:create, admin, ad_salelead
	p, admin/shipment_price_list:update, admin, ad_salelead
	p, admin/shipment_price_list:delete, admin, ad_salelead
	p, admin/shipment_price_list:view, admin, ad_salelead
	# shipment price list promotion
	p, admin/shipment_price_list_promotion:create, admin, ad_salelead
	p, admin/shipment_price_list_promotion:update, admin, ad_salelead
	p, admin/shipment_price_list_promotion:delete, admin, ad_salelead
	p, admin/shipment_price_list_promotion:view, admin, ad_salelead
	# shipment price
	p, admin/shipment_price:create, admin, ad_salelead
	p, admin/shipment_price:update, admin, ad_salelead
	p, admin/shipment_price:delete, admin, ad_salelead
	p, admin/shipment_price:view, admin, ad_salelead
	# admin shop
	p, admin/shop:view, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_accountant, ad_voip
	p, admin/shop:update, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_voip
	# admin user
	p, admin/user:view, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_accountant, ad_voip
	p, admin/user:block, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_voip
	p, admin/user_ref:update, admin, ad_salelead, ad_voip
	# admin subscription
	p, admin/subscription_product:create, admin, ad_voip
	p, admin/subscription_product:view, admin, ad_voip
	p, admin/subscription_product:delete, admin, ad_voip
	p, admin/subscription_plan:create, admin, ad_voip
	p, admin/subscription_plan:update, admin, ad_voip
	p, admin/subscription_plan:view, admin, ad_voip
	p, admin/subscription_plan:delete, admin, ad_voip
	p, admin/subscription:view, admin, ad_voip
	p, admin/subscription:create, admin, ad_voip
	p, admin/subscription:update, admin, ad_voip
	p, admin/subscription:cancel, admin, ad_voip
	p, admin/subscription:active, admin, ad_voip
	p, admin/subscription:delete, admin, ad_voip
	p, admin/invoice:view, admin, ad_voip
	p, admin/invoice:create, admin, ad_voip
	p, admin/invoice_manual_payment:create, admin, ad_voip
	p, admin/invoice:delete, admin, ad_voip
	# admin misc
	p, admin/misc_account:login, admin
	# admin partner
	p, admin/partner:create, admin
	# admin connection
	p, admin/connection:view, admin, ad_sale, ad_salelead, ad_customerservice, ad_accountant, ad_customerservice_lead, ad_voip
	p, admin/connection:confirm, admin
	p, admin/connection:disable, admin
	p, admin/connection_builtin:create, admin
	p, admin/shop_connection_builtin:create, admin, ad_sale, ad_salelead
	p, admin/shop_connection:update, admin
	p, admin/connection_service:view, admin, ad_sale, ad_salelead
	# admin custom region
	p, admin/custom_region:create, admin, ad_salelead
	p, admin/custom_region:update, admin, ad_salelead
	p, admin/custom_region:delete, admin, ad_salelead
	p, admin/custom_region:view, admin, ad_salelead
	# admin transaction
	p, admin/transaction:view, admin, ad_salelead, ad_accountant, ad_voip
`

const CommonPolicy = shopPolicy + adminPolicy
