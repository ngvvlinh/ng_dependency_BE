package partnercarrier

import (
	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	"o.o/backend/pkg/common/apifw/idemp"
	cmService "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/redis"
)

var (
	idempgroup      *idemp.RedisGroup
	connectionQuery connectioning.QueryBus
	connectionAggr  connectioning.CommandBus
	shippingAggr    shipping.CommandBus
)

const PrefixIdempPartnerCarrierAPI = "IdempPartnerCarrierAPI"

func Init(
	sd cmService.Shutdowner,
	rd redis.Store,
	connQuery connectioning.QueryBus,
	connAggr connectioning.CommandBus,
	shippingAggregate shipping.CommandBus,
) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerCarrierAPI, 0)
	sd.Register(idempgroup.Shutdown)

	connectionQuery = connQuery
	connectionAggr = connAggr
	shippingAggr = shippingAggregate
}
