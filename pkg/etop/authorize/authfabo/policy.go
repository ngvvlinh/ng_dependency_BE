package authfabo

import (
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authcommon"
)

const Policy auth.Policy = authcommon.CommonPolicy + `
	# product
	p, shop/product:create, admin, owner, purchasing_management, salesman
	p, shop/product/basic_info:update, admin, owner, purchasing_management, salesman
	p, shop/product:delete, admin, owner, purchasing_management, salesman
	p, shop/product:import, admin, owner, purchasing_management, salesman
	p, shop/product/basic_info:view, admin, owner, salesman, purchasing_management, inventory_management
	p, shop/product/cost_price:view, admin, owner, purchasing_management, salesman
	p, shop/product/cost_price:update, admin, owner, purchasing_management, salesman
	p, shop/product/retail_price:update, admin, owner, salesman
	# facebook conversation
	p, facebook/comment:view, admin, owner, salesman
	p, facebook/comment:create, admin, owner, salesman
	p, facebook/post:create, admin, owner, salesman
	p, facebook/message:view, admin, owner, salesman
	p, facebook/message:create, admin, owner, salesman
	p, facebook/fbuser:view, admin, owner, salesman
	p, facebook/fbuser:update, admin, owner, salesman
	p, facebook/fbuser:create, admin, owner, salesman
	p, facebook/fanpage:create, owner
	p, facebook/fanpage:delete, owner
	p, facebook/fanpage:view, owner, salesman
	p, facebook/shoptag:create, owner, salesman
	p, facebook/shoptag:update, owner, salesman
	p, facebook/shoptag:view, owner, salesman
	p, facebook/shoptag:delete, owner, salesman
	p, facebook/message_template:create, owner, salesman
	p, facebook/message_template:update, owner, salesman
	p, facebook/message_template:view, owner, salesman
	p, facebook/message_template:delete, owner, salesman
	# dashboard
	p, shop/dashboard:view, owner
`
