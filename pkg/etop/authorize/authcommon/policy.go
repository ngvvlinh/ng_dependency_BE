package authcommon

const shopPolicy = `#connection
	p, shop/connection:create, owner
	p, shop/connection:update, owner
	p, shop/connection:delete, owner
	p, shop/connection:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# carrier
    p, shop/carrier:create, admin, owner, salesman
	p, shop/carrier:update, admin, owner, salesman
	p, shop/carrier:delete, admin, owner, salesman
	p, shop/carrier:view, admin, owner, salesman, accountant
	# category
	p, shop/category:create, admin, owner
	p, shop/category:update, admin, owner, purchasing_management
    p, shop/category:delete, admin, owner
    p, shop/category:view, admin, owner, salesman, accountant, purchasing_management, inventory_management
	# collection
	p, shop/collection:create, admin, owner, purchasing_management
	p, shop/collection:update, admin, owner, purchasing_management
    p, shop/collection:view, admin, owner, accountant, salesman, purchasing_management, inventory_management
	# customer
	p, shop/customer:create, admin, owner, salesman
	p, shop/customer:update, admin, owner, salesman
	p, shop/customer:delete, admin, owner, salesman
	p, shop/customer:view, admin, owner, salesman, accountant
	p, shop/customer:manage, admin, owner, salesman
	p, shop/customer_group:manage, admin, owner, salesman
    # dashboard
    p, shop/dashboard:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# external account
    p, shop/external_account:manage, admin, owner
    # fulfillment
	p, shop/fulfillment:create, admin, owner, salesman
	p, shop/fulfillment:update, admin, owner, salesman
	p, shop/fulfillment:cancel, admin, owner, salesman
	p, shop/fulfillment:view, admin, owner, salesman, accountant
	p, shop/fulfillment:export, admin, owner, salesman
	# inventory
	p, shop/inventory:view, admin, owner, inventory_management, purchasing_management
	p, shop/inventory:update, admin, owner, inventory_management, purchasing_management
    p, shop/inventory:create, admin, owner, salesman, inventory_management, purchasing_management
    p, shop/inventory:confirm, admin, owner, inventory_management, salesman
    p, shop/inventory:cancel, admin, owner, inventory_management
	# ledger
	p, shop/ledger:create, admin, owner, accountant
	p, shop/ledger:view, admin, owner, accountant, salesman, purchasing_management
	p, shop/ledger:update, admin, owner, accountant
	p, shop/ledger:delete, admin, owner, accountant
	# money_transaction
	p, shop/money_transaction:view, admin, owner, accountant
	p, shop/money_transaction:export, admin, owner, accountant
	# order
	p, shop/order:create, admin, owner, salesman
	p, shop/order:confirm, admin, owner, salesman
	p, shop/order:complete, admin, owner, salesman
	p, shop/order:update, admin, owner, salesman
	p, shop/order:cancel, admin, owner, salesman
	p, shop/order:import, admin, owner, salesman
	p, shop/order:export, admin, owner, salesman, inventory_management
    # order: hotfix, all users can access order page
	p, shop/order:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# purchase_order
	p, shop/purchase_order:view, admin, owner, accountant, purchasing_management
	p, shop/purchase_order:create, admin, owner, purchasing_management
    p, shop/purchase_order:update, admin, owner, purchasing_management
	p, shop/purchase_order:confirm, admin, owner, purchasing_management
    p, shop/purchase_order:cancel, admin, owner, purchasing_management
	# product
	p, shop/product:create, admin, owner, purchasing_management
	p, shop/product/basic_info:update, admin, owner, purchasing_management
	p, shop/product:delete, admin, owner, purchasing_management
	p, shop/product:import, admin, owner, purchasing_management
	p, shop/product/basic_info:view, admin, owner, salesman, purchasing_management, inventory_management
	p, shop/product/cost_price:view, admin, owner, purchasing_management
	p, shop/product/cost_price:update, admin, owner, purchasing_management
	p, shop/product/retail_price:update, admin, owner
	# receipt
	p, shop/receipt:create, admin, owner, salesman, accountant, purchasing_management
	p, shop/receipt:update, admin, owner, accountant
	p, shop/receipt:confirm, admin, owner, salesman, accountant, purchasing_management
	p, shop/receipt:cancel, admin, owner, accountant
	p, shop/receipt:view, admin, owner, accountant
	# settings
	p, shop/settings/shop_info:update, admin, owner
	p, shop/settings/shop_info:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	p, shop/settings/company_info:view, admin, owner
	p, shop/settings/company_info:update, admin, owner
	p, shop/settings/bank_info:view, admin, owner
	p, shop/settings/bank_info:update, admin, owner
	p, shop/settings/payment_schedule:view, admin, owner
	p, shop/settings/payment_schedule:update, admin, owner
	p, shop/settings/topship_balance:view, admin, owner
	p, shop/settings/wallet_balance:view, admin, owner
	p, shop/settings/shipping_setting:view, admin, owner
	p, shop/settings/shipping_setting:update, admin, owner
	# shipnow FFM
	p, shop/shipnow:view, admin, owner, salesman, accountant
	p, shop/shipnow:create, admin, owner, salesman
	p, shop/shipnow:update, admin, owner, salesman
    p, shop/shipnow:cancel, admin, owner, salesman
	p, shop/shipnow:confirm, admin, owner, salesman
	# shop
	p, shop/balance:view, admin, owner, accountant
	p, shop/account:delete, owner
	p, shop/user/balance:view, owner, accountant
	# staff
	p, relationship/invitation:create, admin, owner, staff_management
	p, relationship/invitation:view, admin, owner, staff_management
	p, relationship/invitation:delete, admin, owner, staff_management
	p, relationship/permission:update, admin, owner, staff_management
	p, relationship/relationship:update, admin, owner, staff_management
	p, relationship/relationship:view, admin, owner, staff_management
	p, relationship/relationship:remove, admin, owner, staff_management
	# stocktake
	p, shop/stocktake:create, admin, owner, purchasing_management, inventory_management
	p, shop/stocktake:update, admin, owner, inventory_management
	p, shop/stocktake:confirm, admin, owner, inventory_management
	p, shop/stocktake:cancel, admin, owner, inventory_management
	p, shop/stocktake:view, admin, owner, accountant, purchasing_management, inventory_management
	p, shop/stocktake:self_update, purchasing_management
	# supplier
	p, shop/supplier:create, admin, owner, purchasing_management
	p, shop/supplier:update, admin, owner, purchasing_management
	p, shop/supplier:delete, admin, owner, purchasing_management
	p, shop/supplier:view, admin, owner, accountant, purchasing_management
	# refund
	p, shop/refund:create, admin, owner, salesman
	p, shop/refund:update, admin, owner, salesman
	p, shop/refund:delete, admin, owner, salesman
	p, shop/refund:view, admin, owner, salesman
	p, shop/refund:confirm, admin, owner, salesman
	p, shop/refund:cancel, admin, owner, salesman
	# purchaserefund
	p, shop/purchaserefund:create, admin, owner, purchasing_management
	p, shop/purchaserefund:update, admin, owner, purchasing_management
	p, shop/purchaserefund:delete, admin, owner, purchasing_management
	p, shop/purchaserefund:view, admin, owner, purchasing_management
	p, shop/purchaserefund:confirm, admin, owner, purchasing_management
	p, shop/purchaserefund:cancel, admin, owner, purchasing_management
	# trading
	p, trading/order:view, admin, owner
	p, trading/order:create, admin, owner
    p, shop/trading/product:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
    p, shop/trading/order:create, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
    p, shop/trading/order:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management
	# payment
	p, shop/payment:create, admin, owner
	p, shop/payment:view, admin, owner
	# webserver
	p, shop/webserver/wswebsite:create, admin, owner
	p, shop/webserver/wswebsite:update, admin, owner
	p, shop/webserver/wswebsite:view, admin, owner
	p, shop/webserver/wsproduct:create, admin, owner
	p, shop/webserver/wsproduct:update, admin, owner
	p, shop/webserver/wsproduct:view, admin, owner
	p, shop/webserver/wscategory:create, admin, owner
	p, shop/webserver/wscategory:update, admin, owner
	p, shop/webserver/wscategory:view, admin, owner
	p, shop/webserver/wspage:create, admin, owner
	p, shop/webserver/wspage:update, admin, owner
	p, shop/webserver/wspage:delete, admin, owner
	p, shop/webserver/wspage:view, admin, owner
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
	# credit
	p, admin/credit:view, admin, ad_accountant, ad_customerservice, ad_customerservice_lead
	p, admin/credit:update, admin, ad_accountant
	p, admin/credit:create, admin, ad_accountant
	p, admin/credit:confirm, admin, ad_accountant
	p, admin/credit:delete, admin, ad_accountant
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
	p, admin/shop:view, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_accountant
	p, admin/shop:update, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale
	# admin user
	p, admin/user:view, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale, ad_accountant
	p, admin/user:block, admin, ad_customerservice, ad_customerservice_lead, ad_salelead, ad_sale
	p, admin/user_ref:update, admin, ad_salelead
	# admin subscription
	p, admin/subscription_product:create, admin
	p, admin/subscription_product:view, admin
	p, admin/subscription_product:delete, admin
	p, admin/subscription_plan:create, admin
	p, admin/subscription_plan:update, admin
	p, admin/subscription_plan:view, admin
	p, admin/subscription_plan:delete, admin
	p, admin/subscription:view, admin
	p, admin/subscription:create, admin
	p, admin/subscription:update, admin
	p, admin/subscription:cancel, admin
	p, admin/subscription:active, admin
	p, admin/subscription:delete, admin
	p, admin/subscription_bill:view, admin
	p, admin/subscription_bill:create, admin
	p, admin/subscription_bill_manual_payment:create, admin
	p, admin/subscription_bill:delete, admin
	# admin misc
	p, admin/misc_account:login, admin
	# admin partner
	p, admin/partner:create, admin
	# admin connection
	p, admin/connection:view, admin, ad_sale, ad_salelead, ad_customerservice, ad_accountant, ad_customerservice_lead
	p, admin/connection:confirm, admin
	p, admin/connection:disable, admin
	p, admin/connection_builtin:create, admin
	p, admin/connection_shop_builtin:create, admin, ad_sale, ad_salelead
	p, admin/connection_shop_builtin:update, admin
	p, admin/connection_service:view, admin, ad_sale, ad_salelead
	# admin custom region
	p, admin/custom_region:create, admin, ad_salelead
	p, admin/custom_region:update, admin, ad_salelead
	p, admin/custom_region:delete, admin, ad_salelead
	p, admin/custom_region:view, admin, ad_salelead
`

const CommonPolicy = shopPolicy + adminPolicy
