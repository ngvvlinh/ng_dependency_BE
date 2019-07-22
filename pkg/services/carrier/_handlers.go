package carrier

import (
	"context"

	"etop.vn/api/main/carrier"
	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/models/shipping"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var ll = l.New()
var _ carrier.ShippingCarrierManager = &Impl{}
var _ carrier.Utilities = &Impl{}

type Impl struct {
	self     shipping.ProcessManagerBus
	location location.QueryBus
}

func New(locationBus location.QueryBus) *Impl {
	im := &Impl{
		location: locationBus,
	}
	im.self = shipping.ProcessManagerBus{bus.New()}
	im.RegisterHandlers(im.self)
	return im
}

func (im *Impl) MessageBus() shipping.ProcessManagerBus {
	return im.self
}

func (im *Impl) RegisterHandlers(bus bus.Bus) {
	bus.AddHandlers(
		im.GetShippingServicesHandler,
		im.CreateExternalShipmentHandler,
		im.CancelExternalShipmentHandler,
	)
}

func (im *Impl) GetShippingServicesHandler(ctx context.Context, query *carrier.GetShippingServiceQuery) error {
	result, err := im.GetShippingServices(ctx, query)
	query.Result = result
	return err
}

func (im *Impl) GetShippingServices(ctx context.Context, query *carrier.GetShippingServiceQuery) ([]*shipping.ShippingService, error) {
	return nil, cm.ErrTODO
}

func (im *Impl) CreateExternalShipmentHandler(ctx context.Context, cmd *carrier.CreateExternalShipmentCommand) error {
	return im.CreateExternalShipment(ctx, cmd)
}

func (im *Impl) CreateExternalShipment(ctx context.Context, cmd *carrier.CreateExternalShipmentCommand) error {
	return cm.ErrTODO
}

func (im *Impl) CancelExternalShipmentHandler(ctx context.Context, cmd *carrier.CancelExternalShipmentCommand) error {
	return im.CancelExternalShipment(ctx, cmd)
}

func (im *Impl) CancelExternalShipment(ctx context.Context, cmd *carrier.CancelExternalShipmentCommand) error {
	return cm.ErrTODO
}

func (im *Impl) ParseServiceCode(carrier string, code string) (serviceName string, err error) {
	return "", nil
}
