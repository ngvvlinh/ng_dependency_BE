package aggregate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/carrying/convert"
	"etop.vn/backend/com/shopping/carrying/model"
	"etop.vn/backend/com/shopping/carrying/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/capi"
)

var _ carrying.Aggregate = &CarrierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type CarrierAggregate struct {
	store    sqlstore.CarrierStoreFactory
	eventBus capi.EventBus
}

func NewCarrierAggregate(eventBus capi.EventBus, db *cmsql.Database) *CarrierAggregate {
	return &CarrierAggregate{
		store:    sqlstore.NewCarrierStore(db),
		eventBus: eventBus,
	}
}

func (a *CarrierAggregate) MessageBus() carrying.CommandBus {
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
	ctx context.Context, id int64, shopID int64,
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
