package aggregate

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping/carrying"
	"o.o/api/shopping/tradering"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/carrying/convert"
	"o.o/backend/com/shopping/carrying/model"
	"o.o/backend/com/shopping/carrying/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ carrying.Aggregate = &CarrierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type CarrierAggregate struct {
	store    sqlstore.CarrierStoreFactory
	eventBus capi.EventBus
}

func NewCarrierAggregate(eventBus capi.EventBus, db com.MainDB) *CarrierAggregate {
	return &CarrierAggregate{
		store:    sqlstore.NewCarrierStore(db),
		eventBus: eventBus,
	}
}

func CarrierAggregateMessageBus(a *CarrierAggregate) carrying.CommandBus {
	b := bus.New()
	return carrying.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *CarrierAggregate) CreateCarrier(
	ctx context.Context, args *carrying.CreateCarrierArgs,
) (*carrying.ShopCarrier, error) {
	carrier := new(carrying.ShopCarrier)
	if err := scheme.Convert(args, carrier); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateCarrier(carrier)
	return carrier, err
}

func (a *CarrierAggregate) UpdateCarrier(
	ctx context.Context, args *carrying.UpdateCarrierArgs,
) (*carrying.ShopCarrier, error) {
	carrier, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetCarrier()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(args, carrier); err != nil {
		return nil, err
	}
	carrierDB := new(model.ShopCarrier)
	if err := scheme.Convert(carrier, carrierDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateCarrierDB(carrierDB)
	return carrier, err
}

func (a *CarrierAggregate) DeleteCarrier(
	ctx context.Context, id dot.ID, shopID dot.ID,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(id).ShopID(shopID).SoftDelete()
	event := &tradering.TraderDeletedEvent{
		EventMeta: meta.NewEvent(),
		ShopID:    shopID,
		TraderID:  id,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return 0, err
	}
	return deleted, err
}
