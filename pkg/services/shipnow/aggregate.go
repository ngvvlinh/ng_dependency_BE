package shipnow

import (
	"context"

	"etop.vn/api/main/location"
	"etop.vn/api/main/shipnow"
	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/shipnow/convert"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.ShipnowAggregate = &Aggregate{}

type Aggregate struct {
	bus      bus.Bus
	location location.Bus

	pm ProcessManager
	s  *sqlstore.ShipnowStore
}

func NewAggregate(db cmsql.Database, pm ProcessManager, eventBus bus.Bus, location location.Bus) *Aggregate {
	return &Aggregate{
		bus:      eventBus,
		location: location,
		pm:       pm,
		s:        sqlstore.NewShipnowStore(db),
	}
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentCommand) (*shipnow.ShipnowFulfillment, error) {
	shipnowFfm := cmd.ShipnowFulfillment
	shipnowFfm.Id = cm.NewID()

	if err := a.pm.HandleCreation(ctx, &shipnowFfm); err != nil {
		return nil, err
	}

	modelShipnowFfm := convert.ShipnowToModel(&shipnowFfm)
	if err := modelShipnowFfm.Validate(); err != nil {
		return nil, err
	}

	// TODO: store to db

	eventData := &shipnowv1.EventData_Created{}
	event := shipnow.NewShipnowEvent(meta.NewUUID(), shipnowFfm.Id, eventData)
	if err := a.bus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return convert.Shipnow(modelShipnowFfm), nil
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentCommand) (*meta.Empty, error) {
	panic("implement me")
}

func (a *Aggregate) CancelShipnowFulfillment(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentCommand) (*meta.Empty, error) {
	panic("implement me")
}
