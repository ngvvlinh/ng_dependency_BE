package v2

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	New,
)
