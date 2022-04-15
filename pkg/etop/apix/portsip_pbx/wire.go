package portsip_pbx

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.Struct(new(PortsipService), "*"),
)
