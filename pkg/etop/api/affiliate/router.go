package affiliate

import (
	service "o.o/api/top/int/affiliate"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/affiliate
// +gen:wrapper:package=affiliate

type Servers []httprpc.Server

func NewServers(
	miscService MiscService,
	accountService AccountService,
) Servers {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
	}
	return servers
}
