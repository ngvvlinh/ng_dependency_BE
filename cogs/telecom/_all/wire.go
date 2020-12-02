// +build wireinject

package _all

import (
	"github.com/google/wire"
	etelecomprovider "o.o/backend/com/etelecom/provider"
)

var WireSet = wire.NewSet(
	etelecomprovider.WireSet,
	SupportedTelecomDriver,
)
