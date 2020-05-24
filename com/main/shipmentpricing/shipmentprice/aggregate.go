package shipmentprice

import (
	"context"

	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipmentpricing/shipmentprice/convert"
	"o.o/backend/com/main/shipmentpricing/shipmentprice/model"
	"o.o/backend/com/main/shipmentpricing/shipmentprice/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
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

func AggregateMessageBus(a *Aggregate) shipmentprice.CommandBus {
	b := bus.New()
	return shipmentprice.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentPrice(ctx context.Context, args *shipmentprice.CreateShipmentPriceArgs) (res *shipmentprice.ShipmentPrice, err error) {
	if args.ShipmentServiceID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn gói vận chuyển")
	}
	if args.ShipmentPriceListID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn bảng giá")
	}

	var pricing shipmentprice.ShipmentPrice
	if err := scheme.Convert(args, &pricing); err != nil {
		return nil, err
	}
	pricing.ID = cm.NewID()
	if err := validateShipmentPrice(&pricing); err != nil {
		return nil, err
	}
	pricing.Status = status3.P

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).CreateShipmentPrice(&pricing)
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
	sp, err := a.shipmentPriceStore(ctx).ID(args.ID).GetShipmentPrice()
	if err != nil {
		return nil, err
	}
	pricing := convert.Apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(args, sp)
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).UpdateShipmentPrice(pricing)
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

func validateShipmentPrice(pricing *shipmentprice.ShipmentPrice) error {
	if len(pricing.RegionTypes) == 0 &&
		len(pricing.CustomRegionIDs) == 0 &&
		len(pricing.ProvinceTypes) == 0 &&
		len(pricing.UrbanTypes) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cấu hình địa điểm áp dụng gói.").WithMeta("field", "region_types")
	}
	if len(pricing.Details) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cấu hình giá").WithMeta("field", "details")
	}
	return nil
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
