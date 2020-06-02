package _shipment

import (
	"github.com/google/wire"

	"o.o/backend/com/main/shipmentpricing"
	"o.o/backend/com/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/etop/logic/shipping_provider"
)

var WireSet = wire.NewSet(
	shipping_provider.WireSet,
	shipmentpricing.WireSet,
	shippingcarrier.WireSet,
	shipping.WireSet,
)
