package partnercarrier

import (
	"o.o/backend/pkg/common/apifw/idemp"
	cmService "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
)

var idempgroup *idemp.RedisGroup

const PrefixIdempPartnerCarrierAPI = "IdempPartnerCarrierAPI"

type Servers []httprpc.Server

func NewServers(
	sd cmService.Shutdowner,
	rd redis.Store,
	miscService *MiscService,
	shipmentConnectionService *ShipmentConnectionService,
	shipmentService *ShipmentService,
) Servers {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerCarrierAPI, 0)
	sd.Register(idempgroup.Shutdown)

	servers := httprpc.MustNewServers(
		miscService.Clone,
		shipmentConnectionService.Clone,
		shipmentService.Clone,
	)
	return servers
}
