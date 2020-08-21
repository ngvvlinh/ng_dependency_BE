// +build wireinject

package shipment_all

import (
	"github.com/google/wire"

	_shipment "o.o/backend/cogs/shipment"
	ghnv2 "o.o/backend/cogs/shipment/ghn/v2"
	"o.o/backend/pkg/integration/shipping/ghn"
)

var WireSet = wire.NewSet(
	_shipment.WireSet,
	ghnv2.WireSet,
	wire.FieldsOf(new(Config), "GHN", "GHNWebhook"),
	ghn.New,
	SupportedCarrierDrivers,
	SupportedShippingCarrierConfig,
	SupportedCarrierDriver,
	SupportedShipmentServices,
)
