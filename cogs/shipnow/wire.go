// +build wireinject

package _shipnow

import (
	"github.com/google/wire"

	shipnowmanager "o.o/backend/com/main/shipnow/carrier"
)

var WireSet = wire.NewSet(
	shipnowmanager.WireSet,
)
