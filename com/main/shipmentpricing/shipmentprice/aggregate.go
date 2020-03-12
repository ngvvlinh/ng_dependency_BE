package shipmentprice

import (
	"context"

	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice/convert"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice/model"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ shipmentprice.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                 cmsql.Transactioner
	redisStore         redis.Store
	shipmentPriceStore sqlstore.ShipmentPriceStoreFactory
}

func NewAggregate(db *cmsql.Database, redisStore redis.Store) *Aggregate {
	return &Aggregate{
		db:                 db,
		redisStore:         redisStore,
		shipmentPriceStore: sqlstore.NewShipmentPriceStore(db),
	}
}

func (a *Aggregate) MessageBus() shipmentprice.CommandBus {
	b := bus.New()
	return shipmentprice.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentPrice(ctx context.Context, args *shipmentprice.CreateShipmentPriceArgs) (res *shipmentprice.ShipmentPrice, err error) {
	if args.ShipmentServiceID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn gói vận chuyển")
	}

	var pricing shipmentprice.ShipmentPrice
	if err := scheme.Convert(args, &pricing); err != nil {
		return nil, err
	}
	pricing.ID = cm.NewID()
	_pricing, err := validateShipmentPrice(&pricing)
	if err != nil {
		return nil, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).CreateShipmentPrice(_pricing)
		if err != nil {
			return err
		}

		// delete redis cache
		return DeleteRedisCache(ctx, a.redisStore)
	})
	return
}

func (a *Aggregate) UpdateShipmentPrice(ctx context.Context, args *shipmentprice.UpdateShipmentPriceArgs) (res *shipmentprice.ShipmentPrice, err error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}
	var pricing shipmentprice.ShipmentPrice
	if err := scheme.Convert(args, &pricing); err != nil {
		return nil, err
	}
	_pricing, err := validateShipmentPrice(&pricing)
	if err != nil {
		return nil, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).UpdateShipmentPrice(_pricing)
		if err != nil {
			return err
		}
		// delete redis cache
		return DeleteRedisCache(ctx, a.redisStore)
	})
	return
}

func (a *Aggregate) DeleteShipmentPrice(ctx context.Context, id dot.ID) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		_, err := a.shipmentPriceStore(ctx).ID(id).SoftDelete()
		if err != nil {
			return err
		}
		// delete redis cache
		return DeleteRedisCache(ctx, a.redisStore)
	})
}

func validateShipmentPrice(pricing *shipmentprice.ShipmentPrice) (*shipmentprice.ShipmentPrice, error) {
	if len(pricing.RegionTypes) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn tuyến vận chuyển: nội tỉnh, nội miền hoặc liên miền").WithMeta("field", "region_types")
	}
	if len(pricing.UrbanTypes) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn khu vực: nội thành, ngoại thành 1 hoặc ngoại thành 2").WithMeta("field", "urban_types")
	}
	if len(pricing.Details) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cấu hình giá").WithMeta("field", "details")
	}
	return pricing, nil
}

func (a *Aggregate) UpdateShipmentPricesPriorityPoint(ctx context.Context, args *shipmentprice.UpdateShipmentPricesPriorityPointArgs) (updated int, err error) {
	if len(args.ShipmentPrices) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing shipment_prices")
	}
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, sp := range args.ShipmentPrices {
			update := &model.ShipmentPrice{
				ID:            sp.ID,
				PriorityPoint: sp.PriorityPoint,
			}
			if _err := a.shipmentPriceStore(ctx).UpdateShipmentPriceDB(update); _err != nil {
				return _err
			}
			updated++
		}
		// delete redis cache
		return DeleteRedisCache(ctx, a.redisStore)
	})
	return
}
