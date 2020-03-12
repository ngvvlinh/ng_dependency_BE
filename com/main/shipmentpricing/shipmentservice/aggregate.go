package shipmentservice

import (
	"context"

	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/shipmentpricing/shipmentservice/convert"
	"etop.vn/backend/com/main/shipmentpricing/shipmentservice/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ shipmentservice.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                   cmsql.Transactioner
	redisStore           redis.Store
	shipmentServiceStore sqlstore.ShipmentServiceStoreFactory
}

func NewAggregate(db *cmsql.Database, redisStore redis.Store) *Aggregate {
	return &Aggregate{
		db:                   db,
		redisStore:           redisStore,
		shipmentServiceStore: sqlstore.NewShipmentServiceStore(db),
	}
}

func (a *Aggregate) MessageBus() shipmentservice.CommandBus {
	b := bus.New()
	return shipmentservice.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentService(ctx context.Context, args *shipmentservice.CreateShipmentServiceArgs) (*shipmentservice.ShipmentService, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection_id")
	}
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	if len(args.ServiceIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing service_ids")
	}
	var service shipmentservice.ShipmentService
	if err := scheme.Convert(args, &service); err != nil {
		return nil, err
	}
	service.ID = cm.NewID()
	service.Status = status3.P

	var res *shipmentservice.ShipmentService
	var err error
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentServiceStore(ctx).CreateShipmentService(&service)
		if err != nil {
			return err
		}
		return DeleteRedisCache(ctx, a.redisStore, res.ConnectionID, res.ServiceIDs)
	})
	return res, err
}

func (a *Aggregate) UpdateShipmentService(ctx context.Context, args *shipmentservice.UpdateShipmentServiceArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	service, err := a.shipmentServiceStore(ctx).ID(args.ID).GetShipmentService()
	if err != nil {
		return err
	}
	var newService shipmentservice.ShipmentService
	if err := scheme.Convert(args, &newService); err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err := a.shipmentServiceStore(ctx).UpdateShipmentService(&newService)
		if err != nil {
			return err
		}

		serviceIDs := newService.ServiceIDs
		serviceIDs = append(serviceIDs, service.ServiceIDs...)
		return DeleteRedisCache(ctx, a.redisStore, service.ConnectionID, serviceIDs)
	})
}

func (a *Aggregate) DeleteShipmentService(ctx context.Context, id dot.ID) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		service, err := a.shipmentServiceStore(ctx).ID(id).GetShipmentService()
		if err != nil {
			return err
		}
		_, err = a.shipmentServiceStore(ctx).ID(id).SoftDelete()
		if err != nil {
			return err
		}
		return DeleteRedisCache(ctx, a.redisStore, service.ConnectionID, service.ServiceIDs)
	})
	return nil
}
