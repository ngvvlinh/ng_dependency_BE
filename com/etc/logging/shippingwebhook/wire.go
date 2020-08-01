package shippingwebhook

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregate,
)
