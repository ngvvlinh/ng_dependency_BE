package shipmentservice

import (
	"context"

	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/sqlstore"
	"o.o/backend/com/main/shipmentpricing/util"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

var _ shipmentservice.QueryService = &QueryService{}

type QueryService struct {
	redisStore           redis.Store
	shipmentServiceStore sqlstore.ShipmentServiceStoreFactory
}

func NewQueryService(db com.MainDB, redisStore redis.Store) *QueryService {
	return &QueryService{
		redisStore:           redisStore,
		shipmentServiceStore: sqlstore.NewShipmentServiceStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) shipmentservice.QueryBus {
	b := bus.New()
	return shipmentservice.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetShipmentService(ctx context.Context, id dot.ID) (*shipmentservice.ShipmentService, error) {
	return q.shipmentServiceStore(ctx).ID(id).GetShipmentService()
}

func (q *QueryService) GetShipmentServiceByServiceID(ctx context.Context, serviceID string, connID dot.ID) (res *shipmentservice.ShipmentService, err error) {
	key := getShipmentServiceRedisKey(ctx, serviceID, connID)
	err = q.redisStore.Get(key, &res)
	if err != nil {
		res, err = q.shipmentServiceStore(ctx).ServiceID(serviceID).ConnectionID(connID).Status(status3.P).GetShipmentService()
		_ = q.redisStore.SetWithTTL(key, res, util.DefaultTTL)
	}
	if res == nil {
		return nil, cm.Errorf(cm.NotFound, nil, "")
	}
	return res, err
}

func (q *QueryService) ListShipmentServices(ctx context.Context, args *shipmentservice.ListShipmentServicesArgs) ([]*shipmentservice.ShipmentService, error) {
	return q.shipmentServiceStore(ctx).OptionalConnectionID(args.ConnectionID).ListShipmentServices()
}
