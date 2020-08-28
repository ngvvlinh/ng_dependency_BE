package carrier

import (
	"github.com/google/wire"
	shipnowcarrier "o.o/api/main/shipnow/carrier"
)

var WireSet = wire.NewSet(
	wire.Bind(new(shipnowcarrier.Manager), new(*ShipnowManager)),
	NewShipnowManager,
)
