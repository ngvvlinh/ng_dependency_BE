package shopshipmentpricelist

import (
	"context"

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
}

func NewAggregate(db com.MainDB) *Aggregate {
	return &Aggregate{
		db:                 db,
		shopPriceListStore: sqlstore.NewShopPriceListStore(db),
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

	var pList shopshipmentpricelist.ShopShipmentPriceList
	if err := scheme.Convert(args, &pList); err != nil {
		return err
	}
	return a.shopPriceListStore(ctx).UpdateShopPriceList(&pList)
}

func (a *Aggregate) DeleteShopShipmentPriceList(ctx context.Context, shopID dot.ID) error {
	_, err := a.shopPriceListStore(ctx).ShopID(shopID).GetShopPriceList()
	if err != nil {
		return err
	}
	_, err = a.shopPriceListStore(ctx).ShopID(shopID).SoftDelete()
	return err
}
