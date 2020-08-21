// +build wireinject

package shipnow_all

import (
	"github.com/google/wire"

	_shipnow "o.o/backend/cogs/shipnow"
	_ahamove "o.o/backend/cogs/shipnow/ahamove"
)

var WireSet = wire.NewSet(
	_shipnow.WireSet,
	_ahamove.WireSet,
	AhamoveConfig,
	AllSupportedShipnowCarriers,
)
