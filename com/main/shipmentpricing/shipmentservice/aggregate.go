package shipmentservice

import (
	"context"

	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/convert"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/model"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
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

func AggregateMessageBus(a *Aggregate) shipmentservice.CommandBus {
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "service_ids không được để trống")
	}
	if err := checkValidOtherCondition(args.OtherCondition); err != nil {
		return nil, err
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
	if err := checkValidOtherCondition(args.OtherCondition); err != nil {
		return err
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

func checkValidOtherCondition(cond *shipmentservice.OtherCondition) error {
	// Require min_weight, max_weight
	// max_weight = -1 -> unlimited
	if cond != nil {
		minWeight := cond.MinWeight
		maxWeight := cond.MaxWeight
		if minWeight < 0 || (maxWeight != -1 && maxWeight < 0) {
			return cm.Errorf(cm.InvalidArgument, nil, "Cấu hình khối lượng không hợp lệ")
		}
	}
	return nil
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

func (a *Aggregate) UpdateShipmentServicesLocationConfig(ctx context.Context, args *shipmentservice.UpdateShipmentServicesLocationConfigArgs) (updated int, _ error) {
	if len(args.IDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing IDs")
	}
	services, err := a.shipmentServiceStore(ctx).IDs(args.IDs...).ListShipmentServices()
	if err != nil {
		return 0, err
	}
	mapServices := make(map[dot.ID]*shipmentservice.ShipmentService, len(services))
	for _, s := range services {
		mapServices[s.ID] = s
	}
	for _, id := range args.IDs {
		if _, ok := mapServices[id]; !ok {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "ShipmentService not found (ID = %v)", id)
		}
	}

	for _, bl := range args.BlacklistLocations {
		if bl.ShippingLocationType == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn loại địa chỉ: lấy hàng (pick) hoặc giao hàng (deliver)")
		}
		if bl.Reason == "" {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng điền nguyên nhận chặn các địa điểm trong blacklist")
		}
	}

	for _, al := range args.AvailableLocations {
		if al.ShippingLocationType == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn loại địa chỉ: lấy hàng (pick) hoặc giao hàng (deliver)")
		}
		if al.FilterType == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn phương thức: bao gồm (include) hoặc loại trừ (exclude)")
		}
	}

	_service := &shipmentservice.ShipmentService{
		AvailableLocations: args.AvailableLocations,
		BlacklistLocations: args.BlacklistLocations,
	}
	var serviceUpdate model.ShipmentService
	if err := scheme.Convert(_service, &serviceUpdate); err != nil {
		return 0, err
	}
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.shipmentServiceStore(ctx).IDs(args.IDs...).UpdateShipmentServiceDB(&serviceUpdate)
		if err != nil {
			return err
		}

		for _, id := range args.IDs {
			s := mapServices[id]
			if err := DeleteRedisCache(ctx, a.redisStore, s.ConnectionID, s.ServiceIDs); err != nil {
				return err
			}
		}
		return nil
	})
	return
}
