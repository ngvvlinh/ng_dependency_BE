// +build wireinject

package subscripting

import (
	"github.com/google/wire"

	"o.o/backend/com/subscripting/invoice"
	"o.o/backend/com/subscripting/subscription"
	"o.o/backend/com/subscripting/subscriptionplan"
	"o.o/backend/com/subscripting/subscriptionproduct"
)

var WireSet = wire.NewSet(
	subscription.WireSet,
	invoice.WireSet,
	subscriptionplan.WireSet,
	subscriptionproduct.WireSet,
)
