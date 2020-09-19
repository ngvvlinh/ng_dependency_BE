package authfabo

import (
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/authcommon"
)

const Policy auth.Policy = authcommon.CommonPolicy + `
	p, facebook/comment:view, admin, owner, salesman
	p, facebook/comment:create, admin, owner, salesman
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
`
