package xshop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/idemp"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/convertpb"
)

var (
	idempgroup        *idemp.RedisGroup
	locationQuery     location.QueryBus
	customerQuery     *customering.QueryBus
	customerAggregate *customering.CommandBus
	addressQuery      *addressing.QueryBus
	addressAggregate  *addressing.CommandBus
	inventoryQuery    *inventory.QueryBus
	catalogQuery      *catalog.QueryBus
	catalogAggregate  *catalog.CommandBus
)

const PrefixIdempShopAPI = "IdempShopAPI"

func init() {
	bus.AddHandlers("apix",
		miscService.VersionInfo,
		miscService.CurrentAccount,
	)
}

func Init(
	sd cmservice.Shutdowner,
	rd redis.Store,
	locationQ location.QueryBus,
	customerQ *customering.QueryBus,
	customerA *customering.CommandBus,
	addressQ *addressing.QueryBus,
	addressA *addressing.CommandBus,
	inventoryQ *inventory.QueryBus,
	catalogQ *catalog.QueryBus,
	catalogA *catalog.CommandBus,
) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	sd.Register(idempgroup.Shutdown)

	locationQuery = locationQ
	customerQuery = customerQ
	customerAggregate = customerA
	addressQuery = addressQ
	addressAggregate = addressA
	inventoryQuery = inventoryQ
	catalogAggregate = catalogA
	catalogQuery = catalogQ
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "shop",
		Version: "1.0.0",
	}
	return nil
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *CurrentAccountEndpoint) error {
	if q.Context.Shop == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = convertpb.PbPublicAccountInfo(q.Context.Shop)
	return nil
}
