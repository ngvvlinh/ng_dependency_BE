// +build wireinject

package subscripting

import (
	"github.com/google/wire"

	"o.o/backend/com/subscripting/subscription"
	"o.o/backend/com/subscripting/subscriptionplan"
	"o.o/backend/com/subscripting/subscriptionproduct"
)

var WireSet = wire.NewSet(
	subscription.WireSet,
	subscriptionplan.WireSet,
	subscriptionproduct.WireSet,
)
