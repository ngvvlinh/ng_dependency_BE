package integration

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewIntegrationServer(
	rd redis.Store,
	miscService *MiscService,
	integrationService *IntegrationService,
) (Servers, func()) {
	idempgroup = idemp.NewRedisGroup(rd, apiroot.PrefixIdempUser, 0)
	servers := httprpc.MustNewServers(
		miscService.Clone,
		integrationService.Clone,
	)
	return servers, idempgroup.Shutdown
}
