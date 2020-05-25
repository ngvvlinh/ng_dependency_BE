package shipmentprice

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/subpricelist"
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
	subPriceListQS     subpricelist.QueryBus
	priceListQS        pricelist.QueryBus
}

func NewAggregate(db com.MainDB, redisStore redis.Store, subPriceListQS subpricelist.QueryBus, priceListQS pricelist.QueryBus) *Aggregate {
	return &Aggregate{
		db:                 db,
		redisStore:         redisStore,
		shipmentPriceStore: sqlstore.NewShipmentPriceStore(db),
		subPriceListQS:     subPriceListQS,
		priceListQS:        priceListQS,
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
	if args.ShipmentSubPriceListID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn bảng giá (SubPriceList)")
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
		return a.deleteCachePriceList(ctx, 0, pricing.ID)
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
	oldShipmentSubPriceList := sp.ShipmentSubPriceListID

	pricing := convert.Apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(args, sp)
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceStore(ctx).UpdateShipmentPrice(pricing)
		if err != nil {
			return err
		}
		// delete redis cache
		if args.ShipmentSubPriceListID != 0 {
			if err := a.deleteCachePriceList(ctx, oldShipmentSubPriceList, 0); err != nil {
				return err
			}
		}
		return a.deleteCachePriceList(ctx, args.ShipmentSubPriceListID, 0)
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
		return a.deleteCachePriceList(ctx, 0, id)
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
		var shipmentPriceIDs []dot.ID
		for _, sp := range args.ShipmentPrices {
			shipmentPriceIDs = append(shipmentPriceIDs, sp.ID)
		}
		shipmentPrices, err := a.shipmentPriceStore(ctx).IDs(shipmentPriceIDs...).ListShipmentPrices()
		if err != nil {
			return err
		}
		var subPriceListIDs []dot.ID
		for _, sp := range shipmentPrices {
			subPriceListIDs = append(subPriceListIDs, sp.ShipmentSubPriceListID)
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
		queryPriceList := &pricelist.ListShipmentPriceListsQuery{
			SubShipmentPriceListIDs: subPriceListIDs,
		}
		if err := a.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
			return err
		}
		for _, pl := range queryPriceList.Result {
			if err := DeleteRedisCache(ctx, a.redisStore, pl.ID); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (a *Aggregate) deleteCachePriceList(ctx context.Context, subShipmentPriceListID dot.ID, shipmentPriceID dot.ID) error {
	if subShipmentPriceListID == 0 {
		shipmentPrice, err := a.shipmentPriceStore(ctx).ID(shipmentPriceID).GetShipmentPrice()
		if err != nil {
			return err
		}
		subShipmentPriceListID = shipmentPrice.ShipmentSubPriceListID
	}

	queryPriceList := &pricelist.ListShipmentPriceListsQuery{
		SubShipmentPriceListIDs: []dot.ID{subShipmentPriceListID},
	}
	if err := a.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
		return err
	}
	for _, pl := range queryPriceList.Result {
		if err := DeleteRedisCache(ctx, a.redisStore, pl.ID); err != nil {
			return err
		}
	}
	return nil
}
