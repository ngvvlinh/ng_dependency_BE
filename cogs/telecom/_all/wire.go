// +build wireinject

package _all

import (
	"github.com/google/wire"
	"o.o/backend/com/etelecom"
	etelecomprovider "o.o/backend/com/etelecom/provider"
	"o.o/backend/com/etelecom/usersetting"
)

var WireSet = wire.NewSet(
	etelecom.WireSet,
	usersetting.WireSet,
	etelecomprovider.WireSet,
	SupportedTelecomDriver,
)
