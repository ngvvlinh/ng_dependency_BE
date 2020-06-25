package authfabo

import "o.o/backend/pkg/etop/authorize/authcommon"

const Policy = authcommon.CommonPolicy + `
	p, facebook/comment:view, admin, owner, salesman
	p, facebook/comment:create, admin, owner, salesman
	p, facebook/message:view, admin, owner, salesman
	p, facebook/message:create, admin, owner, salesman
	p, facebook/fbuser:view, admin, owner, salesman
	p, facebook/fbuser:create, admin, owner, salesman
`
