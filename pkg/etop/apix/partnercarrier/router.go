package partnercarrier

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
)

var idempgroup *idemp.RedisGroup

const PrefixIdempPartnerCarrierAPI = "IdempPartnerCarrierAPI"

type Servers []httprpc.Server

func NewServers(
	rd redis.Store,
	miscService *MiscService,
	shipmentConnectionService *ShipmentConnectionService,
	shipmentService *ShipmentService,
) (Servers, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerCarrierAPI, 0)
	servers := httprpc.MustNewServers(
		miscService.Clone,
		shipmentConnectionService.Clone,
		shipmentService.Clone,
	)
	return servers, idempgroup.Shutdown
}
