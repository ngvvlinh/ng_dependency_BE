package integration

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(IntegrationService), "*"),
	NewIntegrationServer,
)
