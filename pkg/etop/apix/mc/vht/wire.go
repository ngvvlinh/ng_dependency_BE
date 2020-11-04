package vht

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewServers,
	wire.Struct(new(VHTUserService), "*"),
)
