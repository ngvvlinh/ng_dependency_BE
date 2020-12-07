package authetop

import (
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authcommon"
)

const Policy auth.Policy = authcommon.CommonPolicy + `
	# refund
	p, shop/extension:create, admin, owner, staff_management
	p, shop/extension:delete, admin, owner, staff_management
	p, shop/extension:view, admin, owner, staff_management
	p, shop/hotline:view, admin, owner, staff_management
	p, shop/calllog:view, admin, owner, staff_management
`
