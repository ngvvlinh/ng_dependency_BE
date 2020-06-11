package shopshipmentpricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/convert"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ shopshipmentpricelist.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                 *cmsql.Database
	shopPriceListStore sqlstore.ShopPriceListStoreFactory
	priceListQS        pricelist.QueryBus
}

func NewAggregate(db com.MainDB, priceListQS pricelist.QueryBus) *Aggregate {
	return &Aggregate{
		db:                 db,
		shopPriceListStore: sqlstore.NewShopPriceListStore(db),
		priceListQS:        priceListQS,
	}
}

func AggregateMessageBus(a *Aggregate) shopshipmentpricelist.CommandBus {
	b := bus.New()
	return shopshipmentpricelist.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShopShipmentPriceList(ctx context.Context, args *shopshipmentpricelist.CreateShopShipmentPriceListArgs) (*shopshipmentpricelist.ShopShipmentPriceList, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID")
	}
	if args.ShipmentPriceListID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipment price list ID")
	}
	queryPriceList := &pricelist.GetShipmentPriceListQuery{
		ID: args.ShipmentPriceListID,
	}
	if err := a.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
		return nil, err
	}
	if queryPriceList.Result.ConnectionID != args.ConnectionID {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
	}

	var pList shopshipmentpricelist.ShopShipmentPriceList
	if err := scheme.Convert(args, &pList); err != nil {
		return nil, err
	}
	return a.shopPriceListStore(ctx).CreateShopPriceList(&pList)
}

func (a *Aggregate) UpdateShopShipmentPriceList(ctx context.Context, args *shopshipmentpricelist.UpdateShopShipmentPriceListArgs) error {
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID")
	}

	if args.ShipmentPriceListID != 0 {
		queryPriceList := &pricelist.GetShipmentPriceListQuery{
			ID: args.ShipmentPriceListID,
		}
		if err := a.priceListQS.Dispatch(ctx, queryPriceList); err != nil {
			return err
		}
		if queryPriceList.Result.ConnectionID != args.ConnectionID {
			return cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
		}
	}
	var pList shopshipmentpricelist.ShopShipmentPriceList
	if err := scheme.Convert(args, &pList); err != nil {
		return err
	}
	return a.shopPriceListStore(ctx).UpdateShopPriceList(&pList)
}

func (a *Aggregate) DeleteShopShipmentPriceList(ctx context.Context, shopID, connID dot.ID) error {
	_, err := a.shopPriceListStore(ctx).ShopID(shopID).ConnectionID(connID).GetShopPriceList()
	if err != nil {
		return err
	}
	_, err = a.shopPriceListStore(ctx).ShopID(shopID).ConnectionID(connID).SoftDelete()
	return err
}
