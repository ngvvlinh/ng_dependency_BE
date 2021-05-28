package telecom

import "github.com/google/wire"

var WireSet = wire.NewSet(
	BindTelecomStore,
	wire.Struct(new(TelecomStore), "*"),
)
