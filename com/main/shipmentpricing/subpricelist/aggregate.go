package subpricelist

import (
	"context"

	"o.o/api/main/shipmentpricing/subpricelist"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/subpricelist/convert"
	"o.o/backend/com/main/shipmentpricing/subpricelist/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ subpricelist.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                        *cmsql.Database
	eventBus                  capi.EventBus
	shipmentSubPriceListStore sqlstore.ShipmentSubPriceListStoreFactory
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:                        db,
		eventBus:                  eventBus,
		shipmentSubPriceListStore: sqlstore.NewShipmentSubPriceListStore(db),
	}
}

func AggregateMessageBus(q *Aggregate) subpricelist.CommandBus {
	b := bus.New()
	return subpricelist.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipmentSubPriceList(ctx context.Context, args *subpricelist.CreateSubPriceListArgs) (res *subpricelist.ShipmentSubPriceList, err error) {
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	var splist subpricelist.ShipmentSubPriceList
	if err := scheme.Convert(args, &splist); err != nil {
		return nil, err
	}
	splist.ID = cm.NewID()
	splist.Status = status3.P
	return a.shipmentSubPriceListStore(ctx).CreateShipmentSubPriceList(&splist)
}

func (a *Aggregate) UpdateShipmentSubPriceList(ctx context.Context, args *subpricelist.UpdateSubPriceListArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	var subPriceListUpdate subpricelist.ShipmentSubPriceList
	if err := scheme.Convert(args, &subPriceListUpdate); err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.shipmentSubPriceListStore(ctx).UpdateShipmentSubPriceList(&subPriceListUpdate); err != nil {
			return err
		}

		event := &subpricelist.ShipmentSubPriceListUpdatedEvent{
			EventMeta: meta.NewEvent(),
			ID:        args.ID,
		}
		return a.eventBus.Publish(ctx, event)
	})
}

func (a *Aggregate) DeleteShipmentSubPriceList(ctx context.Context, id dot.ID) error {
	eventDeleting := &subpricelist.ShipmentSubPriceListDeletingEvent{
		EventMeta: meta.NewEvent(),
		ID:        id,
	}
	if err := a.eventBus.Publish(ctx, eventDeleting); err != nil {
		return err
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if _, err := a.shipmentSubPriceListStore(ctx).ID(id).SoftDelete(); err != nil {
			return err
		}

		event := &subpricelist.ShipmentSubPriceListDeletedEvent{
			EventMeta: meta.NewEvent(),
			ID:        id,
		}
		return a.eventBus.Publish(ctx, event)
	})
}
