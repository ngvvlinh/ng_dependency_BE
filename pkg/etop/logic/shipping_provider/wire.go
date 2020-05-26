package shipping_provider

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewCtrl,
)
