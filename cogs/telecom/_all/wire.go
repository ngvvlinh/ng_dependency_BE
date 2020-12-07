// +build wireinject

package _all

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom"
	etelecomprovider "o.o/backend/com/etelecom/provider"
)

var WireSet = wire.NewSet(
	etelecom.WireSet,
	etelecomprovider.WireSet,
	SupportedTelecomDriver,
)
