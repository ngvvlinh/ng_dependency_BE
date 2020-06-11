package pricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/meta"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/pricelist/convert"
	"o.o/backend/com/main/shipmentpricing/pricelist/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ pricelist.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                     *cmsql.Database
	eventBus               capi.EventBus
	shipmentPriceListStore sqlstore.ShipmentPriceListStoreFactory
	shopPriceListQS        shopshipmentpricelist.QueryBus
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus, shopPriceListQS shopshipmentpricelist.QueryBus) *Aggregate {
	return &Aggregate{
		db:                     db,
		eventBus:               eventBus,
		shipmentPriceListStore: sqlstore.NewShipmentPriceListStore(db),
		shopPriceListQS:        shopPriceListQS,
	}
}

func AggregateMessageBus(a *Aggregate) pricelist.CommandBus {
	b := bus.New()
	return pricelist.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentPriceList(ctx context.Context, args *pricelist.CreateShipmentPriceListArg) (*pricelist.ShipmentPriceList, error) {
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID")
	}
	var plist pricelist.ShipmentPriceList
	if err := scheme.Convert(args, &plist); err != nil {
		return nil, err
	}
	plist.ID = cm.NewID()
	var res *pricelist.ShipmentPriceList
	var err error
	_err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		res, err = a.shipmentPriceListStore(ctx).CreateShipmentPriceList(&plist)
		if err != nil {
			return err
		}
		if args.IsActive {
			err = a.ActivateShipmentPriceList(ctx, plist.ID, plist.ConnectionID)
		}
		return err
	})
	return res, _err
}

func (a *Aggregate) UpdateShipmentPriceList(ctx context.Context, args *pricelist.UpdateShipmentPriceListArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	var priceList pricelist.ShipmentPriceList
	if err := scheme.Convert(args, &priceList); err != nil {
		return err
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.shipmentPriceListStore(ctx).UpdateShipmentPriceList(&priceList); err != nil {
			return err
		}
		return a.deleteCachePriceList(ctx, args.ID)
	})
}

// ActivateShipmentPriceList
//
// Mỗi connection chỉ có duy nhất 1 bảng giá được active
func (a *Aggregate) ActivateShipmentPriceList(ctx context.Context, id dot.ID, connectionID dot.ID) error {
	if id == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err := a.shipmentPriceListStore(ctx).ConnectionID(connectionID).DeactivePriceList()
		if err != nil {
			return err
		}
		err = a.shipmentPriceListStore(ctx).ID(id).ConnectionID(connectionID).ActivePriceList()
		if err != nil {
			return err
		}

		return a.deleteCachePriceList(ctx, id)
	})
}

func (a *Aggregate) DeleteShipmentPriceList(ctx context.Context, id dot.ID) error {
	priceList, err := a.shipmentPriceListStore(ctx).ID(id).GetShipmentPriceList()
	if err != nil {
		return err
	}
	if priceList.IsActive {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not delete an active price list")
	}

	queryShopPriceList := &shopshipmentpricelist.ListShopShipmentPriceListsQuery{
		ShipmentPriceListID: id,
	}
	if err := a.shopPriceListQS.Dispatch(ctx, queryShopPriceList); err != nil {
		return err
	}
	if len(queryShopPriceList.Result.ShopShipmentPriceLists) != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "This price list is used for shop. Can not delete it.")
	}

	_, err = a.shipmentPriceListStore(ctx).ID(id).SoftDelete()
	if err != nil {
		return err
	}
	return a.deleteCachePriceList(ctx, id)
}

func (a *Aggregate) deleteCachePriceList(ctx context.Context, priceListID dot.ID) error {
	event := &pricelist.DeleteCachePriceListEvent{
		EventMeta:           meta.NewEvent(),
		ShipmentPriceListID: priceListID,
	}
	return a.eventBus.Publish(ctx, event)
}
