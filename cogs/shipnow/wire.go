// +build wireinject

package _shipnow

import (
	"github.com/google/wire"

	shipnowcarrier "o.o/api/main/shipnow/carrier"
	shipnowmanager "o.o/backend/com/main/shipnow/carrier"
	comcarrier "o.o/backend/com/main/shipnowcarrier"
)

var WireSet = wire.NewSet(
	wire.Bind(new(shipnowcarrier.Manager), new(*comcarrier.ShipnowManager)),
	comcarrier.NewManager,
	shipnowmanager.WireSet,
)
