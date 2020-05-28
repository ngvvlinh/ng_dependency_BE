package partnercarrier

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(ShipmentConnectionService), "*"),
	wire.Struct(new(ShipmentService), "*"),
	NewServers,
)
