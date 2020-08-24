package vnp

import (
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(
	shipnowService *VNPostService,
) Servers {
	servers := httprpc.MustNewServers(
		shipnowService.Clone,
	)
	return servers
}
