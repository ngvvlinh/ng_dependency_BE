package shipmentservice

import (
	"context"

	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/shipmentpricing"
	"etop.vn/backend/com/main/shipmentpricing/shipmentservice/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ shipmentservice.QueryService = &QueryService{}

type QueryService struct {
	redisStore           redis.Store
	shipmentServiceStore sqlstore.ShipmentServiceStoreFactory
}

func NewQueryService(db *cmsql.Database, redisStore redis.Store) *QueryService {
	return &QueryService{
		redisStore:           redisStore,
		shipmentServiceStore: sqlstore.NewShipmentServiceStore(db),
	}
}

func (q *QueryService) MessageBus() shipmentservice.QueryBus {
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
		_ = q.redisStore.SetWithTTL(key, res, shipmentpricing.DefaultTTL)
	}
	if res == nil {
		return nil, cm.Errorf(cm.NotFound, nil, "")
	}
	return res, err
}

func (q *QueryService) ListShipmentServices(ctx context.Context, _ *meta.Empty) ([]*shipmentservice.ShipmentService, error) {
	return q.shipmentServiceStore(ctx).ListShipmentServices()
}
