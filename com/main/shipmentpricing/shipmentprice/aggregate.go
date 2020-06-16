package shipmentprice

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
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
	db                 *cmsql.Database
	redisStore         redis.Store
	shipmentPriceStore sqlstore.ShipmentPriceStoreFactory
	priceListQS        pricelist.QueryBus
	shipmentServiceQS  shipmentservice.QueryBus
}

func NewAggregate(db com.MainDB, redisStore redis.Store, priceListQS pricelist.QueryBus, shipmentServiceQS shipmentservice.QueryBus) *Aggregate {
	return &Aggregate{
		db:                 db,
		redisStore:         redisStore,
		shipmentPriceStore: sqlstore.NewShipmentPriceStore(db),
		priceListQS:        priceListQS,
		shipmentServiceQS:  shipmentServiceQS,
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn bảng giá (PriceList)")
	}
	if err := a.validateShipmentPriceConnection(ctx, args.ShipmentServiceID, args.ShipmentPriceListID); err != nil {
		return nil, err
	}
	if err := validateShipmentPriceAdditionalFees(args.AdditionalFees); err != nil {
		return nil, err
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
		return a.deleteCachePriceList(ctx, pricing.ShipmentPriceListID)
	})
	return
}

func (a *Aggregate) validateShipmentPriceConnection(ctx context.Context, shipmentServiceID, shipmentPriceListID dot.ID) error {
	queryShipmentService := &shipmentservice.GetShipmentServiceQuery{
		ID: shipmentServiceID,
	}
	if err := a.shipmentServiceQS.Dispatch(ctx, queryShipmentService); err != nil {
		return err
	}
	queryPriceList := &pricelist.GetShipmentPriceListQuery{
		ID: shipmentPriceListID,
	}
	if err := a.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
		return err
	}
	shipmentService := queryShipmentService.Result
	priceList := queryPriceList.Result
	if shipmentService.ConnectionID != priceList.ConnectionID {
		return cm.Errorf(cm.FailedPrecondition, nil, "Cấu hình giá không hợp lệ. Gói dịch vụ và bảng giá không thuộc cùng một nhà vận chuyển")
	}
	return nil
}

func (a *Aggregate) UpdateShipmentPrice(ctx context.Context, args *shipmentprice.UpdateShipmentPriceArgs) (res *shipmentprice.ShipmentPrice, err error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}
	if err := validateShipmentPriceAdditionalFees(args.AdditionalFees); err != nil {
		return nil, err
	}

	sp, err := a.shipmentPriceStore(ctx).ID(args.ID).GetShipmentPrice()
	if err != nil {
		return nil, err
	}
	oldShipmentPriceList := sp.ShipmentPriceListID
	if args.ShipmentPriceListID != 0 || args.ShipmentServiceID != 0 {
		if err := a.validateShipmentPriceConnection(ctx,
			cm.CoalesceID(args.ShipmentServiceID, sp.ShipmentServiceID),
			cm.CoalesceID(args.ShipmentPriceListID, sp.ShipmentPriceListID)); err != nil {
			return nil, err
		}
	}

	pricing := convert.Apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(args, sp)
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).UpdateShipmentPrice(pricing)
		if err != nil {
			return err
		}
		// delete redis cache
		return a.deleteCachePriceList(ctx, args.ShipmentPriceListID, oldShipmentPriceList)
	})
	return
}

func (a *Aggregate) DeleteShipmentPrice(ctx context.Context, id dot.ID) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		sp, err := a.shipmentPriceStore(ctx).ID(id).GetShipmentPrice()
		if err != nil {
			return err
		}
		if _, err = a.shipmentPriceStore(ctx).ID(id).SoftDelete(); err != nil {
			return err
		}
		// delete redis cache
		return a.deleteCachePriceList(ctx, sp.ShipmentPriceListID)
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

func validateShipmentPriceAdditionalFees(addFees []*shipmentprice.AdditionalFee) error {
	for _, fee := range addFees {
		if fee.FeeType == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "Cấu hình cước phí không hợp lệ. Loại cước phí không hợp lệ.")
		}
		for _, rule := range fee.Rules {
			if rule.MaxValue != shipmentprice.MaximumValue && rule.MinValue > rule.MaxValue {
				return cm.Errorf(cm.InvalidArgument, nil, "Cấu hình cước phí không hợp lệ. (min_value > max_value)")
			}
		}
	}
	return nil
}

func (a *Aggregate) UpdateShipmentPricesPriorityPoint(ctx context.Context, args *shipmentprice.UpdateShipmentPricesPriorityPointArgs) (updated int, err error) {
	if len(args.ShipmentPrices) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing shipment_prices")
	}
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var shipmentPriceIDs []dot.ID
		for _, sp := range args.ShipmentPrices {
			shipmentPriceIDs = append(shipmentPriceIDs, sp.ID)
		}
		shipmentPrices, err := a.shipmentPriceStore(ctx).IDs(shipmentPriceIDs...).ListShipmentPrices()
		if err != nil {
			return err
		}
		var priceListIDs []dot.ID
		for _, sp := range shipmentPrices {
			priceListIDs = append(priceListIDs, sp.ShipmentPriceListID)
			update := &model.ShipmentPrice{
				ID:            sp.ID,
				PriorityPoint: sp.PriorityPoint,
			}
			if _err := a.shipmentPriceStore(ctx).UpdateShipmentPriceDB(update); _err != nil {
				return _err
			}
			updated++
		}

		// delete cache shipmentpricelistID
		return a.deleteCachePriceList(ctx, shipmentPriceIDs...)
	})
	return
}

func (a *Aggregate) deleteCachePriceList(ctx context.Context, shipmentPriceIDs ...dot.ID) error {
	for _, spID := range shipmentPriceIDs {
		if spID != 0 {
			if err := DeleteRedisCache(ctx, a.redisStore, spID); err != nil {
				return err
			}
		}
	}
	return nil
}
