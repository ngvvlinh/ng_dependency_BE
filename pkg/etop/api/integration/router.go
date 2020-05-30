package integration

import (
	service "o.o/api/top/int/integration"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/integration
// +gen:wrapper:package=integration

type Servers []httprpc.Server

func NewIntegrationServer(
	rd redis.Store,
	miscService *MiscService,
	integrationService *IntegrationService,
) (Servers, func()) {
	idempgroup = idemp.NewRedisGroup(rd, api.PrefixIdempUser, 0)
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewIntegrationServiceServer(WrapIntegrationService(integrationService.Clone)),
	}
	return servers, idempgroup.Shutdown
}
