package pricelist

import (
	"context"

	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/shipmentpricing/pricelist/convert"
	"etop.vn/backend/com/main/shipmentpricing/pricelist/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ pricelist.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                     cmsql.Transactioner
	eventBus               capi.EventBus
	shipmentPriceListStore sqlstore.ShipmentPriceListStoreFactory
}

func NewAggregate(db *cmsql.Database, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:                     db,
		eventBus:               eventBus,
		shipmentPriceListStore: sqlstore.NewShipmentPriceListStore(db),
	}
}

func (a *Aggregate) MessageBus() pricelist.CommandBus {
	b := bus.New()
	return pricelist.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentPriceList(ctx context.Context, args *pricelist.CreateShipmentPriceListArg) (*pricelist.ShipmentPriceList, error) {
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
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
			err = a.ActivateShipmentPriceList(ctx, plist.ID)
		}
		return err
	})
	return res, _err
}

func (a *Aggregate) UpdateShipmentPriceList(ctx context.Context, args *pricelist.UpdateShipmentPriceListArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	var pricelist pricelist.ShipmentPriceList
	if err := scheme.Convert(args, &pricelist); err != nil {
		return err
	}
	return a.shipmentPriceListStore(ctx).UpdateShipmentPriceList(&pricelist)
}

func (a *Aggregate) ActivateShipmentPriceList(ctx context.Context, id dot.ID) error {
	if id == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err := a.shipmentPriceListStore(ctx).DeactivePriceList()
		if err != nil {
			return err
		}
		err = a.shipmentPriceListStore(ctx).ID(id).ActivePriceList()
		if err != nil {
			return err
		}

		event := &pricelist.ShipmentPriceListActivatedEvent{
			EventMeta: meta.NewEvent(),
			ID:        id,
		}
		return a.eventBus.Publish(ctx, event)
	})
}

func (a *Aggregate) DeleteShipmentPriceList(ctx context.Context, id dot.ID) error {
	_, err := a.shipmentPriceListStore(ctx).ID(id).SoftDelete()
	return err
}
