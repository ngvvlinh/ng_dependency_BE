package admin

import (
	"o.o/capi/httprpc"
	"o.o/common/l"
)

// +gen:wrapper=o.o/api/top/int/admin
// +gen:wrapper:package=admin

var ll = l.New()

type Servers []httprpc.Server
