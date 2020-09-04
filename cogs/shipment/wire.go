// +build wireinject

package _shipment

import (
	"github.com/google/wire"

	loggingshippingwebhook "o.o/backend/com/etc/logging/shippingwebhook"
	"o.o/backend/com/main/shipmentpricing"
	"o.o/backend/com/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
)

var WireSet = wire.NewSet(
	shipmentpricing.WireSet,
	shippingcarrier.WireSet,
	shipping.WireSet,
	loggingshippingwebhook.WireSet,
)
