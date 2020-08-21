// +build wireinject

package _shipment

import (
	"github.com/google/wire"

	loggingshippingwebhook "o.o/backend/com/etc/logging/shippingwebhook"
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
	loggingshippingwebhook.WireSet,
)
