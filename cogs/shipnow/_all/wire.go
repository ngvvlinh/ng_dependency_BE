// +build wireinject

package shipnow_all

import (
	"github.com/google/wire"

	_shipnow "o.o/backend/cogs/shipnow"
	_ahamove "o.o/backend/cogs/shipnow/ahamove"
	"o.o/backend/com/main/accountshipnow"
)

var WireSet = wire.NewSet(
	_shipnow.WireSet,
	_ahamove.WireSet,
	accountshipnow.WireSet,
	SupportedShipnowCarrierConfig,
	SupportedShipnowCarrierDriver,
	AhamoveConfig,
)
