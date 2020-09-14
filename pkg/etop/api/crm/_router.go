package crm

import (
	cc "o.o/backend/pkg/common/config"
	"o.o/capi/httprpc"
)

func NewCrmServer(m httprpc.Muxer, secret cc.SecretToken) []httprpc.Server {
	servers := httprpc.MustNewServers(
		miscService.Clone,
		crmService.Clone,
		vtigerService.Clone,
		vhtService.Clone,
	)
	return servers
}
