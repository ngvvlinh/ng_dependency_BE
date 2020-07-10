package affiliate

import (
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(
	miscService MiscService,
	accountService AccountService,
) Servers {
	servers := httprpc.MustNewServers(
		miscService.Clone,
		accountService.Clone,
	)
	return servers
}
