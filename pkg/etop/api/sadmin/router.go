package sadmin

import (
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(
	miscService *MiscService,
	userService *UserService,
) Servers {
	servers := httprpc.MustNewServers(
		miscService.Clone,
		userService.Clone,
	)
	return servers
}
