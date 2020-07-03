package authcommon

const CommonPolicy = `#connection
	p, shop/connection:create, owner
	p, shop/connection:update, owner
	p, shop/connection:delete, owner
	p, shop/connection:view, admin, owner, salesman, accountant, purchasing_management, inventory_management, staff_management, analyst
	# carrier
    p, shop/carrier:create, admin, owner, salesman
	p, shop/carrier:update, admin, owner, salesman
	p, shop/carrier:delete, admin, owner, salesman
	p, shop/carrier:view, admin, owner, salesman, accountant
	# category
	p, shop/category:create, admin, owner, salesman
	p, shop/category:update, admin, owner, salesman, purchasing_management
    p, shop/category:delete, admin, owner, salesman
    p, shop/category:view, admin, owner, salesman, accountant, purchasing_management, inventory_management
	# collection
	p, shop/collection:create, admin, owner, salesman, purchasing_management
	p, shop/collection:update, admin, owner, salesman, purchasing_management
    p, shop/collection:view, admin, owner, salesman, accountant, purchasing_management, inventory_management
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
	p, shop/money_transaction:export, admin, owner
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
	p, shop/webserver/wspage:delete, admin, owner,
	p, shop/webserver/wspage:view, admin, owner`